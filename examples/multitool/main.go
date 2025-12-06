package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/message"
	"github.com/hybridgroup/yzma/pkg/template"
)

func main() {
	if err := handleFlags(); err != nil {
		showUsage()
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Load llama.cpp library
	if err := llama.Load(*libPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load llama library: %v\n", err)
		os.Exit(1)
	}

	if !*verbose {
		llama.LogSet(llama.LogSilent())
	}

	llama.Init()

	// Load model
	model, err := llama.ModelLoadFromFile(*modelFile, llama.ModelDefaultParams())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load model: %v\n", err)
		os.Exit(1)
	}
	defer llama.ModelFree(model)

	// Create context
	params := llama.ContextDefaultParams()
	params.NCtx = uint32(*contextSize)

	ctx, err := llama.InitFromModel(model, params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create context: %v\n", err)
		os.Exit(1)
	}
	defer llama.Free(ctx)

	// Get vocabulary and chat template
	vocab := llama.ModelGetVocab(model)
	chatTemplate := llama.ModelChatTemplate(model, "")

	// Create sampler
	sp := llama.DefaultSamplerParams()
	sp.Temp = float32(*temperature)

	sampler := llama.NewSampler(model, llama.DefaultSamplers, sp)

	defer llama.SamplerFree(sampler)

	// Define available tools
	tools := getToolDefinitions()

	// Run the multi-step tool-calling conversation
	runMultiStepToolConversation(ctx, model, vocab, sampler, chatTemplate, tools)
}

// formatToolsForPrompt formats the tool definitions as a JSON string for the system prompt
func formatToolsForPrompt(tools []Tool) string {
	toolsJSON, err := json.MarshalIndent(tools, "", "  ")
	if err != nil {
		return "[]"
	}
	return string(toolsJSON)
}

func runMultiStepToolConversation(ctx llama.Context, model llama.Model, vocab llama.Vocab, sampler llama.Sampler, chatTemplate string, tools []Tool) {
	fmt.Println("=== Multi-Step Tool Calling Example ===")
	fmt.Println()
	fmt.Printf("User: %s\n", *userQuestion)
	fmt.Println()

	// Format tools as JSON for the system prompt
	toolsJSON := formatToolsForPrompt(tools)

	// Build system prompt with tool definitions
	systemPrompt := fmt.Sprintf(`You are a helpful assistant with access to the following tools:

%s

When you need to use a tool, respond with a tool call in the following format:
<tool_call>
{"name": "function_name", "arguments": {"arg1": "value1", "arg2": "value2"}}
</tool_call>

You can make multiple tool calls to solve complex problems step by step.
After receiving all tool results, provide a final answer to the user.
Do not include tool calls in your final answer.`, toolsJSON)

	// Initialize messages with system prompt and user question
	messages := []message.Message{
		message.Chat{
			Role:    "system",
			Content: systemPrompt,
		},
		message.Chat{
			Role:    "user",
			Content: *userQuestion,
		},
	}

	// Maximum number of tool-calling iterations to prevent infinite loops
	maxIterations := 10
	iteration := 0

	for iteration < maxIterations {
		iteration++
		fmt.Printf("=== Iteration %d ===\n", iteration)

		// Apply template and generate response
		prompt, err := template.Apply(chatTemplate, messages, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to apply template: %v\n", err)
			return
		}

		if *verbose {
			fmt.Println("=== Generated Prompt ===")
			fmt.Println(prompt)
			fmt.Println("========================")
		}

		// Clear KV cache before each generation
		mem, _ := llama.GetMemory(ctx)
		llama.MemoryClear(mem, true)

		// Tokenize and generate
		tokens := llama.Tokenize(vocab, prompt, true, false)
		response := generate(ctx, vocab, sampler, tokens)

		fmt.Printf("Assistant: %s\n", response)
		fmt.Println()

		// Parse tool calls from response
		toolCalls := message.ParseToolCalls(response)

		// If no tool calls found, we have the final answer
		if len(toolCalls) == 0 {
			fmt.Println("=== Final Answer (no more tool calls) ===")
			fmt.Printf("%s\n", response)
			break
		}

		fmt.Printf("=== Detected %d Tool Call(s) ===\n", len(toolCalls))

		// Process each tool call
		for i, call := range toolCalls {
			argsJSON, _ := json.Marshal(call.Function.Arguments)
			fmt.Printf("  [%d] Function: %s\n", i+1, call.Function.Name)
			fmt.Printf("      Arguments: %s\n", string(argsJSON))

			// Execute the tool
			result, err := executeToolCall(call)
			if err != nil {
				fmt.Printf("      Error: %v\n", err)
				result = fmt.Sprintf("Error: %v", err)
			} else {
				fmt.Printf("      Result: %s\n", result)
			}

			// Add tool call and response to messages
			messages = append(messages, message.Tool{
				Role:      "assistant",
				ToolCalls: []message.ToolCall{call},
			})
			messages = append(messages, message.ToolResponse{
				Role:    "tool",
				Name:    call.Function.Name,
				Content: result,
			})
		}
		fmt.Println()
	}

	if iteration >= maxIterations {
		fmt.Println("=== Warning: Maximum iterations reached ===")
	}
}

func generate(ctx llama.Context, vocab llama.Vocab, sampler llama.Sampler, tokens []llama.Token) string {
	var response strings.Builder

	batch := llama.BatchGetOne(tokens)
	llama.Decode(ctx, batch)

	for i := 0; i < *predictSize; i++ {
		token := llama.SamplerSample(sampler, ctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}

		buf := make([]byte, 128)
		n := llama.TokenToPiece(vocab, token, buf, 0, true)
		response.Write(buf[:n])

		batch = llama.BatchGetOne([]llama.Token{token})
		llama.Decode(ctx, batch)
	}

	return strings.TrimSpace(response.String())
}

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

	// Run the tool-calling conversation
	runToolConversation(ctx, model, vocab, sampler, chatTemplate, tools)
}

// formatToolsForPrompt formats the tool definitions as a JSON string for the system prompt
func formatToolsForPrompt(tools []Tool) string {
	toolsJSON, err := json.MarshalIndent(tools, "", "  ")
	if err != nil {
		return "[]"
	}
	return string(toolsJSON)
}

func runToolConversation(ctx llama.Context, model llama.Model, vocab llama.Vocab, sampler llama.Sampler, chatTemplate string, tools []Tool) {
	fmt.Println("=== Tool Calling Example ===")
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

After receiving tool results, provide a final answer to the user.`, toolsJSON)

	// Step 1: Create initial message with user question
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

	// Tokenize and generate
	tokens := llama.Tokenize(vocab, prompt, true, false)
	response := generate(ctx, vocab, sampler, tokens)

	fmt.Printf("Assistant: %s\n", response)
	fmt.Println()

	// Step 2: Parse tool calls from response (simplified parsing)
	toolCalls := message.ParseToolCalls(response)

	if len(toolCalls) > 0 {
		fmt.Println("=== Detected Tool Calls ===")
		for _, call := range toolCalls {
			argsJSON, _ := json.Marshal(call.Function.Arguments)
			fmt.Printf("  Function: %s\n", call.Function.Name)
			fmt.Printf("  Arguments: %s\n", string(argsJSON))

			// Execute the tool
			result, err := executeToolCall(call)
			if err != nil {
				fmt.Printf("  Error: %v\n", err)
				continue
			}
			fmt.Printf("  Result: %s\n", result)
			fmt.Println()

			// Add tool response to messages
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

		// Step 3: Generate final response with tool results
		fmt.Println("=== Final Response ===")
		prompt, err = template.Apply(chatTemplate, messages, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to apply template: %v\n", err)
			return
		}

		if *verbose {
			fmt.Println("=== Generated Prompt with Tool Results ===")
			fmt.Println(prompt)
			fmt.Println("==========================================")
		}

		// Clear KV cache and regenerate
		mem, _ := llama.GetMemory(ctx)
		llama.MemoryClear(mem, true)

		tokens = llama.Tokenize(vocab, prompt, true, false)
		finalResponse := generate(ctx, vocab, sampler, tokens)

		fmt.Printf("Assistant: %s\n", finalResponse)
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

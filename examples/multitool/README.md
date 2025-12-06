# multitool

Example showing multi-step tool calls with LLMs.

This example demonstrates how to use tools/function calling in a loop, allowing the model to make multiple tool calls to solve complex problems step by step.

## Usage

```shell
go run ./examples/multitool -model /path/to/model.gguf -lib /path/to/llama/lib -question "the question"
```

## Example

```
$ go run ./examples/multitool -model ~/models/Qwen3-VL-2B-Instruct-Q8_0.gguf -question "Tell me what is (15 + 27) * 3, then tell me what is (50 + 33) * 4, then calculate the flying speed of a swallow. After that then tell me 100 * 43 times that flying speed."                                                                               
=== Multi-Step Tool Calling Example ===

User: Tell me what is (15 + 27) * 3, then tell me what is (50 + 33) * 4, then calculate the flying speed of a swallow. After that then tell me 100 * 43 times that flying speed.

=== Iteration 1 ===
Assistant: First, let's calculate (15 + 27) * 3:

<tool_call>
{"name": "add", "arguments": {"a": 15, "b": 27}}
</tool_call>
<tool_call>
{"name": "multiply", "arguments": {"a": 3, "b": 45}}
</tool_call>

Next, let's calculate (50 + 33) * 4:

<tool_call>
{"name": "add", "arguments": {"a": 50, "b": 33}}
</tool_call>
<tool_call>
{"name": "multiply", "arguments": {"a": 4, "b": 83}}
</tool_call>

Now, let's calculate the flying speed of a swallow. The average flying speed of a swallow is approximately 10 meters per second. 

Then, we'll calculate 100 * 43 times that flying speed:

<tool_call>
{"name": "multiply", "arguments": {"a": 100, "b": 43}}
</tool_call>
<tool_call>
{"name": "multiply", "arguments": {"a": 10, "b": 4300}}
</tool_call>

The final answer is the result of the last calculation. Let's do that.
<tool_call>
{"name": "multiply", "arguments": {"a": 10, "b": 4300}}
</tool_call>

=== Detected 7 Tool Call(s) ===
  [1] Function: add
      Arguments: {"a":"15","b":"27"}
      Result: 42.00
  [2] Function: multiply
      Arguments: {"a":"3","b":"45"}
      Result: 135.00
  [3] Function: add
      Arguments: {"a":"50","b":"33"}
      Result: 83.00
  [4] Function: multiply
      Arguments: {"a":"4","b":"83"}
      Result: 332.00
  [5] Function: multiply
      Arguments: {"a":"100","b":"43"}
      Result: 4300.00
  [6] Function: multiply
      Arguments: {"a":"10","b":"4300"}
      Result: 43000.00
  [7] Function: multiply
      Arguments: {"a":"10","b":"4300"}
      Result: 43000.00

=== Iteration 2 ===
Assistant: The result of (15 + 27) * 3 is 135.  
The result of (50 + 33) * 4 is 332.  
The flying speed of a swallow is not provided in the query, so I cannot calculate the exact flying speed. However, based on typical data, a swallow can fly at speeds of around 10 to 15 meters per second.  
The result of 100 * 43 times the flying speed of a swallow is 43000.00.  

Therefore, the final answer is 43000.00.

=== Final Answer (no more tool calls) ===
The result of (15 + 27) * 3 is 135.  
The result of (50 + 33) * 4 is 332.  
The flying speed of a swallow is not provided in the query, so I cannot calculate the exact flying speed. However, based on typical data, a swallow can fly at speeds of around 10 to 15 meters per second.  
The result of 100 * 43 times the flying speed of a swallow is 43000.00.  

Therefore, the final answer is 43000.00.
```

## Available Tools

- add - Add two numbers together
- multiply - Multiply two numbers together
- subtract - Subtract the second number from the first

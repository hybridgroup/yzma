# tooluse

This example shows how yzma can use models that support tool usage.

## Running

```shell
$ go run ./examples/tooluse/ -model ~/models/qwen2.5-0.5b-instruct-fp16.gguf
=== Tool Calling Example ===

User: What is 15 + 27?

Assistant: 15 + 27 = 42

The result of adding 15 and 27 is 42.
```

```shell
$ go run ./examples/tooluse/ -model ~/models/qwen2.5-0.5b-instruct-fp16.gguf -question="what is 9 times 9?"
=== Tool Calling Example ===                                                                              
                                                                                                          
User: what is 9 times 9?

Assistant: 9 times 9 is 81.

9 times 9 = 81

The final answer is 81.
```

# modelinfo

Shows model information.

## Running

```shell
go run ./examples/modelnfo/ /path/to/model.gguf
```

```shell
$ go run ./examples/modelinfo/ ~/models/gemma-3-1b-it-q4_0.gguf 
Model Description: gemma3 1B Q4_0
Model Size: 997007872 tensors
Model Has Encoder: false
Model Has Decoder: true
Model Is Recurrent: false
Model Is Hybrid: false
Model Metadata (22 entries):
  general.file_type: 2
  tokenizer.ggml.unknown_token_id: 3
  tokenizer.ggml.padding_token_id: 0
  tokenizer.ggml.bos_token_id: 2
  general.quantization_version: 2
  tokenizer.ggml.model: llama
  gemma3.attention.sliding_window: 1024
  general.architecture: gemma3
  tokenizer.ggml.eos_token_id: 1
  gemma3.feed_forward_length: 6912
  gemma3.context_length: 32768
  gemma3.block_count: 26
  gemma3.attention.value_length: 256
  gemma3.rope.scaling.factor: 1.000000
  tokenizer.chat_template: {{ bos_token }} {%- if messages[0]['role'] == 'system' -%} {%- if messages[0]['content'] is string -%} {%- set first_user_prefix = messages[0]['content'] + '\n' -%} {%- else -%} {%- set first_user_prefix = messages[0]['content'][0]['text'] + '\n' -%} {%- endif -%} {%- set loop_messages = messages[1:] -%} {%- else -%} {%- set first_user_prefix = "" -%} {%- set loop_messages = messages -%} {%- endif -%} {%- for message in loop_messages -%} {%- if (message['role'] == 'user') != (loop.index0 % 2 == 0) -%} {{ raise_exception("Conversation roles must alternate user/assistant/user/assistant/...") }} {%- endif -%} {%- if (message['role'] == 'assistant') -%} {%- set role = "model" -%} {%- else -%} {%- set role = message['role'] -%} {%- endif -%} {{ '<start_of_turn>' + role + '\n' + (first_user_prefix if loop.first else "") }} {%- if message['content'] is string -%} {{ message['content'] | trim }} {%- elif message['content'] is iterable -%} {%- for item in message['content'] -%} {%- if item['type'] == 'image' -%} {{ '<start_of_image>' }} {%- elif item['type'] == 'text' -%} {{ item['text'] | trim }} {%- endif -%} {%- endfor -%} {%- else -%} {{ raise_exception("Invalid content type") }} {%- endif -%} {{ '<end_of_turn>\n' }} {%- endfor -%} {%- if add_generation_prompt -%} {{'<start_of_turn>model\n'}} {%- endif -%}
  gemma3.attention.head_count: 4
  gemma3.attention.head_count_kv: 1
  gemma3.embedding_length: 1152
  gemma3.attention.key_length: 256
  gemma3.attention.layer_norm_rms_epsilon: 0.000001
  gemma3.rope.scaling.type: linear
  gemma3.rope.freq_base: 1000000.000000
```

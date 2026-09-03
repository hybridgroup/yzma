# Roadmap

`yzma` currently has support for over 96% of `llama.cpp` functionality.

This is a list of all functions exposed by `llama.cpp` and the current state of
the associated `yzma` wrapper. The WebAssembly column gives the state of the same
function in a browser, which the [WebAssembly support](#webassembly-support)
section at the end explains.

## Completed wrappers

### Backend Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_backend_free` | yes | yes |
| `llama_backend_init` | yes | yes |
| `llama_flash_attn_type_name` | yes | no |
| `llama_ftype_name` | yes | no |
| `llama_load_mode_from_str` | yes | no |
| `llama_load_mode_name` | yes | no |
| `llama_max_devices` | yes | no |
| `llama_max_parallel_sequences` | yes | no |
| `llama_max_tensor_buft_overrides` | yes | no |
| `llama_numa_init` | yes | no |
| `llama_print_system_info` | yes | no |
| `llama_supports_gpu_offload` | yes | no |
| `llama_supports_mlock` | yes | no |
| `llama_supports_mmap` | yes | no |
| `llama_supports_rpc` | yes | no |
| `llama_time_us` | yes | no |

### Model Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_init_from_model` | yes | yes |
| `llama_model_chat_template` | yes | yes |
| `llama_model_cls_label` | yes | no |
| `llama_model_decoder_start_token` | yes | no |
| `llama_model_default_params` | yes | partial |
| `llama_model_desc` | yes | yes |
| `llama_model_free` | yes | yes |
| `llama_model_ftype` | yes | no |
| `llama_model_has_decoder` | yes | no |
| `llama_model_has_encoder` | yes | no |
| `llama_model_is_diffusion` | yes | no |
| `llama_model_is_hybrid` | yes | no |
| `llama_model_is_recurrent` | yes | no |
| `llama_model_load_from_file` | yes | yes |
| `llama_model_load_from_splits` | yes | no |
| `llama_model_meta_count` | yes | no |
| `llama_model_meta_key_by_index` | yes | no |
| `llama_model_meta_key_str` | yes | no |
| `llama_model_meta_val_str_by_index` | yes | no |
| `llama_model_meta_val_str` | yes | no |
| `llama_model_n_cls_out` | yes | no |
| `llama_model_n_ctx_train` | yes | yes |
| `llama_model_n_embd_inp` | yes | no |
| `llama_model_n_embd_out` | yes | no |
| `llama_model_n_embd` | yes | yes |
| `llama_model_n_head_kv` | yes | no |
| `llama_model_n_head` | yes | no |
| `llama_model_n_layer_nextn` | yes | no |
| `llama_model_n_layer` | yes | no |
| `llama_model_n_params` | yes | no |
| `llama_model_n_swa` | yes | no |
| `llama_model_quantize_default_params` | yes | no |
| `llama_model_quantize` | yes | no |
| `llama_model_rope_freq_scale_train` | yes | no |
| `llama_model_rope_type` | yes | no |
| `llama_model_save_to_file` | yes | no |
| `llama_model_size` | yes | no |
| `llama_split_path` | yes | no |
| `llama_split_prefix` | yes | no |

### Vocab Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_detokenize` | yes | partial |
| `llama_model_get_vocab` | yes | yes |
| `llama_token_to_piece` | yes | yes |
| `llama_tokenize` | yes | yes |
| `llama_vocab_bos` | yes | yes |
| `llama_vocab_eos` | yes | yes |
| `llama_vocab_eot` | yes | yes |
| `llama_vocab_fim_mid` | yes | yes |
| `llama_vocab_fim_pad` | yes | yes |
| `llama_vocab_fim_pre` | yes | yes |
| `llama_vocab_fim_rep` | yes | yes |
| `llama_vocab_fim_sep` | yes | yes |
| `llama_vocab_fim_suf` | yes | yes |
| `llama_vocab_get_add_bos` | yes | yes |
| `llama_vocab_get_add_eos` | yes | yes |
| `llama_vocab_get_add_sep` | yes | yes |
| `llama_vocab_get_attr` | yes | yes |
| `llama_vocab_get_score` | yes | yes |
| `llama_vocab_get_suppress_tokens` | yes | yes |
| `llama_vocab_get_text` | yes | yes |
| `llama_vocab_is_control` | yes | yes |
| `llama_vocab_is_eog` | yes | yes |
| `llama_vocab_mask` | yes | yes |
| `llama_vocab_n_tokens` | yes | yes |
| `llama_vocab_nl` | yes | yes |
| `llama_vocab_pad` | yes | yes |
| `llama_vocab_sep` | yes | yes |
| `llama_vocab_type` | yes | yes |

### Context Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_attach_threadpool` | yes | no |
| `llama_context_default_params` | yes | partial |
| `llama_decode` | yes | yes |
| `llama_detach_threadpool` | yes | no |
| `llama_encode` | yes | yes |
| `llama_free` | yes | yes |
| `llama_get_embeddings_ith` | yes | no |
| `llama_get_embeddings_seq` | yes | yes |
| `llama_get_embeddings` | yes | no |
| `llama_get_logits_ith` | yes | no |
| `llama_get_logits` | yes | no |
| `llama_get_memory` | yes | no |
| `llama_get_model` | yes | no |
| `llama_n_batch` | yes | no |
| `llama_n_ctx_seq` | yes | no |
| `llama_n_ctx` | yes | yes |
| `llama_n_rs_seq` | yes | no |
| `llama_n_seq_max` | yes | no |
| `llama_n_threads_batch` | yes | no |
| `llama_n_threads` | yes | no |
| `llama_n_ubatch` | yes | no |
| `llama_pooling_type` | yes | no |
| `llama_set_abort_callback` | yes | no |
| `llama_set_adapter_cvec` | yes | no |
| `llama_set_causal_attn` | yes | no |
| `llama_set_embeddings` | yes | no |
| `llama_set_n_threads` | yes | no |
| `llama_set_warmup` | yes | no |
| `llama_synchronize` | yes | no |

### Backend Sampling Functions (Experimental)

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_get_sampled_candidates_count_ith` | yes | no |
| `llama_get_sampled_candidates_ith` | yes | no |
| `llama_get_sampled_logits_count_ith` | yes | no |
| `llama_get_sampled_logits_ith` | yes | no |
| `llama_get_sampled_probs_count_ith` | yes | no |
| `llama_get_sampled_probs_ith` | yes | no |
| `llama_get_sampled_token_ith` | yes | no |

### Memory Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_memory_can_shift` | yes | no |
| `llama_memory_clear` | yes | yes |
| `llama_memory_seq_add` | yes | no |
| `llama_memory_seq_cp` | yes | no |
| `llama_memory_seq_div` | yes | no |
| `llama_memory_seq_keep` | yes | no |
| `llama_memory_seq_pos_max` | yes | no |
| `llama_memory_seq_pos_min` | yes | no |
| `llama_memory_seq_rm` | yes | no |

### Batch Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_batch_free` | yes | no |
| `llama_batch_get_one` | yes | yes |
| `llama_batch_init` | yes | no |

### Sampling Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_sampler_accept` | yes | yes |
| `llama_sampler_apply` | yes | no |
| `llama_sampler_chain_add` | yes | yes |
| `llama_sampler_chain_default_params` | yes | yes |
| `llama_sampler_chain_get` | yes | yes |
| `llama_sampler_chain_init` | yes | yes |
| `llama_sampler_chain_n` | yes | yes |
| `llama_sampler_chain_remove` | yes | yes |
| `llama_sampler_clone` | yes | yes |
| `llama_sampler_free` | yes | yes |
| `llama_sampler_get_seed` | yes | yes |
| `llama_sampler_init_adaptive_p` | yes | yes |
| `llama_sampler_init_dist` | yes | yes |
| `llama_sampler_init_dry` | yes | yes |
| `llama_sampler_init_grammar_lazy_patterns` | yes | yes |
| `llama_sampler_init_grammar` | yes | yes |
| `llama_sampler_init_greedy` | yes | yes |
| `llama_sampler_init_infill` | yes | yes |
| `llama_sampler_init_logit_bias` | yes | partial |
| `llama_sampler_init_min_p` | yes | yes |
| `llama_sampler_init_mirostat_v2` | yes | yes |
| `llama_sampler_init_mirostat` | yes | yes |
| `llama_sampler_init_penalties` | yes | yes |
| `llama_sampler_init_temp` | yes | yes |
| `llama_sampler_init_temp_ext` | yes | yes |
| `llama_sampler_init_top_k` | yes | yes |
| `llama_sampler_init_top_n_sigma` | yes | yes |
| `llama_sampler_init_top_p` | yes | yes |
| `llama_sampler_init_typical` | yes | yes |
| `llama_sampler_init_xtc` | yes | yes |
| `llama_sampler_name` | yes | yes |
| `llama_sampler_reset` | yes | yes |
| `llama_sampler_sample` | yes | yes |
| `llama_set_sampler` | yes | no |

### Logging Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_log_get` | yes | no |
| `llama_log_set` | yes | partial |

### Performance Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_perf_context` | yes | no |
| `llama_perf_context_print` | yes | no |
| `llama_perf_context_reset` | yes | no |
| `llama_perf_sampler` | yes | no |
| `llama_perf_sampler_print` | yes | no |
| `llama_perf_sampler_reset` | yes | no |

### Chat Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_chat_apply_template` | yes | partial |
| `llama_chat_builtin_templates` | yes | no |

### State Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_state_get_data` | yes | no |
| `llama_state_get_size` | yes | no |
| `llama_state_load_file` | yes | no |
| `llama_state_save_file` | yes | no |
| `llama_state_seq_get_data_ext` | yes | no |
| `llama_state_seq_get_data` | yes | no |
| `llama_state_seq_get_size_ext` | yes | no |
| `llama_state_seq_get_size` | yes | no |
| `llama_state_seq_load_file` | yes | no |
| `llama_state_seq_save_file` | yes | no |
| `llama_state_seq_set_data_ext` | yes | no |
| `llama_state_seq_set_data` | yes | no |
| `llama_state_set_data` | yes | no |

### LoRA Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_adapter_get_alora_invocation_tokens` | yes | no |
| `llama_adapter_get_alora_n_invocation_tokens` | yes | no |
| `llama_adapter_lora_free` | yes | no |
| `llama_adapter_lora_init` | yes | no |
| `llama_adapter_meta_count` | yes | no |
| `llama_adapter_meta_key_by_index` | yes | no |
| `llama_adapter_meta_val_str_by_index` | yes | no |
| `llama_adapter_meta_val_str` | yes | no |
| `llama_set_adapters_lora` | yes | no |

### `mtmd` Functions

Note that these functions are considered by `llama.cpp` to be experimental, and are subject to change.

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `mtmd_batch_add_chunk` | yes | no |
| `mtmd_batch_encode` | yes | no |
| `mtmd_batch_free` | yes | no |
| `mtmd_batch_get_output_embd` | yes | no |
| `mtmd_batch_init` | yes | no |
| `mtmd_bitmap_free` | yes | yes |
| `mtmd_bitmap_get_data` | yes | no |
| `mtmd_bitmap_get_id` | yes | no |
| `mtmd_bitmap_get_n_bytes` | yes | no |
| `mtmd_bitmap_get_nx` | yes | no |
| `mtmd_bitmap_get_ny` | yes | no |
| `mtmd_bitmap_init_from_audio` | yes | no |
| `mtmd_bitmap_init` | yes | yes |
| `mtmd_bitmap_is_audio` | yes | no |
| `mtmd_bitmap_set_id` | yes | no |
| `mtmd_context_params_default` | yes | yes |
| `mtmd_decode_use_mrope` | yes | no |
| `mtmd_decode_use_non_causal` | yes | no |
| `mtmd_default_marker` | yes | no |
| `mtmd_encode_chunk` | yes | no |
| `mtmd_encode` | yes | no |
| `mtmd_free` | yes | yes |
| `mtmd_get_audio_sample_rate` | yes | no |
| `mtmd_get_marker` | yes | yes |
| `mtmd_get_output_embd` | yes | no |
| `mtmd_helper_bitmap_init_from_buf` | yes | no |
| `mtmd_helper_bitmap_init_from_file` | yes | no |
| `mtmd_helper_eval_chunks` | yes | yes |
| `mtmd_helper_support_video` | yes | no |
| `mtmd_helper_video_free` | yes | no |
| `mtmd_helper_video_get_info` | yes | no |
| `mtmd_helper_video_init_from_buf` | yes | no |
| `mtmd_helper_video_init_params_default` | yes | no |
| `mtmd_helper_video_init` | yes | no |
| `mtmd_image_tokens_get_id` | yes | no |
| `mtmd_image_tokens_get_n_pos` | yes | no |
| `mtmd_image_tokens_get_n_tokens` | yes | no |
| `mtmd_image_tokens_get_nx` | yes | no |
| `mtmd_image_tokens_get_ny` | yes | no |
| `mtmd_init_from_file` | yes | yes |
| `mtmd_input_chunk_copy` | yes | no |
| `mtmd_input_chunk_free` | yes | no |
| `mtmd_input_chunk_get_id` | yes | no |
| `mtmd_input_chunk_get_n_pos` | yes | no |
| `mtmd_input_chunk_get_n_tokens` | yes | no |
| `mtmd_input_chunk_get_tokens_image` | yes | no |
| `mtmd_input_chunk_get_tokens_text` | yes | no |
| `mtmd_input_chunk_get_type` | yes | no |
| `mtmd_input_chunks_free` | yes | yes |
| `mtmd_input_chunks_get` | yes | no |
| `mtmd_input_chunks_init` | yes | yes |
| `mtmd_input_chunks_size` | yes | yes |
| `mtmd_log_set` | yes | no |
| `mtmd_support_audio` | yes | yes |
| `mtmd_support_vision` | yes | yes |
| `mtmd_tokenize` | yes | yes |
| `mtmd_tokenize_from_parts` | yes | no |

---

## Functions in `llama.cpp` still needing wrappers

### `llama` Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `llama_model_init_from_user` | no | no |
| `llama_model_load_from_file_ptr` | no | no |
| `llama_opt_epoch` | no | no |
| `llama_opt_init` | no | no |
| `llama_opt_param_filter_all` | no | no |
| `llama_sampler_init` | no | no |

### `mtmd` Functions

| Function | `yzma` | WebAssembly |
| --- | :-: | :-: |
| `mtmd_bitmap_init_lazy` | no | no |
| `mtmd_get_cap_from_file` | no | no |
| `mtmd_helper_video_read_next` | no | no |
| `mtmd_image_tokens_get_decoder_pos` | no | no |

---

## WebAssembly support

The `pkg/llamawasm` package drives a build of `llama.cpp` for WebAssembly. It
does not use the C API directly. A C shim in the
[llama-cpp-builder](https://github.com/hybridgroup/llama-cpp-builder) repository
gives it a small set of calls with a version, which is ABI 5 now. Thus a function
of `llama.cpp` reaches the browser only after the shim exports it.

The WebAssembly column of each table above gives the state of the wrapper.

| Value | Meaning |
| --- | --- |
| yes | The WebAssembly build has this function. |
| partial | The WebAssembly build has it with fewer options. See the notes. |
| no | The shim does not export it. |

93 functions reach WebAssembly, 87 complete and 6 partial, all of them among the
253 that have a wrapper on a host. That is sufficient for text generation,
embeddings, images, chat templates, tool calling with a grammar, and every
sampler that a host has.

### Notes on the partial wrappers

- `llama_detokenize` has no call in the shim. The Go side builds the text from
  `llama_token_to_piece` of each token, and it ignores `removeSpecial`.
- `llama_chat_apply_template` takes one message, which is sufficient for a
  question about an image. Render the template that `llama_model_chat_template`
  gives with the `template` package for a conversation with turns.
- `llama_log_set` cannot take a Go function as a callback. The shim holds the
  callback and the Go side only sets how much llama.cpp prints.
- `llama_sampler_init_logit_bias` takes two slices, the tokens and the biases,
  and not a pointer to an array of `LogitBias`. A struct cannot cross the
  boundary of the module.
- `llama_model_default_params` has `NGpuLayers` only.
- `llama_context_default_params` has `NCtx`, `NBatch`, `NUbatch`, `NThreads`,
  `PoolingType`, and `Embeddings`.

### What WebAssembly still needs

In order of the value that each one adds.

1. **Batches with positions.** `llama_batch_init` and `llama_batch_free`. Only
   `llama_batch_get_one` is available, thus a program cannot put more than one
   sequence in a batch.
2. **The metadata of a model.** The `llama_model_meta_*` calls and the counts of
   layers, heads, and parameters.
3. **Memory with sequences.** The `llama_memory_seq_*` calls and
   `llama_memory_can_shift`. Only `llama_memory_clear` is available.
4. **The parts of `mtmd`.** The shim does the whole pipeline of an image in two
   coarse calls, `mtmd_tokenize` and `mtmd_helper_eval_chunks`. Thus the calls
   that build or examine one piece are absent: the getters of a bitmap, the
   accessors of a chunk and of the tokens of an image, `mtmd_encode`,
   `mtmd_get_output_embd`, and the batch calls. A program in a browser cannot
   place the embeddings of an image itself.

`llama_sampler_apply` is not planned. It takes a `llama_token_data_array`, and no
struct crosses the boundary of the module. A program in a browser cannot reach
the logits in any case, because `llama_get_logits_ith` is absent as well.

Audio, video, LoRA adapters, saved state, quantization, and the performance
counters are not planned for WebAssembly.

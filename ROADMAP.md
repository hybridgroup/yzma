# Roadmap

This is a list of all functions exposed by `llama.cpp` and the current state of the associated `yzma` wrapper.

## Functions in `llama.cpp` with wrappers

### Backend Functions
- [x] `llama_backend_free`
- [x] `llama_backend_init`

### Model Functions
- [x] `llama_init_from_model`
- [x] `llama_model_chat_template`
- [x] `llama_model_cls_label`
- [x] `llama_model_decoder_start_token`
- [x] `llama_model_default_params`
- [x] `llama_model_desc`
- [x] `llama_model_free`
- [x] `llama_model_has_decoder`
- [x] `llama_model_has_encoder`
- [x] `llama_model_is_diffusion`
- [x] `llama_model_is_hybrid`
- [x] `llama_model_is_recurrent`
- [x] `llama_model_load_from_file`
- [x] `llama_model_n_cls_out`
- [x] `llama_model_n_ctx_train`
- [x] `llama_model_n_embd`
- [x] `llama_model_n_head_kv`
- [x] `llama_model_n_head`
- [x] `llama_model_n_layer`
- [x] `llama_model_n_swa`
- [x] `llama_model_rope_freq_scale_train`
- [x] `llama_model_rope_type`
- [x] `llama_model_size`

### Vocab Functions
- [x] `llama_model_get_vocab`
- [x] `llama_token_to_piece`
- [x] `llama_tokenize`
- [x] `llama_vocab_bos`
- [x] `llama_vocab_eos`
- [x] `llama_vocab_eot`
- [x] `llama_vocab_fim_mid`
- [x] `llama_vocab_fim_pad`
- [x] `llama_vocab_fim_pre`
- [x] `llama_vocab_fim_rep`
- [x] `llama_vocab_fim_sep`
- [x] `llama_vocab_fim_suf`
- [x] `llama_vocab_get_add_bos`
- [x] `llama_vocab_get_add_eos`
- [x] `llama_vocab_get_add_sep`
- [x] `llama_vocab_get_attr`
- [x] `llama_vocab_get_score`
- [x] `llama_vocab_get_text`
- [x] `llama_vocab_is_control`
- [x] `llama_vocab_is_eog`
- [x] `llama_vocab_mask`
- [x] `llama_vocab_n_tokens`
- [x] `llama_vocab_nl`
- [x] `llama_vocab_pad`
- [x] `llama_vocab_sep`
- [x] `llama_vocab_type`

### Context Functions
- [x] `llama_context_default_params`
- [x] `llama_decode`
- [x] `llama_encode`
- [x] `llama_free`
- [x] `llama_get_embeddings_ith`
- [x] `llama_get_embeddings_seq`
- [x] `llama_get_memory`
- [x] `llama_perf_context_reset`
- [x] `llama_pooling_type`
- [x] `llama_set_warmup`
- [x] `llama_synchronize`

### Memory Functions
- [x] `llama_memory_can_shift`
- [x] `llama_memory_clear`
- [x] `llama_memory_seq_add`
- [x] `llama_memory_seq_cp`
- [x] `llama_memory_seq_div`
- [x] `llama_memory_seq_keep`
- [x] `llama_memory_seq_pos_max`
- [x] `llama_memory_seq_pos_min`
- [x] `llama_memory_seq_rm`

### Batch Functions
- [x] `llama_batch_free`
- [x] `llama_batch_get_one`
- [x] `llama_batch_init`

### Sampling Functions
- [x] `llama_sampler_accept`
- [x] `llama_sampler_chain_add`
- [x] `llama_sampler_chain_default_params`
- [x] `llama_sampler_chain_init`
- [x] `llama_sampler_free`
- [x] `llama_sampler_init_dist`
- [x] `llama_sampler_init_dry`
- [x] `llama_sampler_init_grammar`
- [x] `llama_sampler_init_greedy`
- [x] `llama_sampler_init_logit_bias`
- [x] `llama_sampler_init_min_p`
- [x] `llama_sampler_init_penalties`
- [x] `llama_sampler_init_temp_ext`
- [x] `llama_sampler_init_top_k`
- [x] `llama_sampler_init_top_n_sigma`
- [x] `llama_sampler_init_top_p`
- [x] `llama_sampler_init_typical`
- [x] `llama_sampler_init_xtc`
- [x] `llama_sampler_sample`

### Logging Functions
- [x] `llama_log_set`

### Chat Functions
- [x] `llama_chat_apply_template`

---

## Functions in `llama.cpp` still needing wrappers

- [ ] `llama_adapter_get_alora_invocation_tokens`
- [ ] `llama_adapter_get_alora_n_invocation_tokens`
- [ ] `llama_adapter_lora_free`
- [ ] `llama_adapter_lora_init`
- [ ] `llama_adapter_meta_count`
- [ ] `llama_adapter_meta_key_by_index`
- [ ] `llama_adapter_meta_val_str_by_index`
- [ ] `llama_adapter_meta_val_str`
- [ ] `llama_apply_adapter_cvec`
- [ ] `llama_attach_threadpool`
- [ ] `llama_clear_adapter_lora`
- [ ] `llama_detach_threadpool`
- [ ] `llama_get_embeddings`
- [ ] `llama_get_logits_ith`
- [ ] `llama_get_logits`
- [ ] `llama_get_model`
- [ ] `llama_max_devices`
- [ ] `llama_max_parallel_sequences`
- [ ] `llama_model_load_from_splits`
- [ ] `llama_model_quantize_default_params`
- [ ] `llama_model_quantize`
- [ ] `llama_n_batch`
- [ ] `llama_n_ctx`
- [ ] `llama_n_seq_max`
- [ ] `llama_n_ubatch`
- [ ] `llama_numa_init`
- [ ] `llama_opt_epoch`
- [ ] `llama_opt_init`
- [ ] `llama_opt_param_filter_all`
- [ ] `llama_rm_adapter_lora`
- [ ] `llama_set_adapter_lora`
- [ ] `llama_set_causal_attn`
- [ ] `llama_set_embeddings`
- [ ] `llama_state_get_data`
- [ ] `llama_state_get_size`
- [ ] `llama_state_load_file`
- [ ] `llama_state_save_file`
- [ ] `llama_state_seq_get_data_ext`
- [ ] `llama_state_seq_get_data`
- [ ] `llama_state_seq_get_size_ext`
- [ ] `llama_state_seq_get_size`
- [ ] `llama_state_seq_load_file`
- [ ] `llama_state_seq_save_file`
- [ ] `llama_state_seq_set_data_ext`
- [ ] `llama_state_seq_set_data`
- [ ] `llama_state_set_data`
- [ ] `llama_supports_gpu_offload`
- [ ] `llama_supports_mlock`
- [ ] `llama_supports_mmap`
- [ ] `llama_supports_rpc`
- [ ] `llama_time_us`

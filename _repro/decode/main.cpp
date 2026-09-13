// Mirrors _repro/batchcompare in C++, so a report to llama.cpp needs no Go.
// It makes a context with the library defaults, tokenizes two tokens and calls
// llama_decode one time. There is no warmup and no graph reuse.
#include "llama.h"

#include <cstdio>
#include <cstring>
#include <string>
#include <vector>

int main(int argc, char ** argv) {
    std::string model_file;
    std::string mode = "getone";

    for (int i = 1; i < argc; i++) {
        if (strcmp(argv[i], "-model") == 0 && i + 1 < argc) {
            model_file = argv[++i];
        } else if (strcmp(argv[i], "-mode") == 0 && i + 1 < argc) {
            mode = argv[++i];
        }
    }

    if (model_file.empty()) {
        fprintf(stderr, "error: -model is necessary\n");
        return 1;
    }

    llama_backend_init();

    llama_model * model = llama_model_load_from_file(model_file.c_str(), llama_model_default_params());
    if (model == nullptr) {
        fprintf(stderr, "error: unable to load model from file %s\n", model_file.c_str());
        return 1;
    }

    llama_context * ctx = llama_init_from_model(model, llama_context_default_params());
    if (ctx == nullptr) {
        fprintf(stderr, "error: unable to initialize context from model\n");
        return 1;
    }

    const llama_vocab * vocab = llama_model_get_vocab(model);

    const char * text = "Hello world";
    std::vector<llama_token> tokens(64);
    const int32_t n = llama_tokenize(vocab, text, strlen(text), tokens.data(), (int32_t) tokens.size(), true, true);
    if (n < 0) {
        fprintf(stderr, "error: unable to tokenize\n");
        return 1;
    }
    tokens.resize(n);

    llama_batch batch;
    if (mode == "getone") {
        batch = llama_batch_get_one(tokens.data(), n);
    } else if (mode == "full") {
        // llama_batch_init leaves the members uninitialized, so every one of
        // them gets a value, which is what the Go Add does.
        batch = llama_batch_init(n, 0, 1);
        batch.n_tokens = n;
        for (int32_t i = 0; i < n; i++) {
            batch.token[i]     = tokens[i];
            batch.pos[i]       = i;
            batch.n_seq_id[i]  = 1;
            batch.seq_id[i][0] = 0;
            batch.logits[i]    = i == n - 1;
        }
    } else {
        fprintf(stderr, "error: unknown mode %s\n", mode.c_str());
        return 1;
    }

    const int32_t ret = llama_decode(ctx, batch);
    if (ret != 0) {
        fprintf(stderr, "error: decode returned %d\n", ret);
        return 1;
    }

    printf("mode %s decoded %d tokens\n", mode.c_str(), n);

    if (mode == "full") {
        llama_batch_free(batch);
    }
    llama_free(ctx);
    llama_model_free(model);
    llama_backend_free();

    return 0;
}

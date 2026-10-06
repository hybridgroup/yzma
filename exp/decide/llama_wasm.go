//go:build js && wasm

package decide

import (
	"errors"
	"fmt"
	"strings"

	llama "github.com/hybridgroup/yzma/pkg/llamawasm"
)

type (
	token = llama.Token
	seqID = llama.SeqId
)

// backend holds the llama.cpp objects of a Decider, here through pkg/llamawasm.
// That package's memory calls take the context.
type backend struct {
	model    llama.Model
	vocab    llama.Vocab
	ctx      llama.Context
	nVocab   int32
	nEmbdOut int32
}

// ctxParams is what a Decider sets on its context.
type ctxParams struct {
	nCtx, nBatch, nUbatch, nSeqMax, nOutputsMax uint32
	threads                                     int32
	embeddings                                  bool
}

// load reads a model that is already in the module's filesystem, for example
// from [llama.FetchModelFile]. It uses the GPU when there is one.
func (b *backend) load(modelPath string) error {
	params := llama.ModelDefaultParams()
	if llama.GPUDevice() != "" {
		params.NGpuLayers = 999
	}

	model, err := llama.ModelLoadFromFile(modelPath, params)
	if err != nil {
		return fmt.Errorf("decide: load model: %w", err)
	}
	if model == 0 {
		return fmt.Errorf("decide: unable to load model %s", modelPath)
	}

	b.model = model
	b.vocab = llama.ModelGetVocab(model)
	b.nEmbdOut = llama.ModelNEmbdOut(model)
	b.nVocab = llama.VocabNTokens(b.vocab)
	return nil
}

// newContext creates the context and returns how many sequences it holds. A
// module before ABI version 10 has no unified cache, so the context then holds
// one sequence and each question is decoded separately.
func (b *backend) newContext(p ctxParams) (uint32, error) {
	params := llama.ContextDefaultParams()
	params.NCtx = p.nCtx
	params.NBatch = p.nBatch
	params.NUbatch = p.nUbatch
	params.NSeqMax = p.nSeqMax
	params.NOutputsMax = p.nOutputsMax
	params.KVUnified = 1
	params.NoPerf = 1
	if p.embeddings {
		params.Embeddings = 1
		params.PoolingType = llama.PoolingTypeNone
	}
	if p.threads > 0 {
		params.NThreads = p.threads
	}

	ctx, err := llama.InitFromModel(b.model, params)
	if errors.Is(err, llama.ErrNoKVUnified) || errors.Is(err, llama.ErrNoBatch) {
		params.NSeqMax = 1
		params.KVUnified = 0
		ctx, err = llama.InitFromModel(b.model, params)
	}
	if err != nil || ctx == 0 {
		return 0, fmt.Errorf("decide: unable to create context: %v", err)
	}
	b.ctx = ctx

	return max(1, params.NSeqMax), nil
}

func (b *backend) open() bool {
	return b.ctx != 0
}

func (b *backend) close() {
	if b.ctx != 0 {
		llama.Free(b.ctx)
		b.ctx = 0
	}
	if b.model != 0 {
		llama.ModelFree(b.model)
		b.model = 0
	}
}

func (b *backend) tokenize(s string, parseSpecial bool) []token {
	return llama.Tokenize(b.vocab, s, false, parseSpecial)
}

func (b *backend) memClear() {
	llama.MemoryClear(b.ctx, true)
}

func (b *backend) memSeqRm(seq seqID) {
	llama.MemorySeqRm(b.ctx, seq, -1, -1)
}

func (b *backend) memSeqCp(src, dst seqID) {
	llama.MemorySeqCp(b.ctx, src, dst, -1, -1)
}

type batch struct {
	llama.Batch
}

func newBatch(n int) *batch {
	return &batch{llama.BatchInit(int32(n), 0, 1)}
}

func (bt *batch) add(tok token, pos int, seq seqID, logits bool) error {
	return bt.Add(tok, llama.Pos(pos), []llama.SeqId{seq}, logits)
}

func (bt *batch) len() int32 {
	return bt.NTokens
}

func (bt *batch) free() {
	llama.BatchFree(bt.Batch)
}

func (b *backend) decode(bt *batch) error {
	ret, err := llama.Decode(b.ctx, bt.Batch)
	if err != nil {
		return err
	}
	if ret != 0 {
		return fmt.Errorf("decide: decode returned %d", ret)
	}
	return nil
}

// logits returns the output logits at batch index i.
func (b *backend) logits(i int32) ([]float32, error) {
	return llama.GetLogitsIth(b.ctx, i, b.nVocab)
}

// meta returns the metadata value of key, and false if the model has none.
func (b *backend) meta(key string) (string, bool) {
	return llama.ModelMetaValStr(b.model, key)
}

// metaPrefix returns the metadata values whose key starts with prefix, by the rest of the key.
func (b *backend) metaPrefix(prefix string) map[string]string {
	out := map[string]string{}
	for i := range llama.ModelMetaCount(b.model) {
		k, ok := llama.ModelMetaKeyByIndex(b.model, i)
		if !ok || !strings.HasPrefix(k, prefix) {
			continue
		}
		if v, ok := llama.ModelMetaValStrByIndex(b.model, i); ok {
			out[k[len(prefix):]] = v
		}
	}
	return out
}

func (b *backend) chatTemplate(name string) string {
	return llama.ModelChatTemplate(b.model, name)
}

// embeddings returns the embeddings output at batch index i.
func (b *backend) embeddings(i int32) ([]float32, error) {
	return llama.GetEmbeddingsIth(b.ctx, i, b.nEmbdOut)
}

func (b *backend) vocabMask() token {
	return llama.VocabMASK(b.vocab)
}

func (b *backend) vocabSep() token {
	return llama.VocabSEP(b.vocab)
}

// piece returns the text of tok, with special tokens written out.
func (b *backend) piece(tok token) string {
	buf := make([]byte, 256)
	n := llama.TokenToPiece(b.vocab, tok, buf, 0, true)
	if n <= 0 {
		return ""
	}
	return string(buf[:min(int(n), len(buf))])
}

package vlm

import (
	"errors"
	"fmt"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
)

// VLM is a Vision Language Model (VLM).
type VLM struct {
	ModelFilename     string
	ProjectorFilename string

	Model llama.Model

	ModelParams     llama.ModelParams
	ContextParams   llama.ContextParams
	ProjectorParams mtmd.ContextParamsType

	Sampler          llama.Sampler
	ModelContext     llama.Context
	ProjectorContext mtmd.Context

	template string
}

// NewVLM creates a new VLM.
func NewVLM(model, projector string) *VLM {
	return &VLM{
		ModelFilename:     model,
		ProjectorFilename: projector,
	}
}

// Close closes the VLM.
func (m *VLM) Close() {
	if m.ProjectorContext != 0 {
		mtmd.Free(m.ProjectorContext)

	}

	if m.ModelContext != 0 {
		llama.Free(m.ModelContext)
	}

	if m.Sampler != 0 {
		llama.SamplerFree(m.Sampler)
	}

	if m.Model != 0 {
		llama.ModelFree(m.Model)
	}
}

// Init initializes the VLM.
func (m *VLM) Init() error {
	var err error
	if m.ModelParams == (llama.ModelParams{}) {
		m.ModelParams = llama.ModelDefaultParams()
	}
	m.Model, err = llama.ModelLoadFromFile(m.ModelFilename, m.ModelParams)
	if err != nil {
		return fmt.Errorf("unable to load model: %w", err)
	}

	if m.ContextParams == (llama.ContextParams{}) {
		m.ContextParams = llama.ContextDefaultParams()
		m.ContextParams.NCtx = 8192
		m.ContextParams.NBatch = 512
	}

	m.ModelContext, err = llama.InitFromModel(m.Model, m.ContextParams)
	if err != nil {
		return fmt.Errorf("unable to initialize model context: %w", err)
	}

	m.template = llama.ModelChatTemplate(m.Model, "")

	if m.Sampler == 0 {
		sp := llama.DefaultSamplerParams()
		sp.Temp = float32(0.8)
		sp.TopK = int32(40)
		sp.TopP = float32(0.9)
		sp.MinP = float32(0.1)

		m.Sampler = llama.NewSampler(m.Model, llama.DefaultSamplers, sp)
	}

	if m.ProjectorParams == (mtmd.ContextParamsType{}) {
		m.ProjectorParams = mtmd.ContextParamsDefault()
	}
	m.ProjectorContext, err = mtmd.InitFromFile(m.ProjectorFilename, m.Model, m.ProjectorParams)
	if err != nil {
		return fmt.Errorf("unable to initialize projector context: %w", err)
	}

	return nil
}

// ChatTemplate applies the model's chat template to the given messages.
// It returns an empty string when llama.cpp does not know the template.
func (m *VLM) ChatTemplate(messages []llama.ChatMessage, add bool) string {
	buf := make([]byte, 16384)
	n := llama.ChatApplyTemplate(m.template, messages, add, buf)
	if int(n) > len(buf) {
		// The result is the full length even when buf is too short.
		buf = make([]byte, n)
		n = llama.ChatApplyTemplate(m.template, messages, add, buf)
	}
	if n < 0 || int(n) > len(buf) {
		return ""
	}

	return string(buf[:n])
}

// Tokenize tokenizes the input text and image bitmap into output chunks.
func (m *VLM) Tokenize(input *mtmd.InputText, bitmaps []mtmd.Bitmap, output mtmd.InputChunks) (err error) {
	if res := mtmd.Tokenize(m.ProjectorContext, output, input, bitmaps); res != 0 {
		err = fmt.Errorf("unable to tokenize: %d", res)
	}
	return
}

// Results generates text results from the given result chunks.
func (m *VLM) Results(chunks mtmd.InputChunks) (string, error) {
	var n llama.Pos
	nBatch := llama.NBatch(m.ModelContext)

	if res := mtmd.HelperEvalChunks(m.ProjectorContext, m.ModelContext, chunks, 0, 0, int32(nBatch), true, &n); res != 0 {
		return "", errors.New("unable to evaluate chunks")
	}

	vocab := llama.ModelGetVocab(m.Model)
	results := ""
	buf := make([]byte, 128)

	for i := 0; i < int(nBatch); i++ {
		token := llama.SamplerSample(m.Sampler, m.ModelContext, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}

		var piece string
		piece, buf = tokenPiece(vocab, token, buf)
		results += piece

		batch := llama.BatchGetOne([]llama.Token{token})
		batch.Pos = &n

		llama.Decode(m.ModelContext, batch)
		n++
	}

	m.Clear()

	return results, nil
}

// tokenPiece returns the text of token. It grows buf when the piece does not fit
// and returns the buffer to use next time.
func tokenPiece(vocab llama.Vocab, token llama.Token, buf []byte) (string, []byte) {
	n := llama.TokenToPiece(vocab, token, buf, 0, true)
	if n < 0 {
		// A negative result is the size the piece needs.
		buf = make([]byte, -n)
		n = llama.TokenToPiece(vocab, token, buf, 0, true)
	}
	if n <= 0 || int(n) > len(buf) {
		return "", buf
	}

	return string(buf[:n]), buf
}

// Clear clears the context memory.
func (m *VLM) Clear() {
	llama.Synchronize(m.ModelContext)
	mem, err := llama.GetMemory(m.ModelContext)
	if err != nil {
		fmt.Println("unable to get memory:", err)
		return
	}
	llama.MemoryClear(mem, true)
}

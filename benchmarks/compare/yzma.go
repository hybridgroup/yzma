package compare

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/ardanlabs/jinja"
	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
)

// YzmaOptions says which library, model and device the in process engine takes.
type YzmaOptions struct {
	Library   string
	Model     string
	Projector string // empty for the text suite
	Device    string // comma separated, as the other benchmarks take it
	NCtx      int
	NBatch    int

	// ImageMinTokens and ImageMaxTokens give the budget of tokens of an image.
	// They decide if the projector splits the image into tiles, thus they
	// decide how much work the engine does. Zero keeps the default.
	ImageMinTokens int32
	ImageMaxTokens int32

	// Threads is how many threads the vision model uses. The default of
	// llama.cpp is 4, while a server takes the count of the machine.
	Threads int32

	// Embeddings makes a context that gives vectors and not tokens. An
	// embedding model needs it, and a context of one kind cannot do the other.
	Embeddings bool
}

// yzmaEngine calls llama.cpp in the same process. It has no server and no
// socket, thus it is the baseline of the comparison.
type yzmaEngine struct {
	opts     YzmaOptions
	model    llama.Model
	ctx      llama.Context
	mctx     mtmd.Context
	template string

	// compiled holds the chat template of the model, ready to render. A server
	// compiles its template one time when it loads the model, thus yzma must
	// not compile it again for each request. The compiling costs 97 us, and the
	// rendering alone costs 5 us.
	compiled *jinja.Template
	sampler  llama.Sampler
	loaded   bool
}

// NewYzma gives the in process engine.
func NewYzma(opts YzmaOptions) Engine {
	return &yzmaEngine{opts: opts}
}

func (e *yzmaEngine) Name() string { return "yzma" }

func (e *yzmaEngine) Load(req Request) error {
	if e.loaded {
		return nil
	}

	if err := llama.Load(e.opts.Library); err != nil {
		return fmt.Errorf("unable to load the library: %w", err)
	}
	if e.opts.Projector != "" {
		if err := mtmd.Load(e.opts.Library); err != nil {
			return fmt.Errorf("unable to load the library: %w", err)
		}
		mtmd.LogSet(llama.LogSilent())
	}

	llama.LogSet(llama.LogSilent())
	llama.Init()

	mparams := llama.ModelDefaultParams()
	mparams.LoadMode = llama.LoadModeNone

	if e.opts.Device != "" {
		devs := []llama.GGMLBackendDevice{}
		for name := range strings.SplitSeq(e.opts.Device, ",") {
			dev := llama.GGMLBackendDeviceByName(name)
			if dev == 0 {
				return fmt.Errorf("unknown device: %s", name)
			}
			devs = append(devs, dev)
		}

		devs = append(devs, 0) // NULL terminator required by llama.cpp
		if err := mparams.SetDevices(devs); err != nil {
			return err
		}
		defer runtime.KeepAlive(devs)
	}

	model, err := llama.ModelLoadFromFile(e.opts.Model, mparams)
	if err != nil {
		return fmt.Errorf("unable to load the model: %w", err)
	}
	e.model = model

	params := llama.ContextDefaultParams()
	params.NCtx = uint32(e.opts.NCtx)
	params.NBatch = uint32(e.opts.NBatch)
	if e.opts.Embeddings {
		params.Embeddings = 1
		params.PoolingType = llama.PoolingTypeMean
	}

	ctx, err := llama.InitFromModel(model, params)
	if err != nil {
		return err
	}
	e.ctx = ctx

	e.template = llama.ModelChatTemplate(model, "")
	if t, err := jinja.Compile(e.template); err == nil {
		e.compiled = t
	}

	// One sampler for the whole run, as a server also keeps one.
	e.sampler = llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(e.sampler, llama.SamplerInitGreedy())

	if e.opts.Projector != "" {
		mprms := mtmd.ContextParamsDefault()
		if e.opts.ImageMinTokens > 0 {
			mprms.ImageMinTokens = e.opts.ImageMinTokens
		}
		if e.opts.ImageMaxTokens > 0 {
			mprms.ImageMaxTokens = e.opts.ImageMaxTokens
		}
		if e.opts.Threads > 0 {
			mprms.Threads = e.opts.Threads
		}

		mctx, err := mtmd.InitFromFile(e.opts.Projector, model, mprms)
		if err != nil {
			return fmt.Errorf("unable to load the projector: %w", err)
		}
		e.mctx = mctx
	}

	if len(req.Image) > 0 && e.mctx == 0 {
		return fmt.Errorf("the image suite needs a projector, give -mmproj")
	}

	e.loaded = true

	// One request before the timed runs, as the servers also get.
	if e.opts.Embeddings {
		_, err = e.Embed(req)

		return err
	}

	warm := req
	warm.MaxTokens = 1
	_, err = e.Generate(warm)

	return err
}

// Embed turns the prompt into a vector. An embedding model takes the plain
// text, thus there is no chat template here and no template on the servers.
func (e *yzmaEngine) Embed(req Request) (Result, error) {
	defer e.reset()

	start := time.Now()

	vocab := llama.ModelGetVocab(e.model)
	tokens := llama.Tokenize(vocab, req.Prompt, true, true)

	if _, err := llama.Decode(e.ctx, llama.BatchGetOne(tokens)); err != nil {
		return Result{}, fmt.Errorf("decode failed: %w", err)
	}

	vec, err := llama.GetEmbeddingsSeq(e.ctx, 0, llama.ModelNEmbd(e.model))
	if err != nil {
		return Result{}, fmt.Errorf("unable to get the embeddings: %w", err)
	}
	if len(vec) == 0 {
		return Result{}, fmt.Errorf("the vector is empty")
	}

	elapsed := time.Since(start)

	return Result{
		PromptTokens: len(tokens),
		Dimensions:   len(vec),
		FirstToken:   elapsed,
		Total:        elapsed,
	}, nil
}

func (e *yzmaEngine) Close() {
	if !e.loaded {
		return
	}

	if e.sampler != 0 {
		llama.SamplerFree(e.sampler)
	}
	if e.mctx != 0 {
		mtmd.Free(e.mctx)
		mtmd.LogSet(llama.LogNormal)
	}

	llama.Free(e.ctx)
	llama.ModelFree(e.model)
	llama.LogSet(llama.LogNormal)
	llama.Close()

	e.loaded = false
}

func (e *yzmaEngine) Generate(req Request) (Result, error) {
	if len(req.Image) > 0 {
		return e.generateWithImage(req)
	}

	return e.generateText(req)
}

func (e *yzmaEngine) generateText(req Request) (Result, error) {
	defer e.reset()

	start := time.Now()

	vocab := llama.ModelGetVocab(e.model)
	prompt := e.render(req.Prompt)

	tokens := llama.Tokenize(vocab, prompt, true, true)
	batch := llama.BatchGetOne(tokens)

	result := Result{PromptTokens: len(tokens)}
	for range req.MaxTokens {
		llama.Decode(e.ctx, batch)

		token := llama.SamplerSample(e.sampler, e.ctx, -1)
		if llama.VocabIsEOG(vocab, token) {
			break
		}
		if token == llama.TokenNull {
			return result, fmt.Errorf("SamplerSample gave TokenNull")
		}

		if result.Tokens == 0 {
			result.FirstToken = time.Since(start)
		}
		result.Tokens++
		result.Text += piece(vocab, token)

		batch = llama.BatchGetOne([]llama.Token{token})
	}

	result.Total = time.Since(start)
	if result.Tokens == 0 {
		return result, fmt.Errorf("the answer has no token")
	}

	return result, nil
}

func (e *yzmaEngine) generateWithImage(req Request) (Result, error) {
	defer e.reset()

	start := time.Now()

	vocab := llama.ModelGetVocab(e.model)
	prompt := e.render(mtmd.DefaultMarker() + req.Prompt)

	// The bitmap comes from the bytes of this run. The servers read the same
	// bytes out of the request, thus both pay for the decoding of the image.
	bitmap := mtmd.BitmapInitFromBuf(e.mctx, &req.Image[0], uint64(len(req.Image)), false, mtmd.InitOptDefault())
	if bitmap.Bitmap == 0 {
		return Result{}, fmt.Errorf("unable to read the image")
	}
	defer mtmd.BitmapFree(bitmap.Bitmap)

	chunks := mtmd.InputChunksInit()
	defer mtmd.InputChunksFree(chunks)

	input := mtmd.NewInputText(prompt, true, true)
	if res := mtmd.Tokenize(e.mctx, chunks, input, []mtmd.Bitmap{bitmap.Bitmap}); res != 0 {
		return Result{}, fmt.Errorf("Tokenize gave %d", res)
	}

	promptTokens := 0
	for i := uint64(0); i < mtmd.InputChunksSize(chunks); i++ {
		promptTokens += int(mtmd.InputChunkGetNTokens(mtmd.InputChunksGet(chunks, i)))
	}

	var n llama.Pos
	nBatch := llama.NBatch(e.ctx)
	if res := mtmd.HelperEvalChunks(e.mctx, e.ctx, chunks, 0, 0, int32(nBatch), true, &n); res != 0 {
		return Result{}, fmt.Errorf("HelperEvalChunks gave %d", res)
	}

	result := Result{PromptTokens: promptTokens}
	for range req.MaxTokens {
		token := llama.SamplerSample(e.sampler, e.ctx, -1)
		if llama.VocabIsEOG(vocab, token) {
			break
		}
		if token == llama.TokenNull {
			return result, fmt.Errorf("SamplerSample gave TokenNull")
		}

		if result.Tokens == 0 {
			result.FirstToken = time.Since(start)
		}
		result.Tokens++
		result.Text += piece(vocab, token)

		batch := llama.BatchGetOne([]llama.Token{token})
		batch.Pos = &n

		llama.Decode(e.ctx, batch)
		n++
	}

	result.Total = time.Since(start)
	if result.Tokens == 0 {
		return result, fmt.Errorf("the answer has no token")
	}

	return result, nil
}

// reset empties the KV cache, thus the next generation starts as this one did.
// The servers do the same between requests.
func (e *yzmaEngine) reset() {
	llama.Synchronize(e.ctx)

	mem, err := llama.GetMemory(e.ctx)
	if err != nil {
		return
	}
	llama.MemoryClear(mem, true)
}

// applyTemplate renders the chat template of the model. It takes the jinja
// path, because that is what the servers do. The built in template of
// llama.cpp is an approximation, and for Gemma 4 it gives 14 tokens where the
// servers give 22. Thinking is off, thus the engines answer the same way.
// render puts the text in the chat template of the model. It renders the
// template that Load compiled, and it does not think, thus the answer matches
// the servers, which also get no thinking.
func (e *yzmaEngine) render(text string) string {
	if e.compiled != nil {
		out, err := e.compiled.Render(map[string]any{
			"messages":              []any{map[string]any{"role": "user", "content": text}},
			"add_generation_prompt": true,
			"enable_thinking":       false,
		})
		if err == nil && out != "" {
			return out
		}
	}

	// A template that jinja cannot compile falls back to the built in one.
	buf := make([]byte, 4096)
	n := llama.ChatApplyTemplate(e.template, []llama.ChatMessage{llama.NewChatMessage("user", text)}, true, buf)
	if n <= 0 {
		return text
	}

	return string(buf[:n])
}

func piece(vocab llama.Vocab, token llama.Token) string {
	buf := make([]byte, 128)
	n := llama.TokenToPiece(vocab, token, buf, 0, true)
	if n <= 0 {
		return ""
	}

	return string(buf[:n])
}

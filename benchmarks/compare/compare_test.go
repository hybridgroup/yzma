package compare

import (
	"flag"
	"fmt"
	"image"
	"os"
	"testing"
	"time"
)

var (
	engineName  string
	suite       string
	library     string
	modelFile   string
	projector   string
	device      string
	serverModel string
	ollamaURL   string
	dmrURL      string
	nCtx        int
	nBatch      int
	maxTokens   int
	seed        uint
	imgMin      int
	imgMax      int
	imageFile   string
	imageSize   string
	threads     int

	engine Engine

	// images give a different image for each run, thus no engine answers
	// from its cache.
	images *Images

	// sent counts the requests of this process. Each -count of a benchmark
	// runs the loop again, thus a count of its own would repeat the requests.
	sent int
)

func init() {
	flag.StringVar(&engineName, "engine", "yzma", "engine to measure (yzma, ollama, dmr)")
	flag.StringVar(&suite, "suite", "text", "suite to run (text, multimodal, embeddings)")
	flag.StringVar(&library, "lib", os.Getenv("YZMA_LIB"), "directory of the llama.cpp library, for yzma")
	flag.StringVar(&modelFile, "model", "", "GGUF file of the model, for yzma")
	flag.StringVar(&projector, "mmproj", "", "GGUF file of the projector, for yzma with images")
	flag.StringVar(&device, "device", "", "comma separated devices, for yzma (CUDA0)")
	flag.StringVar(&serverModel, "server-model", "", "name of the model in the server")
	flag.StringVar(&ollamaURL, "ollama-url", "http://localhost:11434/v1", "OpenAI compatible address of ollama")
	flag.StringVar(&dmrURL, "dmr-url", "http://localhost:12434/engines/v1", "OpenAI compatible address of Docker Model Runner")
	flag.IntVar(&nCtx, "nctx", 8192, "context tokens, for yzma")
	flag.IntVar(&nBatch, "nbatch", 1024, "batch tokens, for yzma")
	flag.IntVar(&maxTokens, "tokens", 16, "tokens to make in one request")
	flag.UintVar(&seed, "seed", 1234, "seed of the sampler")
	flag.StringVar(&imageFile, "image", ImageFile, "image file of the multimodal suite")
	flag.StringVar(&imageSize, "image-size", "", "size to scale the image to, such as 1280x960, empty keeps the size")
	flag.IntVar(&threads, "threads", 0, "threads of the vision model, for yzma, 0 keeps the default")
	flag.IntVar(&imgMin, "image-min-tokens", 0, "least tokens of an image, for yzma, 0 keeps the default")
	flag.IntVar(&imgMax, "image-max-tokens", 0, "most tokens of an image, for yzma, 0 keeps the default")
}

func TestMain(m *testing.M) {
	flag.Parse()

	code := m.Run()

	if engine != nil {
		engine.Close()
	}

	os.Exit(code)
}

// request gives the request of run n. Each run gets a prompt and an image
// that no engine has seen, thus every engine does the whole work every time.
// A server keeps the prompts and images it has seen and answers a repeat of
// them almost at once, while yzma empties its cache after each generation.
func request(n int) (Request, error) {
	req := Request{
		Prompt:    TextVariant(n),
		MaxTokens: maxTokens,
		Seed:      uint32(seed),
	}

	switch suite {
	case "multimodal":
		req.Prompt = ImagePrompt
		img, err := images.Variant(n)
		if err != nil {
			return req, err
		}
		req.Image = img
	case "embeddings":
		req.Prompt = EmbedVariant(n)
		req.Embeddings = true
	}

	return req, nil
}

// next gives a request that this process has not sent yet. Request 0 is the
// one of the load, thus the first one of next is 1.
func next(t skipper) Request {
	sent++
	req, err := request(sent)
	if err != nil {
		t.Fatalf("unable to make request %d: %v", sent, err)
	}

	return req
}

// call sends one request of the suite to the engine.
func call(eng Engine, req Request) (Result, error) {
	if req.Embeddings {
		return eng.Embed(req)
	}

	return eng.Generate(req)
}

type skipper interface {
	Skip(...any)
	Fatalf(string, ...any)
}

// setup makes the engine one time and puts the model in memory.
func setup(t skipper) Engine {
	if engine != nil {
		return engine
	}

	if suite == "multimodal" && images == nil {
		var size image.Point
		if imageSize != "" {
			if _, err := fmt.Sscanf(imageSize, "%dx%d", &size.X, &size.Y); err != nil {
				t.Fatalf("the image size %q is not WIDTHxHEIGHT: %v", imageSize, err)
			}
		}

		made, err := NewImages(imageFile, size)
		if err != nil {
			t.Fatalf("unable to read the image: %v", err)
		}
		images = made
	}

	req, err := request(0)
	if err != nil {
		t.Fatalf("unable to make the first request: %v", err)
	}

	switch engineName {
	case "yzma":
		if library == "" {
			t.Skip("no YZMA_LIB and no -lib, skipping")
		}
		if modelFile == "" {
			t.Skip("no -model, skipping")
		}
		if len(req.Image) > 0 && projector == "" {
			t.Skip("no -mmproj, skipping")
		}

		engine = NewYzma(YzmaOptions{
			Library:   library,
			Model:     modelFile,
			Projector: projector,
			Device:    device,
			NCtx:      nCtx,
			NBatch:    nBatch,

			ImageMinTokens: int32(imgMin),
			ImageMaxTokens: int32(imgMax),
			Threads:        int32(threads),
			Embeddings:     req.Embeddings,
		})
	case "ollama", "dmr":
		if serverModel == "" {
			t.Skip("no -server-model, skipping")
		}

		url := ollamaURL
		if engineName == "dmr" {
			url = dmrURL
		}
		engine = NewServer(engineName, url, serverModel)
	default:
		t.Fatalf("unknown engine: %s", engineName)
	}

	if err := engine.Load(req); err != nil {
		engine = nil
		t.Fatalf("%s did not load: %v", engineName, err)
	}

	return engine
}

func BenchmarkCompare(b *testing.B) {
	eng := setup(b)

	var tokens, promptTokens int
	var ttft, total time.Duration
	runs := 0

	b.ResetTimer()
	for b.Loop() {
		// The request is made before the call, thus its image is not timed.
		result, err := call(eng, next(b))
		if err != nil {
			b.Fatalf("%s failed: %v", eng.Name(), err)
		}

		tokens += result.Tokens
		promptTokens = result.PromptTokens
		ttft += result.FirstToken
		total += result.Total
		runs++
	}

	// The tokens a second come from the time of the requests, not from the
	// time of the whole loop. The loop also holds the work of the harness, and
	// the emptying of the cache of yzma, which costs 4 ms and which a server
	// does inside its own request. Both engines are therefore measured from the
	// call to the answer and from nothing else.
	seconds := total.Seconds()
	if seconds == 0 {
		seconds = b.Elapsed().Seconds()
	}
	// The embeddings suite makes no token, thus it counts the tokens of the
	// prompt that the engine read each second.
	made := tokens
	if suite == "embeddings" {
		made = promptTokens * runs
	}
	b.ReportMetric(float64(made)/seconds, "tokens/s")
	b.ReportMetric(millis(ttft, runs), "ttft_ms")
	b.ReportMetric(millis(total, runs), "total_ms")
	// The count of the prompt tokens says if the engines do the same work.
	b.ReportMetric(float64(promptTokens), "prompt_tokens")
}

func millis(d time.Duration, runs int) float64 {
	if runs == 0 {
		return 0
	}

	return float64(d.Microseconds()) / 1000 / float64(runs)
}

// TestAnswer prints what the engine says. The engines must agree, or the
// numbers of the benchmark compare different work.
func TestAnswer(t *testing.T) {
	eng := setup(t)

	result, err := call(eng, next(t))
	if err != nil {
		t.Fatalf("%s failed: %v", eng.Name(), err)
	}

	if result.Dimensions > 0 {
		fmt.Printf("answer %s/%s: %d prompt tokens, a vector of %d\n",
			eng.Name(), suite, result.PromptTokens, result.Dimensions)

		return
	}

	fmt.Printf("answer %s/%s: %d prompt tokens, %d tokens: %q\n",
		eng.Name(), suite, result.PromptTokens, result.Tokens, result.Text)
}

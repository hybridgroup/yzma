package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	modelFile  *string
	configFile *string
	readout    *string
	state      *string
	question   *string
	options    *string
	qtype      *string
	category   *string
	libPath    *string
	verbose    *bool
	threads    *int
)

func showUsage() {
	fmt.Println(`
Usage:
decide -model [model file path] -config [config path] -readout [jev|jevk5] -lib [llama.cpp .so file path] -state [text] -question [text] -type [choice|score|noul] -options [JSON] -category [name] -v`)
}

func handleFlags() error {
	modelFile = flag.String("model", "", "model file to use")
	configFile = flag.String("config", "", "readout_config.json for jev, jevk5_config.json for jevk5")
	readout = flag.String("readout", "jev", "model family, jev for Jev-Style or jevk5 for JevK5")
	state = flag.String("state", "", "state as text or JSON")
	question = flag.String("question", "", "question text")
	options = flag.String("options", "", `JSON options, ["a","b"] or {"a":"desc"} for choice, ["level 0",...] for score, {"false":"desc","true":"desc"} for noul`)
	qtype = flag.String("type", "", "question type choice, score or noul (default choice with options, noul without)")
	category = flag.String("category", "", "calibration category such as general_sentiment or theme_routing")
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files")
	verbose = flag.Bool("v", false, "verbose logging")
	threads = flag.Int("t", 0, "number of CPU threads (0 = llama.cpp default)")

	flag.Parse()

	if len(*modelFile) == 0 {
		return errors.New("missing model flag")
	}

	if len(*configFile) == 0 {
		return errors.New("missing config flag")
	}

	if *readout != "jev" && *readout != "jevk5" {
		return errors.New("readout flag must be jev or jevk5")
	}

	if len(*question) == 0 {
		return errors.New("missing question flag")
	}

	if len(*options) > 0 && !json.Valid([]byte(*options)) {
		return errors.New("options flag is not valid JSON")
	}

	if len(*libPath) == 0 && os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	return nil
}

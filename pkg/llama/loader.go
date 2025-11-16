package llama

import (
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/loader"
)

// Load loads the shared llama.cpp libraries from the specified path.
func Load(path string) error {
	lib, err := loader.LoadLibrary(path, "ggml")
	if err != nil {
		return err
	}

	if err := loadGGML(lib); err != nil {
		return err
	}

	lib, err = loader.LoadLibrary(path, "ggml-base")
	if err != nil {
		return err
	}

	if err := loadGGMLBase(lib); err != nil {
		return err
	}

	lib, err = loader.LoadLibrary(path, "llama")
	if err != nil {
		return err
	}

	if err := loadBackendFuncs(lib); err != nil {
		return err
	}

	if err := loadModelFuncs(lib); err != nil {
		return err
	}

	if err := loadBatchFuncs(lib); err != nil {
		return err
	}

	if err := loadVocabFuncs(lib); err != nil {
		return err
	}

	if err := loadSamplingFuncs(lib); err != nil {
		return err
	}

	if err := loadChatFuncs(lib); err != nil {
		return err
	}

	if err := loadContextFuncs(lib); err != nil {
		return err
	}

	if err := loadMemoryFuncs(lib); err != nil {
		return err
	}

	if err := loadLogFuncs(lib); err != nil {
		return err
	}

	if err := loadStateFuncs(lib); err != nil {
		return err
	}

	if err := loadLoraFuncs(lib); err != nil {
		return err
	}

	return nil
}

// Init is a convenience function to handle initialization of llama.cpp.
func Init() {
	BackendInit()

	if os.Getenv("YZMA_LIB") != "" {
		GGMLBackendLoadAllFromPath(os.Getenv("YZMA_LIB"))

		return
	}

	GGMLBackendLoadAll()
}

func loadError(name string, err error) error {
	return fmt.Errorf("could not load %q: %w", name, err)
}

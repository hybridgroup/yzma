package llama

import (
	"fmt"
	"os"

	"github.com/jupiterrider/ffi"
)

func Load(lib ffi.Lib) error {
	if err := loadFuncs(lib); err != nil {
		return fmt.Errorf("loadFuncs: %w", err)
	}

	if err := loadModelFuncs(lib); err != nil {
		return fmt.Errorf("loadModelFuncs: %w", err)
	}

	if err := loadBatchFuncs(lib); err != nil {
		return fmt.Errorf("loadBatchFuncs: %w", err)
	}

	if err := loadVocabFuncs(lib); err != nil {
		return fmt.Errorf("loadVocabFuncs: %w", err)
	}

	if err := loadSamplingFuncs(lib); err != nil {
		return fmt.Errorf("loadSamplingFuncs: %w", err)
	}

	if err := loadChatFuncs(lib); err != nil {
		return fmt.Errorf("loadChatFuncs: %w", err)
	}

	if err := loadContextFuncs(lib); err != nil {
		return fmt.Errorf("loadContextFuncs: %w", err)
	}

	if err := loadLogFuncs(lib); err != nil {
		return fmt.Errorf("loadLogFuncs: %w", err)
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

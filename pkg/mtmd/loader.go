package mtmd

import (
	"sync"

	"github.com/hybridgroup/yzma/pkg/loader"
)

var muHelperEvalChunks sync.Mutex

func Load(path string) error {
	lib, err := loader.LoadLibrary(path, "mtmd")
	if err != nil {
		return err
	}

	if err := loadFuncs(lib); err != nil {
		return err
	}

	if err := loadBitmapFuncs(lib); err != nil {
		return err
	}

	return nil
}

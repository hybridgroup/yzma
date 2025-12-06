package download

import (
	"context"

	getter "github.com/hashicorp/go-getter"
)

// GetModel downloads a model from the specified URL to the destination path.
func GetModel(url, dest string) error {
	return getModel(url, dest)
}

func getModel(url, dest string) error {
	client := &getter.Client{
		Ctx:  context.Background(),
		Src:  url,
		Dst:  dest,
		Mode: getter.ClientModeAny,
	}

	if ProgressTracker != nil {
		client.ProgressListener = ProgressTracker
	}

	if err := client.Get(); err != nil {
		return err
	}

	return nil
}

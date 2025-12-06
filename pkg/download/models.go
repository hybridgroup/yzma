package download

import (
	"context"

	getter "github.com/hashicorp/go-getter"
)

// GetModel downloads a model from the specified URL to the destination path.
func GetModel(url, dest string) error {
	return getModel(url, dest, ProgressTracker)
}

// GetModelWithProgress downloads a model from the specified URL to the destination path
// using the provided progress tracker.
func GetModelWithProgress(url, dest string, progress getter.ProgressTracker) error {
	return getModel(url, dest, progress)
}

func getModel(url, dest string, progress getter.ProgressTracker) error {
	client := &getter.Client{
		Ctx:  context.Background(),
		Src:  url,
		Dst:  dest,
		Mode: getter.ClientModeAny,
	}

	if progress != nil {
		client.ProgressListener = progress
	}

	if err := client.Get(); err != nil {
		return err
	}

	return nil
}

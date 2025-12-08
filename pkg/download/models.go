package download

import (
	"context"

	getter "github.com/hashicorp/go-getter"
)

// GetModel downloads a model from the specified URL to the destination path.
func GetModel(url, dest string) error {
	return getModel(context.Background(), url, dest, ProgressTracker)
}

// GetModelWithProgress downloads a model from the specified URL to the destination path
// using the provided progress tracker.
func GetModelWithProgress(url, dest string, progress getter.ProgressTracker) error {
	return getModel(context.Background(), url, dest, progress)
}

// GetModelWithContext downloads a model from the specified URL to the destination path
// using the provided context.
func GetModelWithContext(ctx context.Context, url, dest string, progress getter.ProgressTracker) error {
	return getModel(ctx, url, dest, progress)
}

func getModel(ctx context.Context, url, dest string, progress getter.ProgressTracker) error {
	client := &getter.Client{
		Ctx:  ctx,
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

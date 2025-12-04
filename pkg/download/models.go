package download

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	getter "github.com/hashicorp/go-getter"
)

// GetModel downloads a model from the specified URL to the destination path.
func GetModel(url, dest string, showProgress bool) error {
	return getModel(url, dest, showProgress)
}

func getModel(url, dest string, showProgress bool) error {
	client := &getter.Client{
		Ctx:  context.Background(),
		Src:  url,
		Dst:  dest,
		Mode: getter.ClientModeAny,
	}

	if showProgress {
		progFunc := func(src string, currentSize int64, totalSize int64, mibPerSec float64, complete bool) {
			fmt.Printf("\r\x1b[KDownloading %s... %d MiB of %d MiB (%.2f MiB/s)", src, currentSize/(1024*1024), totalSize/(1024*1024), mibPerSec)
			if complete {
				fmt.Println()
			}
		}

		pr := progressReader{
			dst:      dest,
			progress: progFunc,
		}

		client.ProgressListener = getter.ProgressTracker(&pr)
	}

	if err := client.Get(); err != nil {
		return err
	}

	return nil
}

type progressFunc func(src string, currentSize int64, totalSize int64, mibPerSec float64, complete bool)

type progressReader struct {
	src          string
	dst          string
	currentSize  int64
	totalSize    int64
	lastReported int64
	startTime    time.Time
	reader       io.ReadCloser
	progress     progressFunc
}

func (pr *progressReader) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	if currentSize == totalSize {
		return nil
	}

	if currentSize != totalSize {
		os.Remove(pr.dst)
	}

	pr.src = src
	pr.currentSize = currentSize
	pr.totalSize = totalSize
	pr.startTime = time.Now()
	pr.reader = stream

	return pr
}

const (
	mib    = 1024 * 1024
	mib100 = mib * 100
)

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.currentSize += int64(n)

	if pr.progress != nil && pr.currentSize-pr.lastReported >= mib100 {
		pr.lastReported = pr.currentSize
		pr.progress(pr.src, pr.currentSize, pr.totalSize, pr.mibPerSec(), false)
	}

	return n, err
}

func (pr *progressReader) Close() error {
	if pr.progress != nil {
		pr.progress(pr.src, pr.currentSize, pr.totalSize, pr.mibPerSec(), true)
	}

	return pr.reader.Close()
}

func (pr *progressReader) mibPerSec() float64 {
	elapsed := time.Since(pr.startTime).Seconds()
	if elapsed == 0 {
		return 0
	}

	return float64(pr.currentSize) / mib / elapsed
}

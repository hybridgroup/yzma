package download

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type recordingTracker struct {
	mu   sync.Mutex
	srcs []string
}

func (r *recordingTracker) TrackProgress(src string, _, _ int64, stream io.ReadCloser) io.ReadCloser {
	r.mu.Lock()
	r.srcs = append(r.srcs, src)
	r.mu.Unlock()
	return stream
}

// TestGetModelConcurrent holds every response until all requests arrive, so the
// downloads overlap. Each one must get its own bytes and its own progress tracker.
func TestGetModelConcurrent(t *testing.T) {
	const n = 8
	payload := func(i int) []byte {
		return bytes.Repeat([]byte(fmt.Sprintf("model-%d|", i)), 64*1024)
	}

	var arrived sync.WaitGroup
	arrived.Add(n)
	release := make(chan struct{})
	go func() {
		arrived.Wait()
		close(release)
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var i int
		if _, err := fmt.Sscanf(r.URL.Path, "/model-%d.gguf", &i); err != nil {
			http.NotFound(w, r)
			return
		}
		data := payload(i)
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		if r.Method == http.MethodHead {
			return
		}
		arrived.Done()
		select {
		case <-release:
		case <-time.After(10 * time.Second):
		}
		for len(data) > 0 {
			k := min(32*1024, len(data))
			w.Write(data[:k])
			w.(http.Flusher).Flush()
			data = data[k:]
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	trackers := make([]*recordingTracker, n)
	errs := make([]error, n)
	var wg sync.WaitGroup
	for i := range n {
		trackers[i] = &recordingTracker{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("%s/model-%d.gguf", server.URL, i)
			dest := filepath.Join(dir, fmt.Sprint(i))
			errs[i] = GetModelWithContext(context.Background(), url, dest, trackers[i])
		}()
	}
	wg.Wait()

	for i := range n {
		if errs[i] != nil {
			t.Fatalf("download %d failed: %v", i, errs[i])
		}
		got, err := os.ReadFile(filepath.Join(dir, fmt.Sprint(i), fmt.Sprintf("model-%d.gguf", i)))
		if err != nil {
			t.Fatalf("read %d: %v", i, err)
		}
		if !bytes.Equal(got, payload(i)) {
			t.Fatalf("download %d has wrong content", i)
		}
		want := fmt.Sprintf("model-%d.gguf", i)
		if len(trackers[i].srcs) != 1 || !strings.HasSuffix(trackers[i].srcs[0], want) {
			t.Fatalf("tracker %d saw %v, want one call for %s", i, trackers[i].srcs, want)
		}
	}
}

func TestNewGettersNotShared(t *testing.T) {
	a, b := newGetters(), newGetters()
	for scheme := range a {
		if a[scheme] == b[scheme] {
			t.Fatalf("getter for %q is shared between clients", scheme)
		}
	}
}

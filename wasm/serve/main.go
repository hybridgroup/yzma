// Serve serves the WebAssembly build of yzma on a local HTTP server.
//
// The server sets the Cross-Origin-Opener-Policy and
// Cross-Origin-Embedder-Policy headers. Without them the browser does not
// expose SharedArrayBuffer, and the multithreaded build cannot run.
//
//	go run ./wasm/serve -dir build/wasm -port 8080
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	dir := flag.String("dir", "build/wasm", "directory of files to serve")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	path, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		log.Fatalf("cannot serve %s: %v (run \"make wasm-example\" first)", path, err)
	}

	files := http.FileServer(http.Dir(path))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// These two headers isolate the page, which enables SharedArrayBuffer
		// and so multiple threads.
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		// The browser revalidates each file, so it never runs a stale build.
		w.Header().Set("Cache-Control", "no-cache")

		// A browser only runs a .wasm file served with the correct type.
		if filepath.Ext(r.URL.Path) == ".wasm" {
			w.Header().Set("Content-Type", "application/wasm")
		}

		files.ServeHTTP(w, r)
	})

	address := fmt.Sprintf("localhost:%d", *port)
	log.Printf("serving %s at http://%s", path, address)
	log.Fatal(http.ListenAndServe(address, nil))
}

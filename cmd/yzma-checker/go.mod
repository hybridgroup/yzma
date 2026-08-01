// This is a NESTED module on purpose.
//
// Keeping its own go.mod means the golang.org/x/tools dependency never reaches
// the root yzma module, and `go build ./...` / `go vet ./...` / `go test ./...`
// at the yzma repo root skip this directory entirely.
module github.com/hybridgroup/yzma/cmd/yzma-checker

go 1.26

require golang.org/x/tools v0.48.0

require (
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
)

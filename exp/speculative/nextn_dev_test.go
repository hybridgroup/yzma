package speculative

import (
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestNextNDevZeroHandles(t *testing.T) {
	// The guards must hold before any library is loaded.
	if SetEmbeddingsNextNDev(0, true) {
		t.Fatal("SetEmbeddingsNextNDev accepted a zero context")
	}
	if GetEmbeddingsNextNRows(0, 0, 4, make([]float32, 4)) {
		t.Fatal("GetEmbeddingsNextNRows accepted a zero context")
	}
	if SetEmbeddingsNextNRows(0, 0, 4, make([]float32, 4)) {
		t.Fatal("SetEmbeddingsNextNRows accepted a zero context")
	}
	if SetCtxStateSrc(0, 0) {
		t.Fatal("SetCtxStateSrc accepted a zero context")
	}
	if BatchExtSetEmbdStateRow(0, 0, 0) {
		t.Fatal("BatchExtSetEmbdStateRow accepted a zero batch")
	}
}

func TestRowCount(t *testing.T) {
	tests := []struct {
		n, nEmbd int
		want     int64
		ok       bool
	}{
		{8, 4, 2, true},
		{4, 4, 1, true},
		{0, 4, 0, false},
		{6, 4, 0, false},
		{4, 0, 0, false},
		{4, -1, 0, false},
	}
	for _, tt := range tests {
		got, ok := rowCount(tt.n, tt.nEmbd)
		if got != tt.want || ok != tt.ok {
			t.Errorf("rowCount(%d, %d) = %d, %v, want %d, %v", tt.n, tt.nEmbd, got, ok, tt.want, tt.ok)
		}
	}
}

func TestDeviceRowsAvailable(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// An older llama.cpp build without the device row functions is a valid state.
	t.Logf("DeviceRowsAvailable returned: %v", DeviceRowsAvailable())
}

func TestNextNDevRows(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if !Available() || !DeviceRowsAvailable() {
		t.Skip("device row functions not in this llama.cpp build")
	}

	model, err := llama.ModelLoadFromFile(testModelFileName(t), llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := llama.ContextDefaultParams()
	params.NBatch = 32
	params.NSeqMax = 2
	ctx, err := llama.InitFromModel(model, params)
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer llama.Free(ctx)

	// The device rows are batch rows, so a masked context must be refused.
	SetEmbeddingsNextN(ctx, true, true)
	if SetEmbeddingsNextNDev(ctx, true) {
		t.Fatal("SetEmbeddingsNextNDev accepted a masked context")
	}

	SetEmbeddingsNextN(ctx, true, false)
	if !SetEmbeddingsNextNDev(ctx, true) {
		t.Skip("no device buffer for this model and backend")
	}
	defer SetEmbeddingsNextNDev(ctx, false)

	nEmbd := int(llama.ModelNEmbdOut(model))
	nRows := int64(llama.NBatch(ctx) + llama.NSeqMax(ctx))

	// Scratch rows start at n_batch. Use two of them at a nonzero offset.
	row0 := int64(llama.NBatch(ctx))
	data := make([]float32, 2*nEmbd)
	for i := range data {
		data[i] = float32(i) + 0.5
	}
	if !SetEmbeddingsNextNRows(ctx, row0, nEmbd, data) {
		t.Fatal("SetEmbeddingsNextNRows failed")
	}

	out := make([]float32, 2*nEmbd)
	if !GetEmbeddingsNextNRows(ctx, row0, nEmbd, out) {
		t.Fatal("GetEmbeddingsNextNRows failed")
	}
	for i := range out {
		if out[i] != data[i] {
			t.Fatalf("out[%d] = %v, want %v", i, out[i], data[i])
		}
	}

	// A single row read must match the same row of the pair.
	one := make([]float32, nEmbd)
	if !GetEmbeddingsNextNRows(ctx, row0+1, nEmbd, one) {
		t.Fatal("GetEmbeddingsNextNRows of one row failed")
	}
	for i := range one {
		if one[i] != data[nEmbd+i] {
			t.Fatalf("one[%d] = %v, want %v", i, one[i], data[nEmbd+i])
		}
	}

	if GetEmbeddingsNextNRows(ctx, nRows-1, nEmbd, out) {
		t.Fatal("GetEmbeddingsNextNRows accepted rows past the end")
	}
	if GetEmbeddingsNextNRows(ctx, -1, nEmbd, one) {
		t.Fatal("GetEmbeddingsNextNRows accepted a negative row")
	}
	if SetEmbeddingsNextNRows(ctx, 0, nEmbd, data[:nEmbd+1]) {
		t.Fatal("SetEmbeddingsNextNRows accepted a partial row")
	}

	// Only an MTP context reads state rows.
	if SetCtxStateSrc(ctx, ctx) {
		t.Fatal("SetCtxStateSrc accepted a default context")
	}
	if SetCtxStateSrc(ctx, 0) {
		t.Fatal("SetCtxStateSrc detach returned true")
	}
}

func TestBatchExtSetEmbdStateRow(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if !DeviceRowsAvailable() {
		t.Skip("device row functions not in this llama.cpp build")
	}

	model, err := llama.ModelLoadFromFile(testModelFileName(t), llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	ctx, err := llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer llama.Free(ctx)

	batch, err := llama.BatchExtInit(ctx)
	if err != nil {
		t.Fatalf("BatchExtInit failed: %v", err)
	}
	defer llama.BatchExtFree(batch)

	if BatchExtSetEmbdStateRow(batch, 0, 0) {
		t.Fatal("BatchExtSetEmbdStateRow accepted an empty batch")
	}

	for i := range 2 {
		if _, err := llama.BatchExtAddToken(batch, 0, 1); err != nil {
			t.Fatalf("BatchExtAddToken %d failed: %v", i, err)
		}
	}

	if !BatchExtSetEmbdStateRow(batch, 0, 3) {
		t.Fatal("BatchExtSetEmbdStateRow failed on entry 0")
	}
	if !BatchExtSetEmbdStateRow(batch, 1, 7) {
		t.Fatal("BatchExtSetEmbdStateRow failed on entry 1")
	}
	if BatchExtSetEmbdStateRow(batch, 0, 4) {
		t.Fatal("BatchExtSetEmbdStateRow set a second row on entry 0")
	}
	if BatchExtSetEmbdStateRow(batch, 2, 0) {
		t.Fatal("BatchExtSetEmbdStateRow accepted an index past the end")
	}
	if BatchExtSetEmbdStateRow(batch, 1, -1) {
		t.Fatal("BatchExtSetEmbdStateRow accepted a negative row")
	}
}

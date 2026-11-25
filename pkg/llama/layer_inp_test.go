package llama

import (
	"math"
	"testing"
)

func TestLayerInpInvalidHandles(t *testing.T) {
	if ids := ModelTargetLayerIDs(0); ids != nil {
		t.Fatal("ModelTargetLayerIDs returned IDs for a zero model")
	}
	if err := SetEmbeddingsLayerInp(0, 0, true); err == nil {
		t.Fatal("SetEmbeddingsLayerInp accepted a zero context")
	}
	if _, err := GetEmbeddingsLayerInp(0, 0, 1, 1); err == nil {
		t.Fatal("GetEmbeddingsLayerInp accepted a zero context")
	}
	if SupportsEmbdRow(0) {
		t.Fatal("SupportsEmbdRow returned true for a zero context")
	}
	if _, err := BatchExtAddEmbdRow(0, 0, 0); err == nil {
		t.Fatal("BatchExtAddEmbdRow accepted a zero batch")
	}
}

func TestModelTargetLayerIDs(t *testing.T) {
	model, _ := testLayerInpContext(t)

	if ids := ModelTargetLayerIDs(model); ids != nil {
		t.Fatalf("ModelTargetLayerIDs returned %v for a model that is not a drafter", ids)
	}
}

func TestEmbeddingsLayerInp(t *testing.T) {
	model, ctx := testLayerInpContext(t)

	nLayer := uint32(ModelNLayer(model))
	nEmbd := int(ModelNEmbd(model))
	layers := []uint32{0, nLayer / 2, nLayer - 1}

	if err := SetEmbeddingsLayerInp(ctx, nLayer+1, true); err == nil {
		t.Fatal("SetEmbeddingsLayerInp accepted a layer past the output")
	}
	if _, err := GetEmbeddingsLayerInp(ctx, nLayer+1, 1, nEmbd); err == nil {
		t.Fatal("GetEmbeddingsLayerInp accepted a layer past the output")
	}
	for _, lid := range layers {
		if err := SetEmbeddingsLayerInp(ctx, lid, true); err != nil {
			t.Fatalf("SetEmbeddingsLayerInp(%d) failed: %v", lid, err)
		}
	}

	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, "This is a test of layer inputs", true, true)
	if len(tokens) < 4 {
		t.Fatalf("Tokenize returned %d tokens", len(tokens))
	}

	if _, err := Decode(ctx, BatchGetOne(tokens)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	full := make(map[uint32][]float32)
	for _, lid := range layers {
		feat, err := GetEmbeddingsLayerInp(ctx, lid, len(tokens), nEmbd)
		if err != nil {
			t.Fatalf("GetEmbeddingsLayerInp(%d) failed: %v", lid, err)
		}
		if len(feat) != len(tokens)*nEmbd {
			t.Fatalf("GetEmbeddingsLayerInp(%d) returned %d values, expected %d", lid, len(feat), len(tokens)*nEmbd)
		}
		for i, v := range feat {
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
				t.Fatalf("layer %d value %d is %v", lid, i, v)
			}
		}
		full[lid] = append([]float32(nil), feat...)
	}
	if maxDiff(full[layers[0]], full[layers[1]]) == 0 {
		t.Fatal("two different layers returned the same inputs")
	}

	// Decode the same tokens again in two batches and compare row by row.
	mem, err := GetMemory(ctx)
	if err != nil {
		t.Fatalf("GetMemory failed: %v", err)
	}
	MemoryClear(mem, true)

	split := 2
	parts := [][]Token{tokens[:split], tokens[split:]}
	offset := 0
	for _, part := range parts {
		batch := BatchInit(int32(len(part)), 0, 1)
		for i, tok := range part {
			batch.Add(tok, Pos(offset+i), []SeqId{0}, i == len(part)-1)
		}
		if _, err := Decode(ctx, batch); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		BatchFree(batch)

		for _, lid := range layers {
			feat, err := GetEmbeddingsLayerInp(ctx, lid, len(part), nEmbd)
			if err != nil {
				t.Fatalf("GetEmbeddingsLayerInp(%d) failed: %v", lid, err)
			}
			want := full[lid][offset*nEmbd : (offset+len(part))*nEmbd]
			if d := maxDiff(feat, want) / maxAbs(want); d > 2e-2 {
				t.Fatalf("layer %d rows %d to %d differ from the single batch by %g of the scale", lid, offset, offset+len(part), d)
			}
		}
		offset += len(part)
	}
}

func TestEmbeddingsLayerInpDev(t *testing.T) {
	_, ctx := testLayerInpContext(t)

	if err := SetEmbeddingsLayerInp(ctx, 1, true); err != nil {
		t.Fatalf("SetEmbeddingsLayerInp failed: %v", err)
	}
	SetEmbeddingsLayerInpDev(ctx, true)
	if SupportsEmbdRow(ctx) {
		t.Fatal("SupportsEmbdRow returned true for a context that is not a DFlash drafter")
	}
	SetEmbeddingsLayerInpDev(ctx, false)

	batch, err := BatchExtInit(ctx)
	if err != nil {
		t.Fatalf("BatchExtInit failed: %v", err)
	}
	defer BatchExtFree(batch)

	for want := int32(0); want < 3; want++ {
		idx, err := BatchExtAddEmbdRow(batch, 0, want)
		if err != nil {
			t.Fatalf("BatchExtAddEmbdRow failed: %v", err)
		}
		if idx != want {
			t.Fatalf("BatchExtAddEmbdRow returned index %d, expected %d", idx, want)
		}
	}
}

func maxDiff(a, b []float32) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}
	var d float64
	for i := range a {
		d = max(d, math.Abs(float64(a[i]-b[i])))
	}
	return d
}

func testLayerInpContext(t *testing.T) (Model, Context) {
	t.Helper()
	testSetup(t)

	return testBatchExtContext(t, ContextDefaultParams())
}

func maxAbs(a []float32) float64 {
	var m float64
	for _, v := range a {
		m = max(m, math.Abs(float64(v)))
	}
	return m
}

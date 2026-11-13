//go:build js && wasm

package llamawasm

import (
	"errors"
	"syscall/js"
	"testing"
)

// The fake ABI version 14 module reads the arrays that yzma_process_embd gets.
const fakeABI14Source = `
globalThis.__yzmaABI14 = (function () {
	const heap = new Uint8Array(1 << 16);
	let next = 8;
	let last = null;

	function ints(p, n) {
		return p === 0 ? null : Array.from(new Int32Array(heap.buffer, p, n));
	}
	function floats(p, n) {
		return p === 0 ? null : Array.from(new Float32Array(heap.buffer, p, n));
	}

	return {
		module: {
			HEAPU8: heap,
			_malloc: (n) => { const p = next; next += n + (8 - (n % 8)); return p; },
			_free: () => {},
			_yzma_abi_version: () => 14,
			_yzma_last_error: () => 0,
			_yzma_context_n_batch: () => 4,
			_yzma_context_n_seq_max: () => 2,
			_yzma_process: () => { last = "process"; return 0; },
			_yzma_process_embd: (ctx, type, tokens, pos, nSeq, seqs, nSeqMax, outputs, order, embd, desc, n) => {
				const d = ints(desc, n * 4);
				let nFloats = 0;
				for (let i = 0; i < n; i++) nFloats = Math.max(nFloats, d[i*4] + d[i*4+1], d[i*4+2] + d[i*4+3]);
				last = {
					ctx, type, nSeqMax,
					tokens: ints(tokens, n),
					pos: ints(pos, n * 4),
					nSeq: ints(nSeq, n),
					seqs: ints(seqs, n * nSeqMax),
					outputs: Array.from(new Int8Array(heap.buffer, outputs, n)),
					order: ints(order, n),
					embd: floats(embd, nFloats),
					desc: d,
				};
				return 0;
			},
		},
		last: () => JSON.stringify(last),
	};
})();
`

func fakeABI14(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeABI14Source)
	helper := js.Global().Get("__yzmaABI14")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() {
		mod = previous
		for _, s := range []*scratch{&tokenScratch, &posScratch, &nSeqScratch, &seqScratch, &logitScratch,
			&orderScratch, &embdInScratch, &descScratch} {
			s.ptr, s.size = 0, 0
		}
	})

	return helper
}

func TestProcessEmbd(t *testing.T) {
	helper := fakeABI14(t)

	bt, err := BatchExtInit(Context(7))
	if err != nil {
		t.Fatalf("BatchExtInit gave %v", err)
	}

	idx, err := BatchExtAddToken(bt, 0, 11)
	if err != nil {
		t.Fatalf("BatchExtAddToken gave %v", err)
	}
	BatchExtSetPos(bt, idx, 3)

	img := []float32{0.5, 1.5, 2.5, 3.5}
	idx, err = BatchExtAddEmbd(bt, 1, img, 2)
	if err != nil || idx != 1 {
		t.Fatalf("BatchExtAddEmbd gave %d, %v", idx, err)
	}
	img[0] = 9
	if err := BatchExtSetPos(bt, idx, 4, 1, 2); err != nil {
		t.Fatalf("BatchExtSetPos gave %v", err)
	}
	if err := BatchExtSetEmbdState(bt, idx, []float32{7}, 1); err != nil {
		t.Fatalf("BatchExtSetEmbdState gave %v", err)
	}
	if err := BatchExtSetOutputLogits(bt, idx, true); err != nil {
		t.Fatalf("BatchExtSetOutputLogits gave %v", err)
	}

	idx, err = BatchExtAdd(bt, 0)
	if err != nil || idx != 2 {
		t.Fatalf("BatchExtAdd gave %d, %v", idx, err)
	}
	BatchExtSetPos(bt, idx, 5)

	ret, err := Process(Context(7), ProcessTypeDecode, bt)
	if err != nil || ret != 0 {
		t.Fatalf("Process gave %d, %v", ret, err)
	}

	want := `{"ctx":7,"type":1,"nSeqMax":1,"tokens":[11,-1,-1],"pos":[3,0,0,0,4,1,2,0,5,0,0,0],` +
		`"nSeq":[1,1,1],"seqs":[0,1,0],"outputs":[0,1,0],"order":null,` +
		`"embd":[0.5,1.5,2.5,3.5,7],"desc":[0,0,0,0,0,4,4,1,5,0,5,0]}`
	if got := helper.Call("last").String(); got != want {
		t.Errorf("the shim got\n%s\nwant\n%s", got, want)
	}
}

func TestProcessEmbdTokenOnly(t *testing.T) {
	helper := fakeABI14(t)

	bt, err := BatchExtInit(Context(7))
	if err != nil {
		t.Fatalf("BatchExtInit gave %v", err)
	}
	idx, _ := BatchExtAddToken(bt, 0, 3)
	BatchExtSetPos(bt, idx, 0)

	if _, err := Process(Context(7), ProcessTypeDecode, bt); err != nil {
		t.Fatalf("Process gave %v", err)
	}
	if got := helper.Call("last").String(); got != `"process"` {
		t.Errorf("a batch of tokens went to %s, want yzma_process", got)
	}

	if err := BatchExtSetEmbdToken(bt, idx, []float32{1, 2}, 2); err != nil {
		t.Fatalf("BatchExtSetEmbdToken gave %v", err)
	}
	if _, err := Process(Context(7), ProcessTypeDecode, bt); err != nil {
		t.Fatalf("Process gave %v", err)
	}
	want := `{"ctx":7,"type":1,"nSeqMax":1,"tokens":[3],"pos":[0,0,0,0],"nSeq":[1],"seqs":[0],` +
		`"outputs":[0],"order":null,"embd":[1,2],"desc":[0,2,2,0]}`
	if got := helper.Call("last").String(); got != want {
		t.Errorf("the shim got\n%s\nwant\n%s", got, want)
	}

	BatchExtClear(bt)
	idx, _ = BatchExtAddToken(bt, 0, 4)
	BatchExtSetPos(bt, idx, 1)
	if _, err := Process(Context(7), ProcessTypeDecode, bt); err != nil {
		t.Fatalf("Process gave %v", err)
	}
	if got := helper.Call("last").String(); got != `"process"` {
		t.Errorf("a cleared batch of tokens went to %s, want yzma_process", got)
	}
}

func TestBatchExtEmbdRejects(t *testing.T) {
	fakeABI14(t)

	bt, err := BatchExtInit(Context(7))
	if err != nil {
		t.Fatalf("BatchExtInit gave %v", err)
	}
	if _, err := BatchExtAddEmbd(bt, 0, nil, 2); !errors.Is(err, ErrBatchExtInvalidToken) {
		t.Errorf("no data gave %v", err)
	}
	if _, err := BatchExtAddEmbd(bt, 0, []float32{1, 2, 3}, 2); !errors.Is(err, ErrBatchExtInvalidToken) {
		t.Errorf("a partial row gave %v", err)
	}
	if _, err := BatchExtAddEmbd(bt, 2, []float32{1}, 1); !errors.Is(err, ErrBatchExtInvalidSeqID) {
		t.Errorf("a sequence over n_seq_max gave %v", err)
	}
	if _, err := BatchExtAdd(bt, -1); !errors.Is(err, ErrBatchExtInvalidSeqID) {
		t.Errorf("a negative sequence gave %v", err)
	}

	idx, err := BatchExtAddEmbd(bt, 0, []float32{1}, 1)
	if err != nil {
		t.Fatalf("BatchExtAddEmbd gave %v", err)
	}
	if err := BatchExtSetEmbdToken(bt, idx, []float32{2}, 1); !errors.Is(err, ErrBatchExtRejected) {
		t.Errorf("a second token embedding gave %v", err)
	}
	if err := BatchExtSetEmbdState(bt, idx, []float32{2}, 0); !errors.Is(err, ErrBatchExtRejected) {
		t.Errorf("a state with nEmbd 0 gave %v", err)
	}
	if err := BatchExtSetEmbdState(bt, idx+1, []float32{2}, 1); !errors.Is(err, ErrBatchExtRejected) {
		t.Errorf("a state out of range gave %v", err)
	}
	for range 3 {
		if _, err := BatchExtAdd(bt, 0); err != nil {
			t.Fatalf("BatchExtAdd gave %v", err)
		}
	}
	if _, err := BatchExtAddEmbd(bt, 0, []float32{1}, 1); !errors.Is(err, ErrBatchFull) {
		t.Errorf("a full batch gave %v", err)
	}
}

func TestBatchExtEmbdOldModule(t *testing.T) {
	fakeABI11(t)

	bt, err := BatchExtInit(Context(7))
	if err != nil {
		t.Fatalf("BatchExtInit gave %v", err)
	}
	idx, err := BatchExtAddEmbd(bt, 0, []float32{1, 2}, 2)
	if err != nil {
		t.Fatalf("BatchExtAddEmbd gave %v", err)
	}
	BatchExtSetPos(bt, idx, 0)
	if _, err := Process(Context(7), ProcessTypeDecode, bt); !errors.Is(err, ErrNoBatchExtEmbd) {
		t.Errorf("Process on an ABI 11 module gave %v, want ErrNoBatchExtEmbd", err)
	}
}

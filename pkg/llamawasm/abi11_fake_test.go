//go:build js && wasm

package llamawasm

import (
	"errors"
	"syscall/js"
	"testing"
)

// The fake ABI version 11 module reads the arrays that yzma_process gets.
const fakeABI11Source = `
globalThis.__yzmaABI11 = (function () {
	const heap = new Uint8Array(1 << 16);
	let next = 8;
	let last = null;

	function ints(p, n) {
		return p === 0 ? null : Array.from(new Int32Array(heap.buffer, p, n));
	}

	return {
		module: {
			HEAPU8: heap,
			_malloc: (n) => { const p = next; next += n + (8 - (n % 8)); return p; },
			_free: () => {},
			_yzma_abi_version: () => 11,
			_yzma_last_error: () => 0,
			_yzma_context_n_batch: () => 3,
			_yzma_context_n_seq_max: () => 2,
			_yzma_process: (ctx, type, tokens, pos, nSeq, seqs, nSeqMax, outputs, order, n) => {
				last = {
					ctx, type, nSeqMax,
					tokens: ints(tokens, n),
					pos: ints(pos, n),
					nSeq: ints(nSeq, n),
					seqs: ints(seqs, n * nSeqMax),
					outputs: Array.from(new Int8Array(heap.buffer, outputs, n)),
					order: ints(order, n),
				};
				return 0;
			},
		},
		last: () => JSON.stringify(last),
	};
})();
`

func fakeABI11(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeABI11Source)
	helper := js.Global().Get("__yzmaABI11")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() {
		mod = previous
		for _, s := range []*scratch{&tokenScratch, &posScratch, &nSeqScratch, &seqScratch, &logitScratch, &orderScratch} {
			s.ptr, s.size = 0, 0
		}
	})

	return helper
}

func TestProcessDecisionOrder(t *testing.T) {
	helper := fakeABI11(t)

	bt, err := BatchExtInit(Context(7))
	if err != nil {
		t.Fatalf("BatchExtInit gave %v", err)
	}
	defer BatchExtFree(bt)

	for i, tok := range []Token{11, 12} {
		idx, err := BatchExtAddToken(bt, 1, tok)
		if err != nil {
			t.Fatalf("BatchExtAddToken gave %v", err)
		}
		if err := BatchExtSetPos(bt, idx, Pos(5+i)); err != nil {
			t.Fatalf("BatchExtSetPos gave %v", err)
		}
		if err := BatchExtSetOutputEmbd(bt, idx, i == 1); err != nil {
			t.Fatalf("BatchExtSetOutputEmbd gave %v", err)
		}
	}
	if err := BatchExtAddSeq(bt, 0, 0); err != nil {
		t.Fatalf("BatchExtAddSeq gave %v", err)
	}
	if err := BatchExtSetDecisionOrder(bt, 1, DecisionOrderOption); err != nil {
		t.Fatalf("BatchExtSetDecisionOrder gave %v", err)
	}

	ret, err := Process(Context(7), ProcessTypeDecode, bt)
	if err != nil || ret != 0 {
		t.Fatalf("Process gave %d, %v", ret, err)
	}

	want := `{"ctx":7,"type":1,"nSeqMax":2,"tokens":[11,12],"pos":[5,6],"nSeq":[2,1],"seqs":[1,0,1,0],"outputs":[0,1],"order":[0,4]}`
	if got := helper.Call("last").String(); got != want {
		t.Errorf("the shim got\n%s\nwant\n%s", got, want)
	}
}

func TestProcessNoOrder(t *testing.T) {
	helper := fakeABI11(t)

	bt, err := BatchExtInit(Context(7))
	if err != nil {
		t.Fatalf("BatchExtInit gave %v", err)
	}
	idx, _ := BatchExtAddToken(bt, 0, 3)
	BatchExtSetPos(bt, idx, 0)

	if _, err := Process(Context(7), ProcessTypeEncode, bt); err != nil {
		t.Fatalf("Process gave %v", err)
	}
	want := `{"ctx":7,"type":0,"nSeqMax":1,"tokens":[3],"pos":[0],"nSeq":[1],"seqs":[0],"outputs":[0],"order":null}`
	if got := helper.Call("last").String(); got != want {
		t.Errorf("the shim got\n%s\nwant\n%s", got, want)
	}
}

func TestBatchExtRejects(t *testing.T) {
	fakeABI11(t)

	bt, err := BatchExtInit(Context(7))
	if err != nil {
		t.Fatalf("BatchExtInit gave %v", err)
	}
	if _, err := BatchExtAddToken(bt, 2, 1); !errors.Is(err, ErrBatchExtInvalidSeqID) {
		t.Errorf("a sequence over n_seq_max gave %v", err)
	}
	for range 3 {
		if _, err := BatchExtAddToken(bt, 0, 1); err != nil {
			t.Fatalf("BatchExtAddToken gave %v", err)
		}
	}
	if _, err := BatchExtAddToken(bt, 0, 1); !errors.Is(err, ErrBatchFull) {
		t.Errorf("a full batch gave %v", err)
	}
	if err := BatchExtSetPos(bt, 3, 0); !errors.Is(err, ErrBatchExtRejected) {
		t.Errorf("an index out of range gave %v", err)
	}
	if err := BatchExtSetPos(bt, 0, 0, 1); !errors.Is(err, ErrBatchExtRejected) {
		t.Errorf("two positions gave %v", err)
	}
	if _, err := Process(Context(7), ProcessTypeDecode, bt); !errors.Is(err, ErrBatchExtRejected) {
		t.Errorf("entries with no position gave %v", err)
	}
	if err := BatchExtClear(bt); err != nil {
		t.Fatalf("BatchExtClear gave %v", err)
	}
	if ret, err := Process(Context(7), ProcessTypeDecode, bt); ret != 0 || err != nil {
		t.Errorf("an empty batch gave %d, %v", ret, err)
	}
}

func TestBatchExtOldModule(t *testing.T) {
	fakeABI10(t)

	if _, err := BatchExtInit(Context(7)); !errors.Is(err, ErrNoBatchExt) {
		t.Errorf("BatchExtInit gave %v, want ErrNoBatchExt", err)
	}
}

//go:build js && wasm

package llamawasm

import (
	"errors"
	"syscall/js"
	"testing"
)

// The fake ABI version 15 module reads the arrays of yzma_process_ext and
// runs a sampler that keeps the 2 highest logits.
const fakeABI15Source = `
globalThis.__yzmaABI15 = (function () {
	const heap = new Uint8Array(1 << 16);
	let next = 8;
	let last = null;
	const arrays = { "tokenizer.tts_dicts": ["first", "a much longer second element"] };

	function ints(p, n) {
		return p === 0 ? null : Array.from(new Int32Array(heap.buffer, p, n));
	}
	function floats(p, n) {
		return p === 0 ? null : Array.from(new Float32Array(heap.buffer, p, n));
	}
	function str(p) {
		let end = p;
		while (heap[end] !== 0) end++;
		return new TextDecoder().decode(heap.subarray(p, end));
	}
	function copy(s, buf, cap) {
		const b = new TextEncoder().encode(s);
		if (cap < b.length + 1) return -5;
		heap.set(b, buf);
		heap[buf + b.length] = 0;
		return b.length;
	}

	return {
		module: {
			HEAPU8: heap,
			_malloc: (n) => { const p = next; next += n + (8 - (n % 8)); return p; },
			_free: () => {},
			_yzma_abi_version: () => 15,
			_yzma_last_error: () => 0,
			_yzma_context_n_batch: () => 8,
			_yzma_context_n_seq_max: () => 2,
			_yzma_process: () => { last = "process"; return 0; },
			_yzma_process_embd: () => { last = "process_embd"; return 0; },
			_yzma_process_ext: (ctx, type, tokens, pos, nPos, nSeq, seqs, nSeqMax, outputs, order, group, embd, desc, n) => {
				let nFloats = 0;
				const d = ints(desc, n * 4);
				for (let i = 0; d && i < n; i++) nFloats = Math.max(nFloats, d[i*4] + d[i*4+1], d[i*4+2] + d[i*4+3]);
				last = {
					ctx, type, nPos,
					tokens: ints(tokens, n),
					pos: ints(pos, n * nPos),
					outputs: Array.from(new Int8Array(heap.buffer, outputs, n)),
					order: ints(order, n),
					group: ints(group, n),
					embd: floats(embd, nFloats),
					desc: d,
				};
				return 0;
			},
			_yzma_sampler_apply: (smpl, ids, logits, p, state, n) => {
				const id = new Int32Array(heap.buffer, ids, n);
				const lg = new Float32Array(heap.buffer, logits, n);
				const pr = new Float32Array(heap.buffer, p, n);
				const rows = Array.from(id, (v, i) => [v, lg[i], pr[i]]).sort((a, b) => b[1] - a[1]).slice(0, 2);
				rows.forEach((r, i) => { id[i] = r[0]; lg[i] = r[1]; pr[i] = 0.5; });
				const st = new Int32Array(heap.buffer, state, 2);
				last = { smpl, selected: st[0], sorted: st[1] };
				st[0] = 0;
				st[1] = 1;
				return rows.length;
			},
			_yzma_get_causal_attn: (ctx) => ctx === 3 ? 1 : 0,
			_yzma_get_model: (ctx) => ctx === 3 ? 9 : -2,
			_yzma_model_meta_arr_n: (model, key) => model !== 9 ? -2 : (arrays[str(key)] || { length: -1 }).length,
			_yzma_model_meta_arr_str: (model, key, i, buf, cap) => {
				const a = arrays[str(key)];
				if (!a || i < 0 || i >= a.length) return -1;
				return cap < 1024 ? -5 : copy(a[i], buf, cap);
			},
			_yzma_version: (buf, cap) => copy("b12728", buf, cap),
		},
		last: () => JSON.stringify(last),
	};
})();
`

func fakeABI15(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeABI15Source)
	helper := js.Global().Get("__yzmaABI15")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() {
		mod = previous
		for _, s := range []*scratch{&tokenScratch, &posScratch, &nSeqScratch, &seqScratch, &logitScratch,
			&orderScratch, &embdInScratch, &descScratch, &groupScratch, &applyScratch, &pieceScratch} {
			s.ptr, s.size = 0, 0
		}
	})

	return helper
}

func TestProcessDecisionGroup(t *testing.T) {
	helper := fakeABI15(t)

	bt, err := BatchExtInit(Context(7))
	if err != nil {
		t.Fatalf("BatchExtInit gave %v", err)
	}
	for i, tok := range []Token{4, 5, 6} {
		idx, _ := BatchExtAddToken(bt, 0, tok)
		BatchExtSetPos(bt, idx, Pos(i))
	}
	if err := BatchExtSetDecisionGroup(bt, 1, 1); err != nil {
		t.Fatalf("BatchExtSetDecisionGroup gave %v", err)
	}
	if err := BatchExtSetDecisionGroup(bt, 2, 2); err != nil {
		t.Fatalf("BatchExtSetDecisionGroup gave %v", err)
	}
	BatchExtSetOutputLogits(bt, 2, true)

	if err := BatchExtSetDecisionGroup(bt, 0, -1); !errors.Is(err, ErrBatchExtRejected) {
		t.Errorf("a negative group gave %v, want ErrBatchExtRejected", err)
	}
	if err := BatchExtSetDecisionGroup(bt, 3, 1); !errors.Is(err, ErrBatchExtRejected) {
		t.Errorf("an index past the end gave %v, want ErrBatchExtRejected", err)
	}

	if ret, err := Process(Context(7), ProcessTypeDecode, bt); err != nil || ret != 0 {
		t.Fatalf("Process gave %d, %v", ret, err)
	}
	want := `{"ctx":7,"type":1,"nPos":1,"tokens":[4,5,6],"pos":[0,1,2],"outputs":[0,0,1],` +
		`"order":null,"group":[0,1,2],"embd":null,"desc":null}`
	if got := helper.Call("last").String(); got != want {
		t.Errorf("the shim got\n%s\nwant\n%s", got, want)
	}

	BatchExtClear(bt)
	idx, _ := BatchExtAddToken(bt, 0, 4)
	BatchExtSetPos(bt, idx, 0)
	if _, err := Process(Context(7), ProcessTypeDecode, bt); err != nil {
		t.Fatalf("Process gave %v", err)
	}
	if got := helper.Call("last").String(); got != `"process"` {
		t.Errorf("a batch with no groups went to %s, want yzma_process", got)
	}
}

func TestProcessDecisionGroupEmbd(t *testing.T) {
	helper := fakeABI15(t)

	bt, _ := BatchExtInit(Context(7))
	idx, _ := BatchExtAddToken(bt, 0, 4)
	BatchExtSetPos(bt, idx, 0)
	idx, _ = BatchExtAddEmbd(bt, 0, []float32{1, 2}, 2)
	BatchExtSetPos(bt, idx, 1, 2)
	BatchExtSetDecisionGroup(bt, idx, 1)

	if _, err := Process(Context(7), ProcessTypeDecode, bt); err != nil {
		t.Fatalf("Process gave %v", err)
	}
	want := `{"ctx":7,"type":1,"nPos":4,"tokens":[4,-1],"pos":[0,0,0,0,1,2,0,0],"outputs":[0,0],` +
		`"order":null,"group":[0,1],"embd":[1,2],"desc":[0,0,0,0,0,2,2,0]}`
	if got := helper.Call("last").String(); got != want {
		t.Errorf("the shim got\n%s\nwant\n%s", got, want)
	}
}

func TestDecisionGroupOldModule(t *testing.T) {
	fakeABI14(t)

	bt, _ := BatchExtInit(Context(7))
	idx, _ := BatchExtAddToken(bt, 0, 4)
	if err := BatchExtSetDecisionGroup(bt, idx, 1); !errors.Is(err, ErrBatchExtNoDecisionGroup) {
		t.Errorf("BatchExtSetDecisionGroup gave %v, want ErrBatchExtNoDecisionGroup", err)
	}
}

func TestSamplerApply(t *testing.T) {
	helper := fakeABI15(t)

	data := []TokenData{{Id: 1, Logit: 0.5}, {Id: 2, Logit: 3}, {Id: 3, Logit: 1}, {Id: 4, Logit: -2}}
	cur := TokenDataArray{Data: &data[0], Size: uint64(len(data)), Selected: -1}
	SamplerApply(Sampler(5), &cur)

	if got := helper.Call("last").String(); got != `{"smpl":5,"selected":-1,"sorted":0}` {
		t.Errorf("the shim got %s", got)
	}
	if cur.Size != 2 || cur.Selected != 0 || cur.Sorted != 1 {
		t.Errorf("got size %d, selected %d, sorted %d, want 2, 0, 1", cur.Size, cur.Selected, cur.Sorted)
	}
	want := []TokenData{{Id: 2, Logit: 3, P: 0.5}, {Id: 3, Logit: 1, P: 0.5}}
	for i, w := range want {
		if data[i] != w {
			t.Errorf("candidate %d is %+v, want %+v", i, data[i], w)
		}
	}

	SamplerApply(Sampler(5), nil)
	SamplerApply(Sampler(5), &TokenDataArray{})
}

func TestContextAccessors(t *testing.T) {
	fakeABI15(t)

	if !GetCausalAttn(Context(3)) || GetCausalAttn(Context(4)) {
		t.Error("GetCausalAttn gave the wrong value")
	}
	if m := GetModel(Context(3)); m != 9 {
		t.Errorf("GetModel gave %d, want 9", m)
	}
	if m := GetModel(Context(4)); m != 0 {
		t.Errorf("GetModel of a bad context gave %d, want 0", m)
	}
	if v := Version(); v != "b12728" {
		t.Errorf("Version gave %q", v)
	}
}

func TestModelMetaArr(t *testing.T) {
	fakeABI15(t)

	const key = "tokenizer.tts_dicts"
	if n, err := ModelMetaArrN(Model(9), key); n != 2 || err != nil {
		t.Errorf("ModelMetaArrN gave %d, %v", n, err)
	}
	if n, err := ModelMetaArrN(Model(9), "missing"); n != -1 || err != nil {
		t.Errorf("ModelMetaArrN of a missing key gave %d, %v", n, err)
	}
	if _, err := ModelMetaArrN(Model(8), key); err == nil {
		t.Error("ModelMetaArrN of a bad model gave no error")
	}

	if s, ok, err := ModelMetaArrStr(Model(9), key, 1); s != "a much longer second element" || !ok || err != nil {
		t.Errorf("ModelMetaArrStr gave %q, %v, %v", s, ok, err)
	}
	if s, ok, err := ModelMetaArrStr(Model(9), key, 2); s != "" || ok || err != nil {
		t.Errorf("ModelMetaArrStr past the end gave %q, %v, %v", s, ok, err)
	}
}

func TestABI15CallsOnOlderModule(t *testing.T) {
	fakeABI14(t)

	if GetCausalAttn(Context(3)) || GetModel(Context(3)) != 0 || Version() != "" {
		t.Error("an older module gave a value")
	}
	if _, err := ModelMetaArrN(Model(9), "k"); !errors.Is(err, ErrModelNoMetaArr) {
		t.Errorf("ModelMetaArrN gave %v, want ErrModelNoMetaArr", err)
	}
	if _, _, err := ModelMetaArrStr(Model(9), "k", 0); !errors.Is(err, ErrModelNoMetaArr) {
		t.Errorf("ModelMetaArrStr gave %v, want ErrModelNoMetaArr", err)
	}

	data := []TokenData{{Id: 1, Logit: 1}}
	cur := TokenDataArray{Data: &data[0], Size: 1, Selected: -1}
	SamplerApply(Sampler(5), &cur)
	if cur.Size != 1 || cur.Selected != -1 {
		t.Error("SamplerApply changed the candidates on an older module")
	}
}

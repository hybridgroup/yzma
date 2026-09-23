//go:build js && wasm

package llamawasm

import (
	"syscall/js"
	"testing"
)

// The fake module of ABI version 9 has a real heap, thus the tests cover the
// strings, the doubles, and the buffers of state that cross the boundary.
const fakeABI9Source = `
globalThis.__yzmaABI9 = (function () {
	const heap = new Uint8Array(1 << 20);
	let next = 8;
	const meta = [["general.name", "Tiny"], ["general.license", "x".repeat(40000)]];
	const state = new Uint8Array([1, 2, 3, 4, 5]);
	let restored = null;
	let lastNoPerf = -1;

	function str(s, buf, cap) {
		const b = new TextEncoder().encode(s);
		heap.set(b.subarray(0, Math.max(0, cap - 1)), buf);
		heap[buf + Math.min(b.length, cap - 1)] = 0;
		return b.length;
	}
	function copy(s, buf, cap) {
		if (cap < s.length + 1) return -5;
		return str(s, buf, cap);
	}
	function cstr(p) {
		let e = p;
		while (heap[e] !== 0) e++;
		return new TextDecoder().decode(heap.subarray(p, e));
	}
	function doubles(out, values) {
		new Float64Array(heap.buffer, out, values.length).set(values);
		return values.length;
	}

	return {
		module: {
			HEAPU8: heap,
			_malloc: (n) => { const p = next; next += n + (8 - (n % 8)); return p; },
			_free: () => {},
			_yzma_abi_version: () => 9,
			_yzma_last_error: () => 0,
			_yzma_context_new_ext: (m, nctx, nb, nub, nt, e, pool, nseq, noPerf) => { lastNoPerf = noPerf; return 4; },
			_yzma_model_meta_count: () => meta.length,
			_yzma_model_meta_key_by_index: (m, i, buf, cap) => i < meta.length ? str(meta[i][0], buf, cap) : -1,
			_yzma_model_meta_val_str_by_index: (m, i, buf, cap) => i < meta.length ? str(meta[i][1], buf, cap) : -1,
			_yzma_model_meta_val_str: (m, key, buf, cap) => {
				const kv = meta.find((p) => p[0] === cstr(key));
				return kv ? str(kv[1], buf, cap) : -1;
			},
			_yzma_model_meta_key_str: (k, buf, cap) => k === 1 ? copy("general.sampling.top_k", buf, cap) : 0,
			_yzma_model_cls_label: (m, i, buf, cap) => copy("label" + i, buf, cap),
			_yzma_state_get_size: () => state.length,
			_yzma_state_get_data: (c, dst, size) => { heap.set(state.subarray(0, size), dst); return Math.min(size, state.length); },
			_yzma_state_set_data: (c, src, size) => { restored = Array.from(heap.subarray(src, src + size)); return size; },
			_yzma_state_seq_get_size_ext: (c, seq, flags) => seq * 10 + flags,
			_yzma_state_seq_get_data_ext: (c, dst, size, seq, flags) => { heap.set([seq, flags], dst); return 2; },
			_yzma_state_seq_set_data_ext: (c, src, size, seq, flags) => { restored = [seq, flags].concat(Array.from(heap.subarray(src, src + size))); return size; },
			_yzma_perf_context: (c, out) => c === 1 ? doubles(out, [10.5, 20.25, 30, 40, 7, 3, 2]) : -2,
			_yzma_perf_context_reset: (c) => c === 1 ? 0 : -2,
			_yzma_perf_sampler: (s, out) => s === 1 ? doubles(out, [1.5, 9]) : -2,
			_yzma_perf_sampler_reset: () => 0,
			_yzma_print_system_info: (buf, cap) => copy("CPU : WASM_SIMD = 1 | ", buf, cap),
			_yzma_ftype_name: (f, buf, cap) => copy(f === 15 ? "Q4_K - Medium" : "unknown", buf, cap),
			_yzma_time_us: () => 5000000000123,
			_yzma_max_devices: () => 16,
			_yzma_max_parallel_sequences: () => 256,
			_yzma_supports_gpu_offload: () => 1,
		},
		restored: () => restored,
		lastNoPerf: () => lastNoPerf,
	};
})();
`

func fakeABI9(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeABI9Source)
	helper := js.Global().Get("__yzmaABI9")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() {
		mod = previous
		for _, s := range []*scratch{&pieceScratch, &outScratch} {
			s.ptr, s.size = 0, 0
		}
	})

	return helper
}

// fakeOld puts a module in place that has none of the calls of ABI version 9.
func fakeOld(t *testing.T) {
	t.Helper()

	js.Global().Call("eval", "globalThis.__yzmaOld8 = { HEAPU8: new Uint8Array(64), _malloc: () => 8, _free: () => {} };")

	previous := mod
	mod = js.Global().Get("__yzmaOld8")
	t.Cleanup(func() { mod = previous })
}

//go:build js && wasm

package llamawasm

import (
	"syscall/js"
	"testing"
)

const fakeTemplateSource = `
globalThis.__yzmaTemplate = (function () {
	const heap = new Uint8Array(1 << 20);
	let next = 8;
	const meta = { "tokenizer.chat_template.systemone": "State: {{ state }}" };

	function str(s, buf, cap) {
		const b = new TextEncoder().encode(s);
		heap.set(b.subarray(0, Math.max(0, cap - 1)), buf);
		heap[buf + Math.min(b.length, cap - 1)] = 0;
		return b.length;
	}
	function cstr(p) {
		let e = p;
		while (heap[e] !== 0) e++;
		return new TextDecoder().decode(heap.subarray(p, e));
	}

	return {
		HEAPU8: heap,
		_malloc: (n) => { const p = next; next += n + (8 - (n % 8)); return p; },
		_free: () => {},
		_yzma_abi_version: () => 9,
		_yzma_last_error: () => 0,
		_yzma_model_chat_template: (m, buf, cap) => str("default", buf, cap),
		_yzma_model_meta_val_str: (m, key, buf, cap) => {
			const v = meta[cstr(key)];
			return v === undefined ? -1 : str(v, buf, cap);
		},
	};
})();
`

func TestModelChatTemplateNamed(t *testing.T) {
	js.Global().Call("eval", fakeTemplateSource)

	previous := mod
	mod = js.Global().Get("__yzmaTemplate")
	t.Cleanup(func() {
		mod = previous
		pieceScratch.ptr, pieceScratch.size = 0, 0
	})

	m := Model(1)
	if got := ModelChatTemplate(m, ""); got != "default" {
		t.Errorf("ModelChatTemplate with no name gave %q, want the default template", got)
	}
	if got := ModelChatTemplate(m, "systemone"); got != "State: {{ state }}" {
		t.Errorf("ModelChatTemplate with a name gave %q, want the named template", got)
	}
	if got := ModelChatTemplate(m, "tool_use"); got != "" {
		t.Errorf("ModelChatTemplate with a missing name gave %q, want an empty string", got)
	}
}

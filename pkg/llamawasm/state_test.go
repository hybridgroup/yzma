//go:build js && wasm

package llamawasm

import (
	"bytes"
	"testing"
)

func TestStateGetAndSet(t *testing.T) {
	helper := fakeABI9(t)

	ctx := Context(1)
	size := StateGetSize(ctx)
	if size != 5 {
		t.Fatalf("StateGetSize gave %d, want 5", size)
	}

	buf := make([]byte, size)
	if n := StateGetData(ctx, buf); n != 5 || !bytes.Equal(buf, []byte{1, 2, 3, 4, 5}) {
		t.Errorf("StateGetData gave %d and %v", n, buf)
	}

	if n := StateSetData(ctx, []byte{9, 8, 7}); n != 3 {
		t.Errorf("StateSetData gave %d, want 3", n)
	}
	if got := helper.Call("restored"); got.Length() != 3 || got.Index(0).Int() != 9 {
		t.Errorf("the shim got %v, want [9 8 7]", got)
	}
}

func TestStateSeq(t *testing.T) {
	helper := fakeABI9(t)

	ctx := Context(1)
	if got := StateSeqGetSize(ctx, 2); got != 20 {
		t.Errorf("StateSeqGetSize gave %d, want 20", got)
	}
	if got := StateSeqGetSizeExt(ctx, 2, 1); got != 21 {
		t.Errorf("StateSeqGetSizeExt gave %d, want 21", got)
	}

	buf := make([]byte, 4)
	if n := StateSeqGetDataExt(ctx, buf, 3, 1); n != 2 || buf[0] != 3 || buf[1] != 1 {
		t.Errorf("StateSeqGetDataExt gave %d and %v", n, buf)
	}

	if n := StateSeqSetData(ctx, []byte{42}, 5); n != 1 {
		t.Errorf("StateSeqSetData gave %d, want 1", n)
	}
	got := helper.Call("restored")
	if got.Index(0).Int() != 5 || got.Index(1).Int() != 0 || got.Index(2).Int() != 42 {
		t.Errorf("the shim got %v, want [5 0 42]", got)
	}
}

func TestStateOldModule(t *testing.T) {
	fakeOld(t)

	if StateGetSize(Context(1)) != 0 || StateGetData(Context(1), make([]byte, 4)) != 0 ||
		StateSeqSetData(Context(1), []byte{1}, 0) != 0 {
		t.Error("a module of an earlier version must give 0")
	}
}

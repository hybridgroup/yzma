package mtmd

import (
	_ "image/jpeg"
	"testing"
	"unsafe"
)

func TestBitmap(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	data, x, y, err := openFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("count not open file")
	}

	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	if BitmapGetNBytes(bitmap) != 2073600 {
		t.Fatal("unable to open bitmap")
	}
}

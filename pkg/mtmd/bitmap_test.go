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

func TestBitmapGetNxAndNy(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	data, x, y, err := openFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open file")
	}

	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	nx := BitmapGetNx(bitmap)
	ny := BitmapGetNy(bitmap)

	if nx != x || ny != y {
		t.Fatalf("BitmapGetNx or BitmapGetNy returned incorrect dimensions: nx=%d, ny=%d", nx, ny)
	}

	t.Logf("BitmapGetNx returned: %d, BitmapGetNy returned: %d", nx, ny)
}

func TestBitmapGetData(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	data, x, y, err := openFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open file")
	}

	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	rawData := BitmapGetData(bitmap)
	if rawData == nil || len(rawData) != int(x*y*3) {
		t.Fatalf("BitmapGetData returned incorrect data size: %d", len(rawData))
	}

	t.Logf("BitmapGetData returned data of size: %d", len(rawData))
}

func TestBitmapIsAudio(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	data, x, y, err := openFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open file")
	}

	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	isAudio := BitmapIsAudio(bitmap)
	if isAudio {
		t.Fatal("BitmapIsAudio incorrectly returned true for an image bitmap")
	}

	t.Logf("BitmapIsAudio returned: %v", isAudio)
}

func TestBitmapGetAndSetId(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	data, x, y, err := openFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open file")
	}

	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	id := "test_bitmap_id"
	BitmapSetId(bitmap, id)

	retrievedId := BitmapGetId(bitmap)
	if retrievedId != id {
		t.Fatalf("BitmapGetId returned incorrect ID: %s, expected: %s", retrievedId, id)
	}

	t.Logf("BitmapGetId returned: %s", retrievedId)
}

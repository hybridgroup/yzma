package mtmd

import (
	_ "image/jpeg"
	"os"
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestBitmap(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
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

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
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

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
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

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
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

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
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

func TestBitmapInitFromAudio(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// Create dummy audio data (e.g., 16000 samples for 1 second at 16kHz)
	nSamples := uint64(16000)
	audioData := make([]float32, nSamples)
	for i := range audioData {
		audioData[i] = float32(i) / float32(nSamples) // simple ramp
	}

	bitmap := BitmapInitFromAudio(nSamples, &audioData[0])
	defer BitmapFree(bitmap)

	if bitmap == Bitmap(0) {
		t.Fatal("BitmapInitFromAudio returned an invalid bitmap")
	}

	if !BitmapIsAudio(bitmap) {
		t.Fatal("BitmapIsAudio returned false for audio bitmap")
	}

	t.Logf("BitmapInitFromAudio created bitmap with %d samples", nSamples)
}

func TestBitmapInitPlaceholder(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// BitmapInit with data=0 creates a placeholder bitmap.
	bitmap := BitmapInit(640, 480, 0)
	defer BitmapFree(bitmap)

	if bitmap == Bitmap(0) {
		t.Fatal("BitmapInit returned an invalid bitmap for placeholder")
	}

	if BitmapGetNx(bitmap) != 640 {
		t.Fatalf("BitmapGetNx returned unexpected value: %d", BitmapGetNx(bitmap))
	}
	if BitmapGetNy(bitmap) != 480 {
		t.Fatalf("BitmapGetNy returned unexpected value: %d", BitmapGetNy(bitmap))
	}
}

func TestBitmapIsPlaceholder(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// Placeholder bitmap created with nil data.
	placeholder := BitmapInit(640, 480, 0)
	defer BitmapFree(placeholder)

	if !BitmapIsPlaceholder(placeholder) {
		t.Fatal("BitmapIsPlaceholder returned false for a placeholder bitmap")
	}
}

func TestBitmapIsPlaceholderForRegularBitmap(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open file")
	}

	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	if BitmapIsPlaceholder(bitmap) {
		t.Fatal("BitmapIsPlaceholder returned true for a regular (non-placeholder) bitmap")
	}
}

func TestBitmapIsPlaceholderReturnsFalseForZero(t *testing.T) {
	// No library load needed; BitmapIsPlaceholder short-circuits on zero handle.
	if BitmapIsPlaceholder(Bitmap(0)) {
		t.Fatal("BitmapIsPlaceholder should return false for a zero (nil) Bitmap handle")
	}
}

func TestBitmapGetNBytesPlaceholder(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	bitmap := BitmapInit(640, 480, 0)
	defer BitmapFree(bitmap)

	if BitmapGetNBytes(bitmap) != 0 {
		t.Fatalf("BitmapGetNBytes should return 0 for placeholder, got %d", BitmapGetNBytes(bitmap))
	}
}

func TestBitmapGetDataPlaceholder(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	bitmap := BitmapInit(640, 480, 0)
	defer BitmapFree(bitmap)

	if BitmapGetData(bitmap) != nil {
		t.Fatal("BitmapGetData should return nil for a placeholder bitmap")
	}
}

func TestBitmapInitFromAudioPlaceholder(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	nSamples := uint64(16000)
	bitmap := BitmapInitFromAudio(nSamples, nil)
	defer BitmapFree(bitmap)

	if bitmap == Bitmap(0) {
		t.Fatal("BitmapInitFromAudio returned an invalid bitmap for placeholder")
	}

	if !BitmapIsAudio(bitmap) {
		t.Fatal("BitmapIsAudio returned false for audio placeholder bitmap")
	}

	if !BitmapIsPlaceholder(bitmap) {
		t.Fatal("BitmapIsPlaceholder returned false for audio placeholder bitmap")
	}

	if BitmapGetNx(bitmap) != uint32(nSamples) {
		t.Fatalf("BitmapGetNx returned unexpected value for audio placeholder: %d", BitmapGetNx(bitmap))
	}
}

func testSetupMtmdCtx(t *testing.T) (llama.Model, Context) {
	t.Helper()
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}

	ctx, err := InitFromFile(mmprojFile, model, ContextParamsDefault())
	if err != nil {
		llama.ModelFree(model)
		t.Fatalf("InitFromFile failed: %v", err)
	}
	return model, ctx
}

func TestBitmapInitFromFile(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model, ctx := testSetupMtmdCtx(t)
	defer llama.ModelFree(model)
	defer Free(ctx)

	bitmap := BitmapInitFromFile(ctx, "../../images/domestic_llama.jpg", false)
	defer BitmapFree(bitmap.Bitmap)

	if bitmap.Bitmap == Bitmap(0) {
		t.Fatal("BitmapInitFromFile returned an invalid bitmap")
	}

	if BitmapGetNBytes(bitmap.Bitmap) == 0 {
		t.Fatal("BitmapInitFromFile returned a bitmap with zero bytes")
	}

	t.Logf("BitmapInitFromFile created bitmap with %d bytes", BitmapGetNBytes(bitmap.Bitmap))
}

func TestBitmapInitFromFilePlaceholder(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model, ctx := testSetupMtmdCtx(t)
	defer llama.ModelFree(model)
	defer Free(ctx)

	bitmap := BitmapInitFromFile(ctx, "../../images/domestic_llama.jpg", true)
	defer BitmapFree(bitmap.Bitmap)

	if bitmap.Bitmap == Bitmap(0) {
		t.Fatal("BitmapInitFromFile (placeholder) returned an invalid bitmap")
	}

	if !BitmapIsPlaceholder(bitmap.Bitmap) {
		t.Fatal("BitmapIsPlaceholder returned false for a file-loaded placeholder bitmap")
	}

	if BitmapGetNx(bitmap.Bitmap) == 0 || BitmapGetNy(bitmap.Bitmap) == 0 {
		t.Fatal("BitmapInitFromFile placeholder has zero dimensions")
	}

	t.Logf("BitmapInitFromFile placeholder: nx=%d, ny=%d", BitmapGetNx(bitmap.Bitmap), BitmapGetNy(bitmap.Bitmap))
}

func TestBitmapInitFromFileZeroCtx(t *testing.T) {
	// BitmapInitFromFile should return a zero Bitmap when ctx is 0.
	bitmap := BitmapInitFromFile(Context(0), "../../images/domestic_llama.jpg", false)
	if bitmap.Bitmap != Bitmap(0) {
		t.Fatal("BitmapInitFromFile should return zero Bitmap for zero context")
	}
}

func TestBitmapInitFromFileNonExistent(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model, ctx := testSetupMtmdCtx(t)
	defer llama.ModelFree(model)
	defer Free(ctx)

	bitmap := BitmapInitFromFile(ctx, "/nonexistent/path/image.jpg", false)
	if bitmap.Bitmap != Bitmap(0) {
		BitmapFree(bitmap.Bitmap)
		t.Fatal("BitmapInitFromFile should return zero Bitmap for non-existent file")
	}
}

func TestBitmapInitFromBuf(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model, ctx := testSetupMtmdCtx(t)
	defer llama.ModelFree(model)
	defer Free(ctx)

	fileBytes, err := os.ReadFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatalf("could not read image file: %v", err)
	}

	bitmap := BitmapInitFromBuf(ctx, &fileBytes[0], uint64(len(fileBytes)), false)
	defer BitmapFree(bitmap.Bitmap)

	if bitmap.Bitmap == Bitmap(0) {
		t.Fatal("BitmapInitFromBuf returned an invalid bitmap")
	}

	if BitmapGetNBytes(bitmap.Bitmap) == 0 {
		t.Fatal("BitmapInitFromBuf returned a bitmap with zero bytes")
	}

	t.Logf("BitmapInitFromBuf created bitmap with %d bytes", BitmapGetNBytes(bitmap.Bitmap))
}

func TestBitmapInitFromBufPlaceholder(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model, ctx := testSetupMtmdCtx(t)
	defer llama.ModelFree(model)
	defer Free(ctx)

	fileBytes, err := os.ReadFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatalf("could not read image file: %v", err)
	}

	bitmap := BitmapInitFromBuf(ctx, &fileBytes[0], uint64(len(fileBytes)), true)
	defer BitmapFree(bitmap.Bitmap)

	if bitmap.Bitmap == Bitmap(0) {
		t.Fatal("BitmapInitFromBuf (placeholder) returned an invalid bitmap")
	}

	if !BitmapIsPlaceholder(bitmap.Bitmap) {
		t.Fatal("BitmapIsPlaceholder returned false for a buf-loaded placeholder bitmap")
	}

	if BitmapGetNx(bitmap.Bitmap) == 0 || BitmapGetNy(bitmap.Bitmap) == 0 {
		t.Fatal("BitmapInitFromBuf placeholder has zero dimensions")
	}

	t.Logf("BitmapInitFromBuf placeholder: nx=%d, ny=%d", BitmapGetNx(bitmap.Bitmap), BitmapGetNy(bitmap.Bitmap))
}

func TestBitmapInitFromBufZeroCtx(t *testing.T) {
	fileBytes := []byte{0xFF, 0xD8, 0xFF} // minimal JPEG header bytes
	bitmap := BitmapInitFromBuf(Context(0), &fileBytes[0], uint64(len(fileBytes)), false)
	if bitmap.Bitmap != Bitmap(0) {
		t.Fatal("BitmapInitFromBuf should return zero Bitmap for zero context")
	}
}

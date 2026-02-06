package mtmd

import (
	"image"
	"os"
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func testModelFileName(t *testing.T) string {
	if os.Getenv("YZMA_TEST_MMMODEL") == "" {
		t.Skip("no YZMA_TEST_MMMODEL skipping test")
	}

	return os.Getenv("YZMA_TEST_MMMODEL")
}

func testMMProjFileName(t *testing.T) string {
	if os.Getenv("YZMA_TEST_MMPROJ") == "" {
		t.Skip("no YZMA_TEST_MMPROJ skipping test")
	}

	return os.Getenv("YZMA_TEST_MMPROJ")
}

func testSetup(t *testing.T) {
	if os.Getenv("YZMA_LIB") == "" {
		t.Fatal("no YZMA_LIB set for tests")
	}
	testPath := os.Getenv("YZMA_LIB")

	if err := llama.Load(testPath); err != nil {
		t.Fatal("unable to load library", err.Error())
	}
	if err := Load(testPath); err != nil {
		t.Fatal("unable to load library", err.Error())
	}

	llama.Init()
}

func testCleanup(t *testing.T) {
	llama.BackendFree()
}

func testSetupChunks(t *testing.T, ctx Context, chunks InputChunks) {
	text := NewInputText("Here is an image: <__media__>", true, true)

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open image file")
	}

	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	bitmaps := []Bitmap{bitmap} // Replace with actual bitmap data if available

	result := Tokenize(ctx, chunks, text, bitmaps)
	if result != 0 {
		t.Fatalf("Tokenize failed with result: %d", result)
	}
}

func openImageFile(path string) ([]byte, uint32, uint32, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, 0, err
	}

	// Get the image bounds
	bounds := img.Bounds()
	width := uint32(bounds.Dx())
	height := uint32(bounds.Dy())

	// Create a slice to hold the RGB data
	rgbData := make([]byte, 0, width*height*3)

	// Extract RGB data
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgbData = append(rgbData, byte(r>>8), byte(g>>8), byte(b>>8))
		}
	}

	return rgbData, width, height, nil
}

func benchmarkModelFileName(b *testing.B) string {
	if os.Getenv("YZMA_BENCHMARK_MMMODEL") == "" {
		b.Skip("no YZMA_BENCHMARK_MMMODEL skipping test")
	}

	return os.Getenv("YZMA_BENCHMARK_MMMODEL")
}

func benchmarkProjectorFileName(b *testing.B) string {
	if os.Getenv("YZMA_BENCHMARK_MMPROJ") == "" {
		b.Skip("no YZMA_BENCHMARK_MMPROJ skipping test")
	}

	return os.Getenv("YZMA_BENCHMARK_MMPROJ")
}

func benchmarkSetup(b *testing.B) {
	if os.Getenv("YZMA_LIB") == "" {
		b.Fatal("no YZMA_LIB set for tests")
	}
	testPath := os.Getenv("YZMA_LIB")

	if err := llama.Load(testPath); err != nil {
		b.Fatal("unable to load library", err.Error())
	}
	if err := Load(testPath); err != nil {
		b.Fatal("unable to load library", err.Error())
	}

	llama.LogSet(llama.LogSilent())
	LogSet(llama.LogSilent())

	llama.Init()
}

func benchmarkCleanup(b *testing.B) {
	llama.LogSet(llama.LogNormal)
	LogSet(llama.LogNormal)

	llama.Close()
}

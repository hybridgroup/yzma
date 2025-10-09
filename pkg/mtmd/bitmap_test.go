package mtmd

import (
	"image"
	_ "image/jpeg"
	"os"
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
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

func testSetup(t *testing.T) {
	testPath := "../../lib"

	if err := llama.Load(testPath); err != nil {
		t.Fatal("unable to load library", err.Error())
	}
	if err := Load(testPath); err != nil {
		t.Fatal("unable to load library", err.Error())
	}

	llama.BackendInit()
}

func testCleanup(t *testing.T) {
	llama.BackendFree()
}

func openFile(path string) ([]byte, uint32, uint32, error) {
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

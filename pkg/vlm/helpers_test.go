package vlm

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
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
	if err := mtmd.Load(testPath); err != nil {
		t.Fatal("unable to load library", err.Error())
	}

	llama.Init()
}

func testCleanup(t *testing.T) {
	llama.BackendFree()
}

func openImageFile(path string) ([]byte, uint32, uint32, error) {
	fmt.Println(path)
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

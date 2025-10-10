package mtmd

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

// // if bitmap is image:
// //     length of data must be nx * ny * 3
// //     the data is in RGBRGBRGB... format
// MTMD_API uint32_t              mtmd_bitmap_get_nx     (const mtmd_bitmap * bitmap);
// MTMD_API uint32_t              mtmd_bitmap_get_ny     (const mtmd_bitmap * bitmap);
// MTMD_API const unsigned char * mtmd_bitmap_get_data   (const mtmd_bitmap * bitmap);
// // bitmap ID is optional, but useful for KV cache tracking
// // these getters/setters are dedicated functions, so you can for example calculate the hash of the image based on mtmd_bitmap_get_data()
// MTMD_API const char * mtmd_bitmap_get_id(const mtmd_bitmap * bitmap);
// MTMD_API void         mtmd_bitmap_set_id(mtmd_bitmap * bitmap, const char * id);

// Opaque types (represented as pointers)
type Bitmap uintptr

var (
	// MTMD_API mtmd_bitmap *         mtmd_bitmap_init           (uint32_t nx, uint32_t ny, const unsigned char * data);
	bitmapInitFunc ffi.Fun

	// MTMD_API void                  mtmd_bitmap_free       (mtmd_bitmap * bitmap);
	bitmapFreeFunc ffi.Fun

	// MTMD_API size_t                mtmd_bitmap_get_n_bytes(const mtmd_bitmap * bitmap);
	bitmapGetNBytesFunc ffi.Fun

	// MTMD_API mtmd_bitmap * mtmd_helper_bitmap_init_from_file(mtmd_context * ctx, const char * fname);
	bitmapInitFromFileFunc ffi.Fun

	// MTMD_API mtmd_bitmap * mtmd_helper_bitmap_init_from_buf(mtmd_context * ctx, const unsigned char * buf, size_t len);
	bitmapInitFromBufFunc ffi.Fun
)

func loadBitmapFuncs(lib ffi.Lib) error {
	var err error

	if bitmapInitFunc, err = lib.Prep("mtmd_bitmap_init", &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("mtmd_bitmap_init", err)
	}

	if bitmapFreeFunc, err = lib.Prep("mtmd_bitmap_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("mtmd_bitmap_free", err)
	}

	if bitmapGetNBytesFunc, err = lib.Prep("mtmd_bitmap_get_n_bytes", &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("mtmd_bitmap_get_n_bytes", err)
	}

	if bitmapInitFromFileFunc, err = lib.Prep("mtmd_helper_bitmap_init_from_file", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("mtmd_helper_bitmap_init_from_file", err)
	}

	if bitmapInitFromBufFunc, err = lib.Prep("mtmd_helper_bitmap_init_from_buf", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("mtmd_helper_bitmap_init_from_buf", err)
	}

	return nil
}

// BitmapInit initializes a Bitmap.
func BitmapInit(nx uint32, ny uint32, data uintptr) Bitmap {
	var bitmap Bitmap
	bitmapInitFunc.Call(unsafe.Pointer(&bitmap), &nx, &ny, unsafe.Pointer(&data))

	return bitmap
}

// BitmapFree frees a previous initialized Bitmap.
func BitmapFree(bitmap Bitmap) {
	bitmapFreeFunc.Call(nil, unsafe.Pointer(&bitmap))
}

// BitmapGetNBytes returns the number of bytes in the Bitmap.
func BitmapGetNBytes(bitmap Bitmap) uint32 {
	var result ffi.Arg
	bitmapGetNBytesFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&bitmap))

	return uint32(result)
}

// BitmapInitFromFile initializes a Bitmap from a file.
func BitmapInitFromFile(ctx Context, fname string) Bitmap {
	var bitmap Bitmap
	file := &[]byte(fname + "\x00")[0]
	bitmapInitFromFileFunc.Call(unsafe.Pointer(&bitmap), unsafe.Pointer(&ctx), unsafe.Pointer(&file))

	return bitmap
}

// BitmapInitFromBuf initializes a Bitmap from a buffer of bytes.
func BitmapInitFromBuf(ctx Context, buf *byte, len uint64) Bitmap {
	var bitmap Bitmap
	bitmapInitFromBufFunc.Call(unsafe.Pointer(&bitmap), unsafe.Pointer(&ctx), unsafe.Pointer(&buf), &len)

	return bitmap
}

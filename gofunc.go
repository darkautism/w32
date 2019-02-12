package w32

import (
	"errors"
	"image"
	"unsafe"
)

//GoHBITMAP is a HBITMAP wrapper
type GoHBITMAP struct {
	Hbmp   HBITMAP
	W      int
	H      int
	Pixels unsafe.Pointer
}

//Image Convert GoHBITMAP to Image type
func (ghbmp *GoHBITMAP) Image() *image.RGBA {
	rect := image.Rect(0, 0, ghbmp.W, ghbmp.H)
	i := 0
	img := image.NewRGBA(rect)
	pixels := (*[9999999]uint8)(ghbmp.Pixels)[:ghbmp.W*ghbmp.H*4]
	for y := 0; y < ghbmp.H; y++ {
		for x := 0; x < ghbmp.W; x++ {
			v0 := pixels[i+0]
			v1 := pixels[i+1]
			v2 := pixels[i+2]
			// BGRA => RGBA, and set A to 255
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = v2, v1, v0, 255
			i += 4
		}
	}
	return img
}

//Delete this GoHBITMAP and release resource
func (ghbmp *GoHBITMAP) Delete() {
	ghbmp.H = -1
	ghbmp.W = -1
	DeleteObject(HGDIOBJ(ghbmp.Hbmp))
}

//ScreenShot can take snapshot of screen
func ScreenShot(h HWND, x, y, width, height int) (*GoHBITMAP, error) {
	hdc := GetDC(h)
	if hdc == 0 {
		return nil, errors.New("GetDC failed")
	}
	defer ReleaseDC(h, hdc)
	hDCMem := CreateCompatibleDC(hdc)
	if hDCMem == 0 {
		return nil, errors.New("CreateCompatibleDC failed")
	}
	defer DeleteDC(hDCMem)

	var bmi BITMAPINFO
	bmi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bmi.BmiHeader))
	bmi.BmiHeader.BiPlanes = 1
	bmi.BmiHeader.BiBitCount = 32
	bmi.BmiHeader.BiWidth = int32(width)
	bmi.BmiHeader.BiHeight = int32(-height)
	bmi.BmiHeader.BiCompression = BI_RGB
	bmi.BmiHeader.BiSizeImage = 0

	var p unsafe.Pointer
	bitmap := CreateDIBSection(hdc, &bmi, DIB_RGB_COLORS, &p, HANDLE(0), 0)
	old := SelectObject(hDCMem, HGDIOBJ(bitmap))
	defer SelectObject(hDCMem, old)

	BitBlt(hDCMem, int32(x), int32(y), int32(width), int32(height), hdc, 0, 0, SRCCOPY)
	return &GoHBITMAP{
		Hbmp:   bitmap,
		W:      width,
		H:      height,
		Pixels: p,
	}, nil
}

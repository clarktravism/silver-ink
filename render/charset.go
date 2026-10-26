package render

import (
	"bytes"
	"errors"
	"image"
	_ "image/png"
	"os"

	"github.com/Zyko0/go-sdl3/sdl"
)

var charset *sdl.Texture

func LoadCharset(r *sdl.Renderer, filename string) error {
	tex, err := readTexturePNG(r, filename)
	if err != nil {
		return err
	}
	if charset != nil {
		charset.Destroy()
	}
	charset = tex
	return nil
}

func DestroyCharsets() {
	if charset != nil {
		charset.Destroy()
	}
}

const offPixel = uint32(0)
const onPixel = uint32(0xFFFFFFFF)

func readTexturePNG(r *sdl.Renderer, filename string) (*sdl.Texture, error) {

	// read and validate chr image file
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	size := img.Bounds().Size()
	if size.X != 128 || size.Y != 128 {
		return nil, errors.New("incorrect chr image size")
	}

	// create surface
	surface, _ := sdl.CreateSurface(128, 128, sdl.PIXELFORMAT_RGBA32)
	defer surface.Destroy()
	pixels := surface.Pixels()
	//sPixels := unsafe.Slice((*uint32)(pixels), 128*128)

	// set pixels
	i := 0
	for iy := range 128 {
		for ix := range 128 {
			r, g, b, _ := img.At(ix, iy).RGBA()
			if r == 0 && g == 0 && b == 0 {
				pixels[i] = 0
				pixels[i+1] = 0
				pixels[i+2] = 0
				pixels[i+3] = 0
			} else {
				pixels[i] = 255
				pixels[i+1] = 255
				pixels[i+2] = 255
				pixels[i+3] = 255
			}
			i+=4
		}
	}

	texture, _ := r.CreateTextureFromSurface(surface)
	texture.SetScaleMode(sdl.SCALEMODE_PIXELART)
	return texture, nil
}

var charLookUp = [256]sdl.FRect{
	{X: 0, Y: 0, W: 8, H: 8},
	{X: 8, Y: 0, W: 8, H: 8},
	{X: 16, Y: 0, W: 8, H: 8},
	{X: 24, Y: 0, W: 8, H: 8},
	{X: 32, Y: 0, W: 8, H: 8},
	{X: 40, Y: 0, W: 8, H: 8},
	{X: 48, Y: 0, W: 8, H: 8},
	{X: 56, Y: 0, W: 8, H: 8},
	{X: 64, Y: 0, W: 8, H: 8},
	{X: 72, Y: 0, W: 8, H: 8},
	{X: 80, Y: 0, W: 8, H: 8},
	{X: 88, Y: 0, W: 8, H: 8},
	{X: 96, Y: 0, W: 8, H: 8},
	{X: 104, Y: 0, W: 8, H: 8},
	{X: 112, Y: 0, W: 8, H: 8},
	{X: 120, Y: 0, W: 8, H: 8},
	{X: 0, Y: 8, W: 8, H: 8},
	{X: 8, Y: 8, W: 8, H: 8},
	{X: 16, Y: 8, W: 8, H: 8},
	{X: 24, Y: 8, W: 8, H: 8},
	{X: 32, Y: 8, W: 8, H: 8},
	{X: 40, Y: 8, W: 8, H: 8},
	{X: 48, Y: 8, W: 8, H: 8},
	{X: 56, Y: 8, W: 8, H: 8},
	{X: 64, Y: 8, W: 8, H: 8},
	{X: 72, Y: 8, W: 8, H: 8},
	{X: 80, Y: 8, W: 8, H: 8},
	{X: 88, Y: 8, W: 8, H: 8},
	{X: 96, Y: 8, W: 8, H: 8},
	{X: 104, Y: 8, W: 8, H: 8},
	{X: 112, Y: 8, W: 8, H: 8},
	{X: 120, Y: 8, W: 8, H: 8},
	{X: 0, Y: 16, W: 8, H: 8},
	{X: 8, Y: 16, W: 8, H: 8},
	{X: 16, Y: 16, W: 8, H: 8},
	{X: 24, Y: 16, W: 8, H: 8},
	{X: 32, Y: 16, W: 8, H: 8},
	{X: 40, Y: 16, W: 8, H: 8},
	{X: 48, Y: 16, W: 8, H: 8},
	{X: 56, Y: 16, W: 8, H: 8},
	{X: 64, Y: 16, W: 8, H: 8},
	{X: 72, Y: 16, W: 8, H: 8},
	{X: 80, Y: 16, W: 8, H: 8},
	{X: 88, Y: 16, W: 8, H: 8},
	{X: 96, Y: 16, W: 8, H: 8},
	{X: 104, Y: 16, W: 8, H: 8},
	{X: 112, Y: 16, W: 8, H: 8},
	{X: 120, Y: 16, W: 8, H: 8},
	{X: 0, Y: 24, W: 8, H: 8},
	{X: 8, Y: 24, W: 8, H: 8},
	{X: 16, Y: 24, W: 8, H: 8},
	{X: 24, Y: 24, W: 8, H: 8},
	{X: 32, Y: 24, W: 8, H: 8},
	{X: 40, Y: 24, W: 8, H: 8},
	{X: 48, Y: 24, W: 8, H: 8},
	{X: 56, Y: 24, W: 8, H: 8},
	{X: 64, Y: 24, W: 8, H: 8},
	{X: 72, Y: 24, W: 8, H: 8},
	{X: 80, Y: 24, W: 8, H: 8},
	{X: 88, Y: 24, W: 8, H: 8},
	{X: 96, Y: 24, W: 8, H: 8},
	{X: 104, Y: 24, W: 8, H: 8},
	{X: 112, Y: 24, W: 8, H: 8},
	{X: 120, Y: 24, W: 8, H: 8},
	{X: 0, Y: 32, W: 8, H: 8},
	{X: 8, Y: 32, W: 8, H: 8},
	{X: 16, Y: 32, W: 8, H: 8},
	{X: 24, Y: 32, W: 8, H: 8},
	{X: 32, Y: 32, W: 8, H: 8},
	{X: 40, Y: 32, W: 8, H: 8},
	{X: 48, Y: 32, W: 8, H: 8},
	{X: 56, Y: 32, W: 8, H: 8},
	{X: 64, Y: 32, W: 8, H: 8},
	{X: 72, Y: 32, W: 8, H: 8},
	{X: 80, Y: 32, W: 8, H: 8},
	{X: 88, Y: 32, W: 8, H: 8},
	{X: 96, Y: 32, W: 8, H: 8},
	{X: 104, Y: 32, W: 8, H: 8},
	{X: 112, Y: 32, W: 8, H: 8},
	{X: 120, Y: 32, W: 8, H: 8},
	{X: 0, Y: 40, W: 8, H: 8},
	{X: 8, Y: 40, W: 8, H: 8},
	{X: 16, Y: 40, W: 8, H: 8},
	{X: 24, Y: 40, W: 8, H: 8},
	{X: 32, Y: 40, W: 8, H: 8},
	{X: 40, Y: 40, W: 8, H: 8},
	{X: 48, Y: 40, W: 8, H: 8},
	{X: 56, Y: 40, W: 8, H: 8},
	{X: 64, Y: 40, W: 8, H: 8},
	{X: 72, Y: 40, W: 8, H: 8},
	{X: 80, Y: 40, W: 8, H: 8},
	{X: 88, Y: 40, W: 8, H: 8},
	{X: 96, Y: 40, W: 8, H: 8},
	{X: 104, Y: 40, W: 8, H: 8},
	{X: 112, Y: 40, W: 8, H: 8},
	{X: 120, Y: 40, W: 8, H: 8},
	{X: 0, Y: 48, W: 8, H: 8},
	{X: 8, Y: 48, W: 8, H: 8},
	{X: 16, Y: 48, W: 8, H: 8},
	{X: 24, Y: 48, W: 8, H: 8},
	{X: 32, Y: 48, W: 8, H: 8},
	{X: 40, Y: 48, W: 8, H: 8},
	{X: 48, Y: 48, W: 8, H: 8},
	{X: 56, Y: 48, W: 8, H: 8},
	{X: 64, Y: 48, W: 8, H: 8},
	{X: 72, Y: 48, W: 8, H: 8},
	{X: 80, Y: 48, W: 8, H: 8},
	{X: 88, Y: 48, W: 8, H: 8},
	{X: 96, Y: 48, W: 8, H: 8},
	{X: 104, Y: 48, W: 8, H: 8},
	{X: 112, Y: 48, W: 8, H: 8},
	{X: 120, Y: 48, W: 8, H: 8},
	{X: 0, Y: 56, W: 8, H: 8},
	{X: 8, Y: 56, W: 8, H: 8},
	{X: 16, Y: 56, W: 8, H: 8},
	{X: 24, Y: 56, W: 8, H: 8},
	{X: 32, Y: 56, W: 8, H: 8},
	{X: 40, Y: 56, W: 8, H: 8},
	{X: 48, Y: 56, W: 8, H: 8},
	{X: 56, Y: 56, W: 8, H: 8},
	{X: 64, Y: 56, W: 8, H: 8},
	{X: 72, Y: 56, W: 8, H: 8},
	{X: 80, Y: 56, W: 8, H: 8},
	{X: 88, Y: 56, W: 8, H: 8},
	{X: 96, Y: 56, W: 8, H: 8},
	{X: 104, Y: 56, W: 8, H: 8},
	{X: 112, Y: 56, W: 8, H: 8},
	{X: 120, Y: 56, W: 8, H: 8},
	{X: 0, Y: 64, W: 8, H: 8},
	{X: 8, Y: 64, W: 8, H: 8},
	{X: 16, Y: 64, W: 8, H: 8},
	{X: 24, Y: 64, W: 8, H: 8},
	{X: 32, Y: 64, W: 8, H: 8},
	{X: 40, Y: 64, W: 8, H: 8},
	{X: 48, Y: 64, W: 8, H: 8},
	{X: 56, Y: 64, W: 8, H: 8},
	{X: 64, Y: 64, W: 8, H: 8},
	{X: 72, Y: 64, W: 8, H: 8},
	{X: 80, Y: 64, W: 8, H: 8},
	{X: 88, Y: 64, W: 8, H: 8},
	{X: 96, Y: 64, W: 8, H: 8},
	{X: 104, Y: 64, W: 8, H: 8},
	{X: 112, Y: 64, W: 8, H: 8},
	{X: 120, Y: 64, W: 8, H: 8},
	{X: 0, Y: 72, W: 8, H: 8},
	{X: 8, Y: 72, W: 8, H: 8},
	{X: 16, Y: 72, W: 8, H: 8},
	{X: 24, Y: 72, W: 8, H: 8},
	{X: 32, Y: 72, W: 8, H: 8},
	{X: 40, Y: 72, W: 8, H: 8},
	{X: 48, Y: 72, W: 8, H: 8},
	{X: 56, Y: 72, W: 8, H: 8},
	{X: 64, Y: 72, W: 8, H: 8},
	{X: 72, Y: 72, W: 8, H: 8},
	{X: 80, Y: 72, W: 8, H: 8},
	{X: 88, Y: 72, W: 8, H: 8},
	{X: 96, Y: 72, W: 8, H: 8},
	{X: 104, Y: 72, W: 8, H: 8},
	{X: 112, Y: 72, W: 8, H: 8},
	{X: 120, Y: 72, W: 8, H: 8},
	{X: 0, Y: 80, W: 8, H: 8},
	{X: 8, Y: 80, W: 8, H: 8},
	{X: 16, Y: 80, W: 8, H: 8},
	{X: 24, Y: 80, W: 8, H: 8},
	{X: 32, Y: 80, W: 8, H: 8},
	{X: 40, Y: 80, W: 8, H: 8},
	{X: 48, Y: 80, W: 8, H: 8},
	{X: 56, Y: 80, W: 8, H: 8},
	{X: 64, Y: 80, W: 8, H: 8},
	{X: 72, Y: 80, W: 8, H: 8},
	{X: 80, Y: 80, W: 8, H: 8},
	{X: 88, Y: 80, W: 8, H: 8},
	{X: 96, Y: 80, W: 8, H: 8},
	{X: 104, Y: 80, W: 8, H: 8},
	{X: 112, Y: 80, W: 8, H: 8},
	{X: 120, Y: 80, W: 8, H: 8},
	{X: 0, Y: 88, W: 8, H: 8},
	{X: 8, Y: 88, W: 8, H: 8},
	{X: 16, Y: 88, W: 8, H: 8},
	{X: 24, Y: 88, W: 8, H: 8},
	{X: 32, Y: 88, W: 8, H: 8},
	{X: 40, Y: 88, W: 8, H: 8},
	{X: 48, Y: 88, W: 8, H: 8},
	{X: 56, Y: 88, W: 8, H: 8},
	{X: 64, Y: 88, W: 8, H: 8},
	{X: 72, Y: 88, W: 8, H: 8},
	{X: 80, Y: 88, W: 8, H: 8},
	{X: 88, Y: 88, W: 8, H: 8},
	{X: 96, Y: 88, W: 8, H: 8},
	{X: 104, Y: 88, W: 8, H: 8},
	{X: 112, Y: 88, W: 8, H: 8},
	{X: 120, Y: 88, W: 8, H: 8},
	{X: 0, Y: 96, W: 8, H: 8},
	{X: 8, Y: 96, W: 8, H: 8},
	{X: 16, Y: 96, W: 8, H: 8},
	{X: 24, Y: 96, W: 8, H: 8},
	{X: 32, Y: 96, W: 8, H: 8},
	{X: 40, Y: 96, W: 8, H: 8},
	{X: 48, Y: 96, W: 8, H: 8},
	{X: 56, Y: 96, W: 8, H: 8},
	{X: 64, Y: 96, W: 8, H: 8},
	{X: 72, Y: 96, W: 8, H: 8},
	{X: 80, Y: 96, W: 8, H: 8},
	{X: 88, Y: 96, W: 8, H: 8},
	{X: 96, Y: 96, W: 8, H: 8},
	{X: 104, Y: 96, W: 8, H: 8},
	{X: 112, Y: 96, W: 8, H: 8},
	{X: 120, Y: 96, W: 8, H: 8},
	{X: 0, Y: 104, W: 8, H: 8},
	{X: 8, Y: 104, W: 8, H: 8},
	{X: 16, Y: 104, W: 8, H: 8},
	{X: 24, Y: 104, W: 8, H: 8},
	{X: 32, Y: 104, W: 8, H: 8},
	{X: 40, Y: 104, W: 8, H: 8},
	{X: 48, Y: 104, W: 8, H: 8},
	{X: 56, Y: 104, W: 8, H: 8},
	{X: 64, Y: 104, W: 8, H: 8},
	{X: 72, Y: 104, W: 8, H: 8},
	{X: 80, Y: 104, W: 8, H: 8},
	{X: 88, Y: 104, W: 8, H: 8},
	{X: 96, Y: 104, W: 8, H: 8},
	{X: 104, Y: 104, W: 8, H: 8},
	{X: 112, Y: 104, W: 8, H: 8},
	{X: 120, Y: 104, W: 8, H: 8},
	{X: 0, Y: 112, W: 8, H: 8},
	{X: 8, Y: 112, W: 8, H: 8},
	{X: 16, Y: 112, W: 8, H: 8},
	{X: 24, Y: 112, W: 8, H: 8},
	{X: 32, Y: 112, W: 8, H: 8},
	{X: 40, Y: 112, W: 8, H: 8},
	{X: 48, Y: 112, W: 8, H: 8},
	{X: 56, Y: 112, W: 8, H: 8},
	{X: 64, Y: 112, W: 8, H: 8},
	{X: 72, Y: 112, W: 8, H: 8},
	{X: 80, Y: 112, W: 8, H: 8},
	{X: 88, Y: 112, W: 8, H: 8},
	{X: 96, Y: 112, W: 8, H: 8},
	{X: 104, Y: 112, W: 8, H: 8},
	{X: 112, Y: 112, W: 8, H: 8},
	{X: 120, Y: 112, W: 8, H: 8},
	{X: 0, Y: 120, W: 8, H: 8},
	{X: 8, Y: 120, W: 8, H: 8},
	{X: 16, Y: 120, W: 8, H: 8},
	{X: 24, Y: 120, W: 8, H: 8},
	{X: 32, Y: 120, W: 8, H: 8},
	{X: 40, Y: 120, W: 8, H: 8},
	{X: 48, Y: 120, W: 8, H: 8},
	{X: 56, Y: 120, W: 8, H: 8},
	{X: 64, Y: 120, W: 8, H: 8},
	{X: 72, Y: 120, W: 8, H: 8},
	{X: 80, Y: 120, W: 8, H: 8},
	{X: 88, Y: 120, W: 8, H: 8},
	{X: 96, Y: 120, W: 8, H: 8},
	{X: 104, Y: 120, W: 8, H: 8},
	{X: 112, Y: 120, W: 8, H: 8},
	{X: 120, Y: 120, W: 8, H: 8},
}

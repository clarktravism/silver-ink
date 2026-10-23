package render

import (
	"bytes"
	"errors"
	"image"
	_ "image/png"
	"os"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

type Palette [16]Color

type Color struct {
	R uint8
	G uint8
	B uint8
}

func init() {
	palette = VGAPalette
}

var palette Palette

func LoadPalette(filename string) error {
	p, err := readPalettePNG(filename)
	if err != nil {
		return err
	}
	palette = p
	return nil
}

func ClearBackground(r *sdl.Renderer) {
	col := palette[0]
	sdl.SetRenderDrawColor(r, col.R, col.G, col.B, sdl.AlphaOpaque)
	sdl.RenderClear(r)
}


func SetColor(value Color, color int) {
	palette[color] = value
}

func GetColor(color int, index int) Color {
	return palette[color]
}

func readPalettePNG(filename string) (Palette, error) {

	// read and validate chr image file
	b, err := os.ReadFile(filename)
	if err != nil {
		return Palette{}, err
	}
	img, _, err := image.Decode(bytes.NewBuffer(b))
	if err != nil {
		return Palette{}, err
	}
	size := img.Bounds().Size()
	if size.X * size.Y != 16 {
		return Palette{}, errors.New("incorrect pal image size")
	}

	// create the palette
	p := Palette{}
	i := 0
	for iy := range(size.Y) {
		for ix := range(size.X) {
			r, g, b, _ := img.At(ix, iy).RGBA()
			p[i] = Color{R: uint8(r), G: uint8(g), B: uint8(b)}
			i++
		}
	}

	return p, nil
}

var VGAPalette = Palette{
	{0, 0, 0},
	{0, 0, 170},
	{0, 170, 0},
	{0, 170, 170},
	{170, 0, 0},
	{170, 0, 170},
	{170, 85, 0},
	{170, 170, 170},
	{85, 85, 85},
	{85, 85, 255},
	{85, 255, 85},
	{85, 255, 255},
	{255, 85, 85},
	{255, 85, 255},
	{255, 255, 85},
	{255, 255, 255},
}

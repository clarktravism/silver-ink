package render

import (
	"silver-ink/data"
	"github.com/jupiterrider/purego-sdl3/sdl"
)

const buffer_batch_size = 32

type Buffer struct {
	r    *sdl.Renderer
	offsetX float32
	offsetY float32
	rects [buffer_batch_size]sdl.FRect
	width int
	height int
	size int
	Cells []data.Cell
}

func NewBuffer(width int, height int, renderer *sdl.Renderer) *Buffer {
	return &Buffer{
		width: width,
		height: height,
		size: width * height,
		Cells: make([]data.Cell, width*height),
		r: renderer,
	}
}

func (buf *Buffer) SetRenderer(renderer *sdl.Renderer) {
	buf.r = renderer
}

func (buf *Buffer) SetOffset(x, y float32) {
	buf.offsetX = x
	buf.offsetY = y
}

func (buf *Buffer) Render() {
	buf.renderBG()
	buf.renderFG()
}

func (buf *Buffer) renderBG() {
	rects := buf.rects[:0]
	for c := data.Color(1); c < 16; c++ {
		sdl.SetRenderDrawColor(buf.r, palette[c].R, palette[c].G, palette[c].B, sdl.AlphaOpaque)
		//iterate over cells within the view
		y := float32(0)
		x := float32(0)
		i := 0
		for iy := 0; iy < buf.height; iy++ {
			for ix := 0; ix < buf.width; ix++ {
				if buf.Cells[i].Color.BG() == c {
					rects = append(rects, sdl.FRect{
						X: x + buf.offsetX, Y: y + buf.offsetY, W: 8, H: 8,
					})
					//clear/flush rects if full
					if len(rects) == buffer_batch_size {
						sdl.RenderFillRects(buf.r, rects[:])
						rects = buf.rects[:0]
					}
				}
				x += 8
				i++
			}
			x = 0
			y += 8
		}
		//finish flush before moving to next color
		sdl.RenderFillRects(buf.r, rects)
		rects = buf.rects[:0]
	}
}

func (buf *Buffer) renderFG() {
	for c := range data.Color(16) {
		sdl.SetTextureColorMod(charset, palette[c].R, palette[c].G, palette[c].B)
		//iterate over cells within the view
		rect := sdl.FRect{X: buf.offsetX, Y: buf.offsetY, W:8, H:8}
		i := 0
		for iy := 0; iy < buf.height; iy++ {
			for ix := 0; ix < buf.width; ix++ {
				if buf.Cells[i].Char != 0 && buf.Cells[i].Color.FG() == c {
					sdl.RenderTexture(buf.r, charset, &charLookUp[buf.Cells[i].Char], &rect)
				}
				rect.X += 8
				i++
			}
			rect.X = buf.offsetX
			rect.Y += 8
		}
	}
}

func (buf *Buffer) Stride() int {
	return buf.width
}

func (buf *Buffer) GetCells() []data.Cell {
	return buf.Cells
}

func (buf *Buffer) SetCell(index int, cell data.Cell) {
	if index < 0 || index >= len(buf.Cells) { return }
	buf.Cells[index] = cell
}

func (buf *Buffer) GetCell(index int) data.Cell {
	if index < 0 || index >= len(buf.Cells) { return data.Cell{} }
	return buf.Cells[index]
}

func (buf *Buffer) Index(x int, y int) int {
	return int(x + (y * buf.width))
}

func (buf *Buffer) VectorIndex(vector data.Vector) int {
	return buf.Index(vector.X, vector.Y)
}

func (buf *Buffer) WriteText(index int, color data.Color, text string) int {
	for _, b := range []byte(text) {
		buf.SetCell(index, data.Cell{Char: b, Color: color})
		index++
	}
	return index
}

func (buf *Buffer) FillRectangle(index int, size data.Vector, cell data.Cell) {
	for iy := 0; iy < size.Y; iy++ {
		i := index
		for ix := 0; ix < size.X; ix++ {
			buf.SetCell(i, cell)
			i++
		}
		index += buf.Stride()
	}
}
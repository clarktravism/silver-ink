package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

var cellview_background = data.Cell{Char: 0, Color: 0}

type CellView struct {
	Position data.Vector
	Cell     data.Cell

	Click engine.Signal[engine.MouseState]

	ShowHover  bool
	HoverColor data.Color

	rect  data.Rect
	hover bool
}

func (cv *CellView) Init() {
	cv.rect = data.Rect{
		Position: cv.Position,
		Size:     data.Vector{X: 3, Y: 3},
	}
}

func (cv *CellView) Rect() data.Rect {
	return cv.rect
}

func (cv *CellView) ProcessMouse(ms engine.MouseState) {
	inside := cv.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(cv.rect.Position)
	if cv.ShowHover && inside != cv.hover {
		cv.hover = inside
		engine.Dirty()
	}
	if inside && ms.LeftButtonClick {
		cv.Click.Emit(ms)
	}
}

func (cv *CellView) SetCell(value data.Cell) {
	if cv.Cell != value {
		cv.Cell = value
		engine.Dirty()
	}
}

func (cv *CellView) Render(buffer *render.Buffer, offset data.Vector) {
	pos := cv.Position.Add(offset)

	// hover
	background := cellview_background
	cell := cv.Cell
	if cv.hover {
		background.Color = cv.HoverColor
		cell.Color = cv.HoverColor
	}

	biy := buffer.Index(pos.X, pos.Y)
	for iy := 0; iy < 3; iy++ {
		bix := biy
		for ix := 0; ix < 3; ix++ {
			buffer.SetCell(bix, background)
			bix++
		}
		biy += buffer.Stride()
	}

	buffer.SetCell(buffer.Index(pos.X+1, pos.Y+1), cell)
}

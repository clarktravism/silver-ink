package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

var colorselect_ind_color = data.NewColor(data.ColorWHITE, data.ColorBLACK)

type ColorSelect struct {
	Position data.Vector
	Color    data.Color

	// signals
	ColorChanged engine.Signal[ColorChangedArgs]

	rect  data.Rect
	hover bool
}

// signal args
type ColorChangedArgs struct {
	Color data.Color
}

func (cs *ColorSelect) SetColor(value data.Color) {
	if cs.Color != value {
		cs.Color = value
		engine.Dirty()
		cs.ColorChanged.Emit(ColorChangedArgs{Color: value})
	}
}

// control overrides

func (cs *ColorSelect) Init() {
	cs.rect = data.Rect{
		Position: cs.Position,
		Size:     data.Vector{X: 16, Y: 2},
	}
}

func (cs *ColorSelect) Rect() data.Rect {
	return cs.rect
}

func (cs *ColorSelect) ProcessMouse(ms engine.MouseState) {
	inside := cs.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(cs.rect.Position)
	if inside != cs.hover {
		cs.hover = inside
	}
	if inside && (ms.LeftButtonClick != ms.RightButtonClick) {
		color := cs.Color
		newColor := data.Color(ms.Position.X)
		if ms.LeftButtonClick {
			cs.SetColor(data.NewColor(newColor, color.BG()))
		} else if ms.RightButtonClick {
			cs.SetColor(data.NewColor(color.FG(), newColor))
		}
	}
}

func (cs *ColorSelect) Render(buffer *render.Buffer, offset data.Vector) {
	pos := cs.Position.Add(offset)

	// colors
	bi := buffer.Index(pos.X, pos.Y)
	for i := data.Color(0); i < 16; i++ {
		buffer.SetCell(bi, data.Cell{Char: 0, Color: data.NewColor(0, i)})
		bi++
	}

	// indicators
	bi = buffer.Index(pos.X, pos.Y+1)
	for i := data.Color(0); i < 16; i++ {
		// set char for ind
		char := byte(0)
		if cs.Color.FG() == i {
			if cs.Color.BG() == i {
				char = 'X'
			} else {
				char = 'F'
			}
		} else if cs.Color.BG() == i {
			char = 'B'
		}
		buffer.SetCell(bi, data.Cell{Char: char, Color: colorselect_ind_color})
		bi++
	}
}

package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

var charselect_color = data.NewColor(data.ColorGRAY, data.ColorBLACK)
var charselect_not_hover_color = data.NewColor(data.ColorDKGRAY, data.ColorBLACK)
var charselect_hover_color = data.NewColor(data.ColorYELLOW, data.ColorBLACK)

type CharSelect struct {
	Position  data.Vector
	Char      byte
	hoverChar byte

	//signals
	CharChanged engine.Signal[CharChangedArgs]

	rect  data.Rect
	hover bool
}

// signal args

type CharChangedArgs struct {
	Char byte
}

func (cs *CharSelect) SetChar(value byte) {
	if cs.Char != value {
		cs.Char = value
		engine.Dirty()
		cs.CharChanged.Emit(CharChangedArgs{Char: value})
	}
}

// control overrides

func (cs *CharSelect) Init() {
	cs.rect = data.Rect{
		Position: cs.Position,
		Size:     data.Vector{X: 16, Y: 16},
	}
}

func (cs *CharSelect) Rect() data.Rect {
	return cs.rect
}

func (cs *CharSelect) ProcessMouse(ms engine.MouseState) {
	inside := cs.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(cs.rect.Position)
	if inside != cs.hover {
		cs.hover = inside
		engine.Dirty()
	}
	if inside {
		char := byte(ms.Position.X + (ms.Position.Y * 16))
		cs.hoverChar = char
		if ms.LeftButtonClick {
			cs.SetChar(char)
		}
		engine.Dirty()
	}
}

func (cs *CharSelect) Render(buffer *render.Buffer, offset data.Vector) {
	pos := cs.Position.Add(offset)
	char := byte(0)
	biy := buffer.Index(pos.X, pos.Y)
	for iy := 0; iy < 16; iy++ {
		bix := biy
		for ix := 0; ix < 16; ix++ {

			// color based on hover
			color := charselect_color
			if cs.hover {
				if cs.hoverChar == char {
					color = charselect_hover_color
				} else {
					color = charselect_not_hover_color
				}
			}

			buffer.SetCell(bix, data.Cell{Char: char, Color: color})
			bix++
			char++
		}
		biy += buffer.Stride()
	}
}

// signal implementions
func (cs *CharSelect) onMouseMoved(ms engine.MouseState) {
	pos := ms.Position.Sub(cs.Position)
	char := byte(pos.X + (pos.Y * 16))
	if cs.hoverChar != char {
		cs.hoverChar = char
		engine.Dirty()
	}
}

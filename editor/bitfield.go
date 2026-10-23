package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

// button with side label
type BitField struct {
	Position data.Vector
	Value    data.BitField

	ValueChanged engine.Signal[BitFieldChangedArgs]

	rect  data.Rect
	hover bool
}

type BitFieldChangedArgs struct {
	Value data.BitField
}

func (bf *BitField) Init() {
	bf.rect = data.Rect{
		Position: bf.Position,
		Size:     data.Vector{X: 8, Y: 1},
	}
}

func (bf *BitField) Rect() data.Rect {
	return bf.rect
}

func (bf *BitField) ProcessMouse(ms engine.MouseState) {
	inside := bf.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(bf.rect.Position)
	if inside != bf.hover {
		bf.hover = inside
		engine.Dirty()
	}
	if inside && ms.LeftButtonClick {
		p := 7 - ms.Position.X
		b := data.BitField(1 << p)
		bf.Value = bf.Value ^ b
		bf.ValueChanged.Emit(BitFieldChangedArgs{
			Value: bf.Value,
		})
		engine.Dirty()
	}
}

func (bf *BitField) Render(buffer *render.Buffer, offset data.Vector) {
	offset = offset.Add(bf.rect.Position)
	bi := buffer.Index(offset.X, offset.Y)
	for i := 7; i >= 0; i-- {
		cell := data.Cell{Char: byte('1'+i), Color: data.ColorDKGRAY}
		if bf.Value&(1<<i) != 0 {
			cell.Color = data.ColorWHITE
		}
		buffer.SetCell(bi, cell)
		bi++
	}
}

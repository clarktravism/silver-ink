package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

type BrushView struct {
	Position data.Vector
	Brush    Brush

	cellView  CellView
	typeLabel engine.Label
	nameLabel engine.Label

	rect data.Rect
}

func (bv *BrushView) Init() {
	bv.rect = data.Rect{
		Position: bv.Position,
		Size:     data.Vector{X: 16, Y: 3},
	}

	bv.typeLabel = engine.Label{
		Position: data.Vector{X: 0, Y: 0},
		Color:    data.NewColor(data.ColorWHITE, data.ColorDKGRAY),
		Length:   12,
	}
	bv.nameLabel = engine.Label{
		Position: data.Vector{X: 0, Y: 1},
		Color:    data.NewColor(data.ColorWHITE, data.ColorDKGRAY),
		Length:   12,
	}

	bv.cellView = CellView{
		Position:  data.Vector{X: 13, Y: 0},
		ShowHover: false,
	}
	bv.cellView.Init()
}

func (bv *BrushView) Rect() data.Rect {
	return bv.rect
}

func (bv *BrushView) ProcessMouse(engine.MouseState) {}

func (bv *BrushView) SetBrush(value Brush) {
	if bv.Brush != value {
		bv.Brush = value
		bv.updateBrushValues()
		engine.Dirty()
	}
}

func (bv *BrushView) Refresh() {
	bv.updateBrushValues()
}

func (bv *BrushView) updateBrushValues() {
	bv.typeLabel.SetText(bv.Brush.Type())
	bv.nameLabel.SetText(bv.Brush.Name().String())
	bv.cellView.SetCell(bv.Brush.Cell)
}

func (bv *BrushView) Render(buffer *render.Buffer, offset data.Vector) {
	pos := bv.Position.Add(offset)
	bv.typeLabel.Render(buffer, pos)
	bv.nameLabel.Render(buffer, pos)
	bv.cellView.Render(buffer, pos)
}

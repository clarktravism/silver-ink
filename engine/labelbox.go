package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type LabelBox struct {
	Rectangle data.Rect
	Text      string
	Color     data.Color
	Style     data.BoxStyle

	box Box

	Dirty Signal[struct{}]
}

func (lb *LabelBox) Init() {
	lb.box.Rectangle = data.Rect{Position: data.VectorZero, Size: lb.Rectangle.Size}
	lb.box.Style = lb.Style
}

func (lb *LabelBox) ProcessMouse(MouseState) {}

func (lb *LabelBox) Rect() data.Rect { return data.RectZero }

func (lb *LabelBox) SetText(text string) {
	if lb.Text != text {
		lb.Text = text
		lb.Dirty.Emit(struct{}{})
	}
}

func (lb *LabelBox) SetColor(color data.Color) {
	if lb.Color != color {
		lb.Color = color
		lb.Dirty.Emit(struct{}{})
	}
}

func (lb *LabelBox) Render(buffer *render.Buffer, offset data.Vector) {
	offset = lb.Rectangle.Position.Add(offset)
	lb.box.Render(buffer, offset)

	bi := buffer.Index(offset.X+1, offset.Y)
	buffer.WriteText(bi, lb.Color, lb.Text)
}

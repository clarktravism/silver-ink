package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type Label struct {
	Position data.Vector
	Text     string
	Color    data.Color
	Length   int
}

func (l *Label) Init() {}

func (l *Label) ProcessMouse(MouseState) {}

func (l *Label) Rect() data.Rect { return data.RectZero }

func (l *Label) SetText(text string) {
	if len(text) > l.Length {
		text = text[:l.Length]
	}
	if l.Text != text {
		l.Text = text
	}
}

func (l *Label) SetColor(color data.Color) {
	if l.Color != color {
		l.Color = color
	}
}

func (l *Label) Render(buffer *render.Buffer, offset data.Vector) {
	offset = l.Position.Add(offset)
	bi := buffer.Index(offset.X, offset.Y)
	buffer.WriteText(bi, l.Color, l.Text)
}

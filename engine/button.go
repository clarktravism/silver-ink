package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type ButtonStyle struct {
	Background      data.Color
	HoverBackground data.Color
}

// button with side label
type Button struct {
	Position data.Vector
	Style    ButtonStyle
	Text     string

	Click Signal[MouseState]

	rect  data.Rect
	hover bool
}

func (b *Button) Init() {
	b.rect = data.Rect{
		Position: b.Position,
		Size:     data.Vector{X: len(b.Text), Y: 1},
	}
}

func (b *Button) Rect() data.Rect {
	return b.rect
}

func (b *Button) ProcessMouse(ms MouseState) {
	inside := b.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(b.rect.Position)
	if inside != b.hover {
		b.hover = inside
		Dirty()
	}
	if inside && ms.LeftButtonClick {
		b.Click.Emit(ms)
	}
}

func (b *Button) Render(buffer *render.Buffer, offset data.Vector) {
	offset = offset.Add(b.rect.Position)
	bi := buffer.Index(offset.X, offset.Y)
	if b.hover {
		bi = buffer.WriteText(bi, b.Style.HoverBackground, b.Text)
	} else {
		bi = buffer.WriteText(bi, b.Style.Background, b.Text)
	}
}

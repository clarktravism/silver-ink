package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type LabelButtonStyle struct {
	LeftBackground       data.Color
	RightBackground      data.Color
	LeftHoverBackground  data.Color
	RightHoverBackground data.Color
}

// button with side label
type LabelButton struct {
	Position  data.Vector
	Style     LabelButtonStyle
	LeftText  string
	RightText string

	Click Signal[MouseState]

	rect  data.Rect
	hover bool
}

func (lb *LabelButton) Init() {
	lb.rect = data.Rect{
		Position: lb.Position,
		Size:     data.Vector{X: len(lb.LeftText) + len(lb.RightText), Y: 1},
	}
}

func (lb *LabelButton) Rect() data.Rect {
	return lb.rect
}

func (lb *LabelButton) ProcessMouse(ms MouseState) {
	inside := lb.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(lb.rect.Position)
	if inside != lb.hover {
		lb.hover = inside
		Dirty()
	}
	if inside && ms.LeftButtonClick {
		lb.Click.Emit(ms)
	}
}

func (lb *LabelButton) Render(buffer *render.Buffer, offset data.Vector) {
	offset = offset.Add(lb.rect.Position)
	bi := buffer.Index(offset.X, offset.Y)
	if lb.LeftText != "" {
		if lb.hover {
			bi = buffer.WriteText(bi, lb.Style.LeftHoverBackground, lb.LeftText)
		} else {
			bi = buffer.WriteText(bi, lb.Style.LeftBackground, lb.LeftText)
		}
	}
	if lb.RightText != "" {
		if lb.hover {
			bi = buffer.WriteText(bi, lb.Style.RightHoverBackground, lb.RightText)
		} else {
			bi = buffer.WriteText(bi, lb.Style.RightBackground, lb.RightText)
		}
	}
}

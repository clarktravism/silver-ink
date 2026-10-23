package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type CheckBoxStyle struct {
	CheckedOn  data.Cell
	CheckedOff data.Cell
}

type CheckBox struct {
	Position data.Vector
	Style    CheckBoxStyle

	Checked bool

	CheckedChanged Signal[CheckChangedArgs]

	rect  data.Rect
	hover bool
}

type CheckChangedArgs struct {
	Checked bool
}

func (cb *CheckBox) Init() {
	cb.rect = data.Rect{
		Position: cb.Position,
		Size:     data.Vector{X: 1, Y: 1},
	}
}

func (cb *CheckBox) Rect() data.Rect {
	return cb.rect
}

func (cb *CheckBox) ProcessMouse(ms MouseState) {
	inside := cb.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(cb.rect.Position)
	if inside != cb.hover {
		cb.hover = inside
		Dirty()
	}
	if inside && ms.LeftButtonClick {
		cb.Checked = !cb.Checked
		cb.CheckedChanged.Emit(CheckChangedArgs{
			Checked: cb.Checked,
		})
		Dirty()
	}
}

func (cb *CheckBox) Render(buffer *render.Buffer, offset data.Vector) {
	offset = offset.Add(cb.Position)
	cell := cb.Style.CheckedOff
	if cb.Checked {
		cell = cb.Style.CheckedOn
	}

	bi := buffer.Index(offset.X, offset.Y)
	buffer.SetCell(bi, cell)
}

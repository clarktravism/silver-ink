package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type TextField struct {
	Position   data.Vector
	Background data.Cell
	Cursor     data.Cell
	TextInput
	rect data.Rect
}

func (tf *TextField) Init() {
	tf.rect = data.Rect{
		Position: tf.Position,
		Size:     data.Vector{X: tf.MaxLength, Y: 1},
	}
}

func (tf *TextField) Rect() data.Rect {
	return tf.rect
}

func (tf *TextField) ProcessMouse(ms MouseState) {
	inside := tf.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(tf.rect.Position)
	if inside && ms.LeftButtonClick {
		tf.TextInput.TakeFocus()
		Dirty()
	}
}

func (tf *TextField) onTextInputAccept() {
	Dirty()
}

func (tf *TextField) onTextInputCancel() {
	Dirty()
}

func (tf *TextField) Render(buffer *render.Buffer, offset data.Vector) {
	offset = offset.Add(tf.Position)
	index := buffer.Index(offset.X, offset.Y)
	buffer.FillRectangle(index, tf.rect.Size, tf.Background)
	buffer.WriteText(index, tf.Background.Color, tf.Text)
	if tf.TextInput.Focused && len(tf.Text) < tf.MaxLength {
		buffer.SetCell(index+len(tf.Text), tf.Cursor)
	}
}

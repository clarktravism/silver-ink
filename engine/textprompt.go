package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type TextPrompt struct {
	Length      int
	Title       string
	Text        string
	AllowCancel bool
	Style       data.BoxStyle

	position  data.Vector
	frame     Box
	textField TextField

	Callback func(TextPromptCallback)
}

type TextPromptCallback struct {
	Accepted    bool
	AllowCancel bool
	Text        string
}

func (tp *TextPrompt) Push() {
	tp.init()
	PushStack(tp)
}

func (tp *TextPrompt) init() {
	size := data.Vector{X: tp.Length + 2, Y: 5}
	tp.position = CenterPos(size)

	// frame
	tp.frame.Rectangle = data.Rect{
		Position: data.VectorZero,
		Size:     size,
	}
	tp.frame.Title = tp.Title
	tp.frame.Style = tp.Style

	// textfield
	tp.textField = TextField{
		Position:   data.Vector{X: 1, Y: 3},
		Background: data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorBLACK)},
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: TextInput{
			MaxLength: tp.Length,
			Text:      tp.Text,
		},
	}
	tp.textField.Init()
	tp.textField.OnAccept = tp.onAccept
	tp.textField.OnCancel = tp.onCancel
	tp.textField.TakeFocus()
}

func (tp *TextPrompt) onAccept() {
	if tp.textField.Text == "" && !tp.AllowCancel {
		tp.textField.TakeFocus()
		return
	}
	PopStack()
	tp.Callback(TextPromptCallback{
		Accepted:    true,
		Text:        tp.textField.Text,
		AllowCancel: tp.AllowCancel,
	})
}

func (tp *TextPrompt) onCancel() {
	if !tp.AllowCancel {
		tp.textField.TakeFocus()
		return
	}
	PopStack()
	tp.Callback(TextPromptCallback{
		Accepted:    false,
		AllowCancel: tp.AllowCancel,
	})
}

// processor
func (tp *TextPrompt) Render(buffer *render.Buffer, top bool) {
	offset := tp.position
	tp.frame.Render(buffer, offset)
	tp.textField.Render(buffer, offset)
}

func (tp *TextPrompt) Active() {
	FPS(20)
	tp.textField.TakeFocus()
}

func (tp *TextPrompt) Process() {}

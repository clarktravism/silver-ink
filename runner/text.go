package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

var colorCycle = []data.Color{11, 10, 9, 15, 14, 13, 12}

var Text text

var TextBoxIsOpen bool

type text struct {
	lines []engine.ScrollBoxLine
	id int
	textBoxID int
	title string
	msgTicks int
	msgColor int
	msgText string
}

func ResetText() {
	Text.lines = Text.lines[:0]
	Text.id = -1
	Text.textBoxID = TextBoxFile.DialogID
	Text.title = ""
	Text.msgTicks = 0
	Text.msgColor = 0
	Text.msgText = ""
}

func (t *text) AddLine(line engine.ScrollBoxLine) {
	t.lines = append(t.lines, line)
}

func (t *text) ShowText(id int) {
	if len(t.lines) == 1 { // one line
		t.SetMessage(t.lines[0].Text)
		t.lines = t.lines[:0]
	} else if len(t.lines) > 1 { // multi line
		t.id = id
		t.PushTextBox()
		t.lines = t.lines[:0]
	}
}

func (t *text) SetMessage(text string) {
	t.msgText = text
	t.msgTicks = 40
	t.msgColor = 0
}

func (t *text) PushTextBox() {
	if t.textBoxID == 0 {
		t.textBoxID = TextBoxFile.DialogID // default
	}
	scrollbox := CreateTextBox(t.textBoxID, "")
	scrollbox.Title = t.title
	scrollbox.Callback = t.TextBoxCallback
	for _, line := range t.lines {
		scrollbox.AppendLine(line)
	}
	t.textBoxID = TextBoxFile.DialogID 
	TextBoxIsOpen = true
	scrollbox.Push()
}

func (t *text) TextBoxCallback(callback engine.ScrollBoxCallback) {
	TextBoxIsOpen = false
	if callback.Accepted && !callback.Link.Empty() {
		ObjectEvent(t.id, callback.Link, 0)
	}
}

// called every cycle to process message
func (t *text) Process() {
	if t.msgTicks > 0 {
		t.msgTicks--
		t.msgColor++
		if t.msgColor >= len(colorCycle) {
			t.msgColor = 0
		}
	}
}

// draws the message
func (t *text) Render(buffer *render.Buffer) {
	if t.msgTicks > 0 {
		offx := (view.Camera.Size.X >> 1) - (len(t.msgText) >> 1)
		bi := buffer.Index(offx, view.Camera.Position.Y + view.Camera.Size.Y - 1)
		buffer.WriteText(bi, colorCycle[t.msgColor], t.msgText)
	}
}


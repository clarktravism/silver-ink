package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

type SavePrompt struct {
	position   data.Vector
	frame      engine.Box
	textField  engine.TextField
	saveButton engine.Button
	noButton   engine.Button

	Callback func(SavePromptCallback)
}

type SavePromptCallback struct {
	Accepted bool
	Save     bool
}

func (sp *SavePrompt) Push() {
	sp.init()
	engine.PushStack(sp)
}

func (sp *SavePrompt) init() {
	size := data.Vector{X: 20, Y: 5}
	sp.position = engine.CenterPos(size)

	// frame
	sp.frame.Rectangle = data.Rect{
		Position: data.VectorZero,
		Size:     size,
	}
	sp.frame.Title = "Save World"
	sp.frame.Style = text_prompt_style

	sp.saveButton = engine.Button{
		Position: data.Vector{X: 15, Y: 3},
		Style:    green_button_style,
		Text:     " Save ",
	}
	sp.saveButton.Click.Listen(sp.onSave)

	sp.saveButton = engine.Button{
		Position: data.Vector{X: 2, Y: 3},
		Style:    red_button_style,
		Text:     "  No  ",
	}
	sp.noButton.Click.Listen(sp.onNo)
}

func (sp *SavePrompt) onSave(engine.MouseState) {
	engine.PopStack()
	sp.Callback(SavePromptCallback{
		Accepted: true,
		Save:     true,
	})
}

func (sp *SavePrompt) onNo(engine.MouseState) {
	engine.PopStack()
	sp.Callback(SavePromptCallback{
		Accepted: true,
		Save:     false,
	})
}

// processor
func (sp *SavePrompt) Render(buffer *render.Buffer, top bool) {
	offset := sp.position
	sp.frame.Render(buffer, offset)
	sp.saveButton.Render(buffer, offset)
	sp.noButton.Render(buffer, offset)
}

func (sp *SavePrompt) Active() {
	engine.FPS(20)
}

func (sp *SavePrompt) Process() {}

package main

import (
	"log"
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
	"strconv"
)

var editworldinfo_textcolor = data.NewColor(data.ColorWHITE, data.ColorDKGRAY)
var editworldinfo_textfield_background = data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorBLACK)}

type EditWorldInfo struct {
	Position data.Vector

	box1 engine.Box //outer frame
	box2 engine.Box //counters

	// controls
	container        engine.Container
	titleTextField   engine.TextField
	startTextField   engine.TextField
	fpsTextField     engine.TextField
	inputTextField   engine.TextField
	counterAddButton engine.Button
	cancelButton     engine.Button
	acceptButton     engine.Button
	counterControl   CounterControl

	// dialog callback
	Callback func()
}

func (eo *EditWorldInfo) Push() {
	eo.init()
	engine.PushStack(eo)
}

func (eo *EditWorldInfo) init() {
	size := data.Vector{X: 40, Y: 37}

	eo.Position = engine.CenterPos(size)

	// container
	eo.container = engine.Container{
		Rectangle: data.Rect{
			Position: data.VectorZero,
			Size:     size,
		},
	}

	// box1
	eo.box1 = engine.Box{
		Title: "World",
		Rectangle: data.Rect{
			Position: data.VectorZero,
			Size:     size,
		},
		Style: dialog_box_style,
	}

	// box2
	eo.box2 = engine.Box{
		Rectangle: data.Rect{
			Position: data.Vector{X: 2, Y: 12},
			Size:     data.Vector{X: 36, Y: 21},
		},
		Style: inside_box_style,
	}

	// text fields
	eo.titleTextField = engine.TextField{
		Position:   data.Vector{X: 14, Y: 4},
		Background: editworldinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.titleTextField.Text = engine.WorldFile.TitleBoard.String()

	eo.startTextField = engine.TextField{
		Position:   data.Vector{X: 14, Y: 6},
		Background: editworldinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.startTextField.Text = engine.WorldFile.StartBoard.String()

	eo.fpsTextField = engine.TextField{
		Position:   data.Vector{X: 14, Y: 8},
		Background: editworldinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 3,
			TextType:  engine.TextNumbers,
		},
	}
	eo.fpsTextField.Text = strconv.Itoa(engine.WorldFile.FPS)

	eo.inputTextField = engine.TextField{
		Position:   data.Vector{X: 14, Y: 10},
		Background: editworldinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 3,
			TextType:  engine.TextNumbers,
		},
	}
	eo.inputTextField.Text = strconv.Itoa(int(engine.WorldFile.InputCycle))

	// buttons
	eo.counterAddButton = engine.Button{
		Position: data.Vector{X: 13, Y: 13},
		Style:    gray_button_style,
		Text:     " Add ",
	}
	eo.cancelButton = engine.Button{
		Position: data.Vector{X: 1, Y: 35},
		Style:    red_button_style,
		Text:     " X ",
	}
	eo.acceptButton = engine.Button{
		Position: data.Vector{X: 36, Y: 35},
		Style:    green_button_style,
		Text:     " O ",
	}

	// misc
	eo.counterControl = CounterControl{
		Position: data.Vector{X: 4, Y: 15},
		Counters: engine.WorldFile.Counters.Copy(),
	}

	// add to container / init
	eo.container.Init()
	eo.container.AddChild(&eo.titleTextField)
	eo.container.AddChild(&eo.startTextField)
	eo.container.AddChild(&eo.fpsTextField)
	eo.container.AddChild(&eo.inputTextField)
	eo.container.AddChild(&eo.counterAddButton)
	eo.container.AddChild(&eo.counterControl)
	eo.container.AddChild(&eo.cancelButton)
	eo.container.AddChild(&eo.acceptButton)

	// hook up buttons
	eo.counterAddButton.Click.Listen(eo.counterControl.onAddButtonClick)
	eo.acceptButton.Click.Listen(eo.onAcceptButtonClick)
	eo.cancelButton.Click.Listen(eo.onCancelButtonClick)
}

// button events
func (ewi *EditWorldInfo) onAcceptButtonClick(engine.MouseState) {
	ewi.accept()
}

func (ewi *EditWorldInfo) accept() {
	engine.WorldFile.TitleBoard = data.NewName(ewi.titleTextField.Text)
	engine.WorldFile.StartBoard = data.NewName(ewi.startTextField.Text)
	engine.WorldFile.FPS, _ = strconv.Atoi(ewi.fpsTextField.Text)
	inputCycle, _ := strconv.Atoi(ewi.inputTextField.Text)
	engine.WorldFile.InputCycle = uint8(inputCycle)
	engine.WorldFile.Counters = ewi.counterControl.GetCounters()
	if err := engine.WorldFile.Save(); err != nil {
		log.Printf("error while saving world file: %v\r\n", err)
	}

	engine.PopStack()
	if ewi.Callback != nil {
		ewi.Callback()
	}
}

func (ewi *EditWorldInfo) onCancelButtonClick(ms engine.MouseState) {
	ewi.cancel()
}

func (ewi *EditWorldInfo) cancel() {
	engine.PopStack()
	if ewi.Callback != nil {
		ewi.Callback()
	}
}

// stack
func (eo *EditWorldInfo) ProcessMouse(ms engine.MouseState) {
	ms.Position = ms.Position.Sub(eo.Position)
	eo.container.ProcessMouse(ms)
}

func (eo *EditWorldInfo) ProcessAction(action data.Name) {
	//TD; hotkeys
}

func (eo *EditWorldInfo) Render(buffer *render.Buffer, top bool) {
	offset := eo.Position

	eo.box1.Render(buffer, offset)
	eo.box2.Render(buffer, offset)

	// static text
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 7, Y: 4})), editboardinfo_textcolor, "Title:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 34, Y: 4})), editboardinfo_textcolor, ".sib")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 7, Y: 6})), editboardinfo_textcolor, "Start:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 34, Y: 6})), editboardinfo_textcolor, ".sib")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 9, Y: 8})), editboardinfo_textcolor, "FPS:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 2, Y: 10})), editboardinfo_textcolor, "InputCycle:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 3, Y: 13})), editboardinfo_textcolor, "Counters:")

	eo.container.Render(buffer, offset)
}

func (eo *EditWorldInfo) Active() {
	engine.FPS(30)
}

func (ewi *EditWorldInfo) Process() {
	for _, ms := range engine.PullMouseState() {
		ms.Position = ms.Position.Sub(ewi.Position)
		ewi.container.ProcessMouse(ms)
	}

	if !engine.CycleTick(4) {
		return
	}
	if action := engine.GetAction(); action != data.DefaultName {
		switch action {
		case engine.INPUT_CANCEL:
			ewi.cancel()
		case engine.INPUT_ACCEPT:
			ewi.accept()
		}
	}
}

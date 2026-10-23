package main

import (
	"fmt"
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

type EditCell struct {
	Position data.Vector
	Cell     data.Cell

	// boxes / labels
	frame         engine.Box
	charLabelBox  engine.LabelBox
	colorLabelBox engine.LabelBox
	cellBox       engine.Box

	// controls
	container engine.Container

	charSelect   CharSelect
	colorSelect  ColorSelect
	cellView     CellView
	cancelButton engine.Button
	acceptButton engine.Button

	// dialog callback
	Callback func(EditCellCallback)

	rect data.Rect
}

type EditCellCallback struct {
	Accepted bool
	Cell     data.Cell
}

func (ec *EditCell) Push() {
	ec.init()
	engine.PushStack(ec)
}

func (ec *EditCell) init() {
	ec.Position = engine.CenterPos(data.Vector{X: 25, Y: 27})

	ec.rect = data.Rect{
		Position: ec.Position,
		Size:     data.Vector{X: 25, Y: 27},
	}

	ec.container = engine.Container{
		Rectangle: data.Rect{
			Position: data.VectorZero,
			Size:     ec.rect.Size,
		},
	}

	// box1 - outer box
	ec.frame = engine.Box{
		Title: "Edit Cell",
		Rectangle: data.Rect{
			Position: data.VectorZero,
			Size:     data.Vector{X: 25, Y: 27},
		},
		Style: dialog_box_style,
	}

	// inside boxes
	ec.charLabelBox = engine.LabelBox{
		Rectangle: data.NewRect(1, 3, 18, 18),
		Text:      " Char ",
		Color:     sidepanel_boxpanel_color,
		Style:     sidepanel_box_style,
	}
	ec.charLabelBox.Init()

	ec.colorLabelBox = engine.LabelBox{
		Rectangle: data.NewRect(1, 21, 18, 4),
		Text:      " Color ",
		Color:     sidepanel_boxpanel_color,
		Style:     sidepanel_box_style,
	}
	ec.colorLabelBox.Init()

	ec.cellBox = engine.Box{
		Rectangle: data.NewRect(19, 3, 5, 5),
		Style:     inside_box_style,
	}

	// controls
	ec.charSelect = CharSelect{
		Position: data.Vector{X: 2, Y: 4},
	}
	ec.charSelect.CharChanged.Listen(ec.onCharChanged)
	ec.charSelect.SetChar(ec.Cell.Char)
	ec.setCharLabel()

	ec.colorSelect = ColorSelect{
		Position: data.Vector{X: 2, Y: 22},
	}
	ec.colorSelect.ColorChanged.Listen(ec.onColorChanged)
	ec.colorSelect.SetColor(ec.Cell.Color)
	ec.setColorLabel()

	ec.cellView = CellView{
		Position: data.Vector{X: 20, Y: 4},
	}
	ec.cellView.SetCell(ec.Cell)

	ec.cancelButton = engine.Button{
		Position: data.Vector{X: 1, Y: 25},
		Text:     " X ",
		Style:    red_button_style,
	}
	ec.cancelButton.Click.Listen(ec.onCancelButtonClick)

	ec.acceptButton = engine.Button{
		Position: data.Vector{X: 21, Y: 25},
		Text:     " O ",
		Style:    green_button_style,
	}
	ec.acceptButton.Click.Listen(ec.onAcceptButtonClick)

	// add to container / init
	ec.container.Init()
	ec.container.AddChild(&ec.charSelect)
	ec.container.AddChild(&ec.colorSelect)
	ec.container.AddChild(&ec.cancelButton)
	ec.container.AddChild(&ec.acceptButton)
	ec.container.AddChild(&ec.cellView)
}

// buttons
func (ec *EditCell) onAcceptButtonClick(engine.MouseState) {
	ec.accept()
}

func (ec *EditCell) accept() {
	ec.Callback(EditCellCallback{
		Accepted: true,
		Cell: data.Cell{
			Char:  ec.charSelect.Char,
			Color: ec.colorSelect.Color,
		},
	})
	engine.PopStack()
}

func (ec *EditCell) onCancelButtonClick(engine.MouseState) {
	ec.cancel()
}

func (ec *EditCell) cancel() {
	ec.Callback(EditCellCallback{
		Accepted: false,
	})
	engine.PopStack()
}

func (ec *EditCell) onCharChanged(args CharChangedArgs) {
	cell := ec.cellView.Cell
	cell.Char = args.Char
	ec.cellView.SetCell(cell)
	ec.setCharLabel()
}

func (ec *EditCell) onColorChanged(args ColorChangedArgs) {
	cell := ec.cellView.Cell
	cell.Color = args.Color
	ec.cellView.SetCell(cell)
	ec.setColorLabel()
}

// stack
func (ec *EditCell) ProcessMouse(ms engine.MouseState) {
	ms.Position = ms.Position.Sub(ec.Position)
	ec.container.ProcessMouse(ms)
}

func (ec *EditCell) ProcessAction(action data.Name) {
	//TD; hotkeys?
}

func (ec *EditCell) Render(buffer *render.Buffer, top bool) {
	offset := ec.Position

	// render non-interative elements
	ec.frame.Render(buffer, offset)
	ec.charLabelBox.Render(buffer, offset)
	ec.colorLabelBox.Render(buffer, offset)
	ec.cellBox.Render(buffer, offset)
	
	ec.container.Render(buffer, offset)
}

func (ec *EditCell) Active() {
	engine.FPS(30)
}

func (ec *EditCell) Process() {
	for _, ms := range engine.PullMouseState() {
		ms.Position = ms.Position.Sub(ec.Position)
		ec.container.ProcessMouse(ms)
	}

	if !engine.CycleTick(4) {
		return
	}
	if action := engine.GetAction(); action != data.DefaultName {
		switch action {
		case engine.INPUT_CANCEL:
			ec.cancel()
		case engine.INPUT_ACCEPT:
			ec.accept()
		}
	}
}

func (ec *EditCell) setCharLabel() {
	ec.charLabelBox.Text = fmt.Sprintf(" Char %v ", ec.Cell.Char)
}

func (ec *EditCell) setColorLabel() {
	ec.colorLabelBox.Text = fmt.Sprintf(" Color x%02X ", ec.Cell.Color)
}
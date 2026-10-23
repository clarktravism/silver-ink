package main

import (
	"log"
	"os"
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
	"strconv"
)

var editboardinfo_textcolor = data.NewColor(data.ColorWHITE, data.ColorDKGRAY)
var editboardinfo_textfield_background = data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorBLACK)}

type EditBoardInfo struct {
	Position data.Vector

	box1 engine.Box //outer frame
	box2 engine.Box //counters
	box3 engine.Box //exits

	// controls
	container        engine.Container
	nameTextField    engine.TextField
	cameraXTextField engine.TextField
	cameraYTextField engine.TextField
	charsetTextField engine.TextField
	paletteTextField engine.TextField
	viewTextField    engine.TextField
	northTextField   engine.TextField
	southTextField   engine.TextField
	westTextField    engine.TextField
	eastTextField    engine.TextField
	persistCheckBox  engine.CheckBox
	//darkCheckBox     engine.CheckBox
	counterAddButton engine.Button
	cancelButton     engine.Button
	acceptButton     engine.Button
	counterControl   CounterControl

	// dialog callback
	Callback func()
}

func (eo *EditBoardInfo) Push() {
	eo.init()
	engine.PushStack(eo)
}

func (eo *EditBoardInfo) init() {
	size := data.Vector{X: 40, Y: 50}

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
		Title: "Board Info",
		Rectangle: data.Rect{
			Position: data.VectorZero,
			Size:     size,
		},
		Style: dialog_box_style,
	}

	// box2
	eo.box2 = engine.Box{
		Rectangle: data.Rect{
			Position: data.Vector{X: 2, Y: 25},
			Size:     data.Vector{X: 36, Y: 21},
		},
		Style: inside_box_style,
	}

	// box3
	eo.box3 = engine.Box{
		Rectangle: data.Rect{
			Position: data.Vector{X: 2, Y: 16},
			Size:     data.Vector{X: 36, Y: 9},
		},
		Style: inside_box_style,
	}

	// text fields
	eo.nameTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 4},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.nameTextField.Text = engine.Board.Name.String()

	eo.cameraXTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 6},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 3,
			TextType:  engine.TextCaps,
		},
	}
	eo.cameraXTextField.Text = strconv.Itoa(engine.Board.Camera.X)

	eo.cameraYTextField = engine.TextField{
		Position:   data.Vector{X: 17, Y: 6},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 3,
			TextType:  engine.TextCaps,
		},
	}
	eo.cameraYTextField.Text = strconv.Itoa(engine.Board.Camera.Y)

	eo.charsetTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 10},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.charsetTextField.Text = engine.Board.Charset.String()

	eo.paletteTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 12},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.paletteTextField.Text = engine.Board.Palette.String()

	eo.viewTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 14},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 3,
			TextType:  engine.TextNumbers,
		},
	}
	eo.viewTextField.Text = strconv.Itoa(engine.Board.ViewID)

	eo.northTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 17},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.northTextField.Text = engine.Board.North.String()

	eo.southTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 19},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.southTextField.Text = engine.Board.South.String()

	eo.eastTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 21},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.eastTextField.Text = engine.Board.East.String()

	eo.westTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 23},
		Background: editboardinfo_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.westTextField.Text = engine.Board.West.String()

	// buttons
	eo.counterAddButton = engine.Button{
		Position: data.Vector{X: 13, Y: 26},
		Style:    gray_button_style,
		Text:     " Add ",
	}
	eo.cancelButton = engine.Button{
		Position: data.Vector{X: 1, Y: 48},
		Style:    red_button_style,
		Text:     " X ",
	}
	eo.acceptButton = engine.Button{
		Position: data.Vector{X: 36, Y: 48},
		Style:    green_button_style,
		Text:     " O ",
	}

	// misc
	eo.persistCheckBox = engine.CheckBox{
		Position: data.Vector{X: 13, Y: 8},
		Style:    checkbox_style,
	}
	eo.persistCheckBox.Checked = engine.Board.Flags&data.PersistFlag == data.PersistFlag

	eo.counterControl = CounterControl{
		Position: data.Vector{X: 4, Y: 28},
		Counters: engine.Board.Counters.Copy(),
	}

	// add to container / init
	eo.container.Init()
	eo.container.AddChild(&eo.nameTextField)
	eo.container.AddChild(&eo.cameraXTextField)
	eo.container.AddChild(&eo.cameraYTextField)
	eo.container.AddChild(&eo.charsetTextField)
	eo.container.AddChild(&eo.paletteTextField)
	eo.container.AddChild(&eo.viewTextField)
	eo.container.AddChild(&eo.persistCheckBox)
	//eo.container.AddChild(&eo.darkCheckBox)
	eo.container.AddChild(&eo.northTextField)
	eo.container.AddChild(&eo.southTextField)
	eo.container.AddChild(&eo.eastTextField)
	eo.container.AddChild(&eo.westTextField)
	eo.container.AddChild(&eo.counterAddButton)
	eo.container.AddChild(&eo.counterControl)
	eo.container.AddChild(&eo.cancelButton)
	eo.container.AddChild(&eo.acceptButton)

	// hook up buttons
	eo.counterAddButton.Click.Listen(eo.counterControl.onAddButtonClick)
	eo.acceptButton.Click.Listen(eo.onAcceptButtonClick)
	eo.cancelButton.Click.Listen(eo.onCancelButtonClick)
}

const board_ext = ".sib"
const png_ext = ".png"

// button events
func (eo *EditBoardInfo) onAcceptButtonClick(engine.MouseState) {
	eo.accept()
}

func (eo *EditBoardInfo) accept() {

	// adjust name/path
	name := data.NewName(eo.nameTextField.Text)
	if name != engine.Board.Name {
		if err := RenameBoard(engine.Board.Name, name); err != nil {
			log.Printf("error while renaming board: %v\r\n", err)
			return
		}
		engine.Board.Name = name
	}

	engine.Board.Charset = data.NewName(eo.charsetTextField.Text)
	if err := engine.LoadCharset(engine.Board.Charset); err != nil {
		log.Println(err)
	}

	engine.Board.Palette = data.NewName(eo.paletteTextField.Text)
	if err := engine.LoadPalette(engine.Board.Palette); err != nil {
		log.Println(err)
	}

	if viewId, err := strconv.Atoi(eo.viewTextField.Text); err == nil {
		engine.Board.ViewID = viewId
	}

	engine.Board.North = data.NewName(eo.northTextField.Text)
	engine.Board.South = data.NewName(eo.southTextField.Text)
	engine.Board.East = data.NewName(eo.eastTextField.Text)
	engine.Board.West = data.NewName(eo.westTextField.Text)

	if eo.persistCheckBox.Checked {
		engine.Board.Flags = engine.Board.Flags | data.PersistFlag
	} else {
		engine.Board.Flags = engine.Board.Flags &^ data.PersistFlag
	}

	engine.Board.Counters = eo.counterControl.GetCounters()

	engine.PopStack()
	if eo.Callback != nil {
		eo.Callback()
	}
}

func (eo *EditBoardInfo) onCancelButtonClick(ms engine.MouseState) {
	eo.cancel()
}

func (eo *EditBoardInfo) cancel() {
	engine.PopStack()
	if eo.Callback != nil {
		eo.Callback()
	}
}

// stack
func (eo *EditBoardInfo) ProcessMouse(ms engine.MouseState) {
	ms.Position = ms.Position.Sub(eo.Position)
	eo.container.ProcessMouse(ms)
}

func (eo *EditBoardInfo) ProcessAction(action data.Name) {
	//TD; hotkeys
}

func (eo *EditBoardInfo) Render(buffer *render.Buffer, top bool) {
	offset := eo.Position

	eo.box1.Render(buffer, offset)
	eo.box2.Render(buffer, offset)
	eo.box3.Render(buffer, offset)

	// static text
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 7, Y: 4})), editboardinfo_textcolor, "Name:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 33, Y: 4})), editboardinfo_textcolor, ".sib")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 5, Y: 6})), editboardinfo_textcolor, "Camera:")
	buffer.SetCell(buffer.VectorIndex(offset.Add(data.Vector{X: 16, Y: 6})), data.Cell{Char: ',', Color: editboardinfo_textcolor})
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 4, Y: 8})), editboardinfo_textcolor, "Persist:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 4, Y: 10})), editboardinfo_textcolor, "Charset:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 33, Y: 10})), editboardinfo_textcolor, ".png")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 4, Y: 12})), editboardinfo_textcolor, "Palette:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 33, Y: 12})), editboardinfo_textcolor, ".png")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 7, Y: 14})), editboardinfo_textcolor, "View:")
	//buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 24, Y: 10})), editboardinfo_textcolor, "Dark:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 6, Y: 17})), editboardinfo_textcolor, "North:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 6, Y: 19})), editboardinfo_textcolor, "South:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 7, Y: 21})), editboardinfo_textcolor, "East:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 7, Y: 23})), editboardinfo_textcolor, "West:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 33, Y: 17})), editboardinfo_textcolor, ".sib")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 33, Y: 19})), editboardinfo_textcolor, ".sib")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 33, Y: 21})), editboardinfo_textcolor, ".sib")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 33, Y: 23})), editboardinfo_textcolor, ".sib")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 3, Y: 26})), editboardinfo_textcolor, "Counters:")

	eo.container.Render(buffer, offset)
}

func (eo *EditBoardInfo) Active() {
	engine.FPS(30)
}

func (eo *EditBoardInfo) Process() {
	for _, ms := range engine.PullMouseState() {
		ms.Position = ms.Position.Sub(eo.Position)
		eo.container.ProcessMouse(ms)
	}

	if !engine.CycleTick(4) {
		return
	}
	if action := engine.GetAction(); action != data.DefaultName {
		switch action {
		case engine.INPUT_CANCEL:
			eo.cancel()
		case engine.INPUT_ACCEPT:
			eo.accept()
		}
	}
}

func RenameBoard(old data.Name, new data.Name) error {
	oldfilename := old.String() + data.BoardsEXT
	oldPath := engine.WorldPath.Join(data.BoardsFolder, oldfilename)
	newfilename := new.String() + data.BoardsEXT
	newPath := engine.WorldPath.Join(data.BoardsFolder, newfilename)
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	return nil
}

package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
	"strconv"
)

var editobject_box2_style = data.BoxStyle{
	BorderTopRight:    data.Cell{Char: 191, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderTopLeft:     data.Cell{Char: 218, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottomRight: data.Cell{Char: 217, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottomLeft:  data.Cell{Char: 192, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderTop:         data.Cell{Char: 196, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottom:      data.Cell{Char: 196, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderLeft:        data.Cell{Char: 179, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderRight:       data.Cell{Char: 179, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	Background:        data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorDKGRAY)},
}
var editobject_textcolor = data.NewColor(data.ColorWHITE, data.ColorDKGRAY)
var editobject_textfield_background = data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorBLACK)}

type EditObject struct {
	Position data.Vector
	Prefab   *data.Prefab

	box1 engine.Box
	box2 engine.Box

	// controls
	container         engine.Container
	nameTextField     engine.TextField
	groupTextField    engine.TextField
	cycleTextField    engine.TextField
	scriptTextField   engine.TextField
	targetTextField   engine.TextField
	priorityCheckBox  engine.CheckBox
	showBGCheckBox    engine.CheckBox
	canPushCheckBox   engine.CheckBox
	pushableNCheckBox engine.CheckBox
	pushableSCheckBox engine.CheckBox
	pushableECheckBox engine.CheckBox
	pushableWCheckBox engine.CheckBox
	followCheckBox    engine.CheckBox
	playerCheckBox    engine.CheckBox
	cellView          CellView
	counterAddButton  engine.Button
	collisionBitField BitField
	cancelButton      engine.Button
	acceptButton      engine.Button
	counterControl    CounterControl

	// dialog callback
	Callback func(EditObjectCallback)
}

type EditObjectCallback struct {
	Accepted bool
	Data     data.Prefab
}

// copy data to/from
func (eo *EditObject) CopyObjectData(obj data.Object) {
	eo.CopyPrefabData(obj.Prefab())
}

func (eo *EditObject) CopyPrefabData(pf data.Prefab) {
	eo.nameTextField.Text = pf.Name.String()
	eo.groupTextField.Text = pf.Group.String()
	eo.cycleTextField.Text = strconv.Itoa(int(pf.Cycle))
	eo.scriptTextField.Text = pf.ScriptName.String()
	eo.targetTextField.Text = pf.Target.String()
	eo.cellView.Cell = pf.Cell
	eo.priorityCheckBox.Checked = pf.DrawFlags.Has(data.DrawFlagPriority)
	eo.showBGCheckBox.Checked = pf.DrawFlags.Has(data.DrawFlagShowBG)
	eo.collisionBitField.Value = pf.Collision
	eo.canPushCheckBox.Checked = pf.PushFlags.Has(data.CanPush)
	eo.pushableNCheckBox.Checked = pf.PushFlags.Has(data.PushableNorth)
	eo.pushableSCheckBox.Checked = pf.PushFlags.Has(data.PushableSouth)
	eo.pushableECheckBox.Checked = pf.PushFlags.Has(data.PushableEast)
	eo.pushableWCheckBox.Checked = pf.PushFlags.Has(data.PushableWest)
	eo.followCheckBox.Checked = pf.Follow
	eo.playerCheckBox.Checked = pf.Player

	//counters
	for i := 0; i < len(pf.Counters.Counters) && i < 9; i++ {
		item := &eo.counterControl.Items[i]
		counter := pf.Counters.Counters[i]
		item.enabled = true
		item.nameTextField.Text = counter.Name.String()
		item.valueTextField.Text = strconv.Itoa(counter.Value)
	}
}

func (eo *EditObject) toPrefabData() data.Prefab {
	cycle, _ := strconv.Atoi(eo.cycleTextField.Text)

	prefab := data.Prefab{
		Name:       data.NewName(eo.nameTextField.Text),
		Group:      data.NewName(eo.groupTextField.Text),
		Cycle:      uint8(cycle),
		ScriptName: data.NewName(eo.scriptTextField.Text),
		Target:     data.NewName(eo.targetTextField.Text),
		Cell:       eo.cellView.Cell,
		Collision:  eo.collisionBitField.Value,
	}
	//drawflags
	if eo.priorityCheckBox.Checked {
		prefab.DrawFlags.Set(data.DrawFlagPriority)
	}
	if eo.showBGCheckBox.Checked {
		prefab.DrawFlags.Set(data.DrawFlagShowBG)
	}
	// push flags
	if eo.canPushCheckBox.Checked {
		prefab.PushFlags.Set(data.CanPush)
	}
	if eo.pushableNCheckBox.Checked {
		prefab.PushFlags.Set(data.PushableNorth)
	}
	if eo.pushableSCheckBox.Checked {
		prefab.PushFlags.Set(data.PushableSouth)
	}
	if eo.pushableECheckBox.Checked {
		prefab.PushFlags.Set(data.PushableEast)
	}
	if eo.pushableWCheckBox.Checked {
		prefab.PushFlags.Set(data.PushableWest)
	}
	prefab.Player = eo.playerCheckBox.Checked
	prefab.Follow = eo.followCheckBox.Checked

	//counters
	for i := 0; i < 9; i++ {
		item := &eo.counterControl.Items[i]
		if item.enabled {
			value, _ := strconv.Atoi(item.valueTextField.Text)
			prefab.Counters.SetCounter(
				data.NewName(item.nameTextField.Text),
				value,
			)
		}
	}
	return prefab
}

func (eo *EditObject) Push() {
	eo.init()
	if eo.Prefab != nil {
		eo.CopyPrefabData(*eo.Prefab)
	} else {
		eo.nameTextField.Text = "Object"
		eo.cellView.SetCell(data.Cell{Char: 1, Color: data.ColorWHITE})
	}
	engine.PushStack(eo)
}

func (eo *EditObject) init() {
	size := data.Vector{X: 36, Y: 50}

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
		Title: "Edit Object",
		Rectangle: data.Rect{
			Position: data.VectorZero,
			Size:     size,
		},
		Style: dialog_box_style,
	}

	// box2
	eo.box2 = engine.Box{
		Rectangle: data.Rect{
			Position: data.Vector{X: 1, Y: 25},
			Size:     data.Vector{X: 34, Y: 21},
		},
		Style: editobject_box2_style,
	}

	// text fields
	eo.nameTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 4},
		Background: editobject_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.groupTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 6},
		Background: editobject_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.cycleTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 8},
		Background: editobject_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 3,
			TextType:  engine.TextNumbers,
		},
	}
	eo.scriptTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 10},
		Background: editobject_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}
	eo.targetTextField = engine.TextField{
		Position:   data.Vector{X: 13, Y: 12},
		Background: editobject_textfield_background,
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}

	// buttons
	eo.counterAddButton = engine.Button{
		Position: data.Vector{X: 12, Y: 26},
		Style:    gray_button_style,
		Text:     " Add ",
	}
	eo.cancelButton = engine.Button{
		Position: data.Vector{X: 1, Y: 48},
		Style:    red_button_style,
		Text:     " X ",
	}
	eo.acceptButton = engine.Button{
		Position: data.Vector{X: 32, Y: 48},
		Style:    green_button_style,
		Text:     " O ",
	}

	// misc
	eo.priorityCheckBox = engine.CheckBox{
		Position: data.Vector{X: 13, Y: 16},
		Style:    checkbox_style,
	}
	eo.showBGCheckBox = engine.CheckBox{
		Position: data.Vector{X: 26, Y: 16},
		Style:    checkbox_style,
	}
	eo.canPushCheckBox = engine.CheckBox{
		Position: data.Vector{X: 13, Y: 18},
		Style:    checkbox_style,
	}
	eo.pushableNCheckBox = engine.CheckBox{
		Position: data.Vector{X: 26, Y: 18},
		Style: engine.CheckBoxStyle{
			CheckedOn:  data.Cell{Char: 'N', Color: data.ColorWHITE},
			CheckedOff: data.Cell{Char: 'N', Color: data.ColorDKGRAY},
		},
	}
	eo.pushableSCheckBox = engine.CheckBox{
		Position: data.Vector{X: 27, Y: 18},
		Style: engine.CheckBoxStyle{
			CheckedOn:  data.Cell{Char: 'S', Color: data.ColorWHITE},
			CheckedOff: data.Cell{Char: 'S', Color: data.ColorDKGRAY},
		},
	}
	eo.pushableECheckBox = engine.CheckBox{
		Position: data.Vector{X: 28, Y: 18},
		Style: engine.CheckBoxStyle{
			CheckedOn:  data.Cell{Char: 'E', Color: data.ColorWHITE},
			CheckedOff: data.Cell{Char: 'E', Color: data.ColorDKGRAY},
		},
	}
	eo.pushableWCheckBox = engine.CheckBox{
		Position: data.Vector{X: 29, Y: 18},
		Style: engine.CheckBoxStyle{
			CheckedOn:  data.Cell{Char: 'W', Color: data.ColorWHITE},
			CheckedOff: data.Cell{Char: 'W', Color: data.ColorDKGRAY},
		},
	}
	eo.followCheckBox = engine.CheckBox{
		Position: data.Vector{X: 13, Y: 20},
		Style:    checkbox_style,
	}
	eo.playerCheckBox = engine.CheckBox{
		Position: data.Vector{X: 26, Y: 20},
		Style:    checkbox_style,
	}
	eo.cellView = CellView{
		Position:   data.Vector{X: 13, Y: 22},
		ShowHover:  true,
		HoverColor: data.NewColor(data.ColorBLACK, data.ColorYELLOW),
	}
	eo.collisionBitField = BitField{
		Position: data.Vector{X: 13, Y: 14},
	}
	eo.counterControl = CounterControl{
		Position: data.Vector{X: 2, Y: 28},
	}

	// add to container / init
	eo.container.Init()
	eo.container.AddChild(&eo.nameTextField)
	eo.container.AddChild(&eo.groupTextField)
	eo.container.AddChild(&eo.cycleTextField)
	eo.container.AddChild(&eo.scriptTextField)
	eo.container.AddChild(&eo.targetTextField)
	eo.container.AddChild(&eo.counterAddButton)
	eo.container.AddChild(&eo.priorityCheckBox)
	eo.container.AddChild(&eo.showBGCheckBox)
	eo.container.AddChild(&eo.canPushCheckBox)
	eo.container.AddChild(&eo.pushableNCheckBox)
	eo.container.AddChild(&eo.pushableSCheckBox)
	eo.container.AddChild(&eo.pushableECheckBox)
	eo.container.AddChild(&eo.pushableWCheckBox)
	eo.container.AddChild(&eo.followCheckBox)
	eo.container.AddChild(&eo.playerCheckBox)
	eo.container.AddChild(&eo.cellView)
	eo.container.AddChild(&eo.collisionBitField)
	eo.container.AddChild(&eo.counterControl)
	eo.container.AddChild(&eo.cancelButton)
	eo.container.AddChild(&eo.acceptButton)

	// hook up buttons
	eo.counterAddButton.Click.Listen(eo.counterControl.onAddButtonClick)
	eo.acceptButton.Click.Listen(eo.onAcceptButtonClick)
	eo.cancelButton.Click.Listen(eo.onCancelButtonClick)
	eo.cellView.Click.Listen(eo.onCellViewClick)
}

// button events
func (eo *EditObject) onAcceptButtonClick(ms engine.MouseState) {
	eo.accept()
}

func (eo *EditObject) accept() {
	eo.Callback(EditObjectCallback{
		Accepted: true,
		Data:     eo.toPrefabData(),
	})
	engine.PopStack()
}

func (eo *EditObject) onCancelButtonClick(ms engine.MouseState) {
	eo.cancel()
}

func (eo *EditObject) cancel() {
	eo.Callback(EditObjectCallback{
		Accepted: false,
	})
	engine.PopStack()
}

func (eo *EditObject) onCellViewClick(ms engine.MouseState) {
	if !ms.LeftButtonClick {
		return
	}

	ec := EditCell{
		Callback: eo.onEditCellCallback,
		Cell:     eo.cellView.Cell,
	}
	ec.Push()
}

func (eo *EditObject) onEditCellCallback(callback EditCellCallback) {
	if callback.Accepted {
		eo.cellView.SetCell(callback.Cell)
	}
}

// stack
func (eo *EditObject) ProcessMouse(ms engine.MouseState) {
	ms.Position = ms.Position.Sub(eo.Position)
	eo.container.ProcessMouse(ms)
}

func (eo *EditObject) ProcessAction(action data.Name) {
	//TD; hotkeys
}

func (eo *EditObject) Render(buffer *render.Buffer, top bool) {
	offset := eo.Position

	eo.box1.Render(buffer, offset)
	eo.box2.Render(buffer, offset)

	// static text
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 7, Y: 4})), editobject_textcolor, "Name:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 6, Y: 6})), editobject_textcolor, "Group:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 6, Y: 8})), editobject_textcolor, "Cycle:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 5, Y: 10})), editobject_textcolor, "Script:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 5, Y: 12})), editobject_textcolor, "Target:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 2, Y: 14})), editobject_textcolor, "Collision:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 3, Y: 16})), editobject_textcolor, "Priority:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 17, Y: 16})), editobject_textcolor, "Show BG:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 3, Y: 18})), editobject_textcolor, "Can Push:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 16, Y: 18})), editobject_textcolor, "Pushable:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 5, Y: 20})), editobject_textcolor, "Follow:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 18, Y: 20})), editobject_textcolor, "Player:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 7, Y: 22})), editobject_textcolor, "Cell:")
	buffer.WriteText(buffer.VectorIndex(offset.Add(data.Vector{X: 2, Y: 26})), editobject_textcolor, "Counters:")

	eo.container.Render(buffer, offset)
}

func (eo *EditObject) Active() {
	engine.FPS(30)
}

func (eo *EditObject) Process() {
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

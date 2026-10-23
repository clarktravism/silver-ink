package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
	"strconv"
)

const max_counters_items = 9

type CounterControl struct {
	Position  data.Vector
	Counters  data.Counters
	container engine.Container
	Items     [max_counters_items]CounterControlItem
	rect      data.Rect
}

func (cc *CounterControl) Init() {

	cc.rect = data.Rect{
		Position: cc.Position,
		Size:     data.Vector{X: 33, Y: 20},
	}

	cc.container = engine.Container{
		Rectangle: cc.rect,
	}
	cc.container.Init()

	for i := range cc.Items {
		cc.Items[i] = CounterControlItem{
			Index:    i,
			Position: data.Vector{X: 0, Y: i * 2},
		}
		cc.container.AddChild(&cc.Items[i])
		cc.Items[i].Disabled.Listen(cc.onItemDisabled)
	}
	cc.setCounters(cc.Counters)
}

func (cc *CounterControl) Rect() data.Rect {
	return cc.rect
}

func (cc *CounterControl) Render(buffer *render.Buffer, offset data.Vector) {
	cc.container.Render(buffer, offset)
}

func (cc *CounterControl) ProcessMouse(ms engine.MouseState) {
	cc.container.ProcessMouse(ms)
}

// when the 'Add' button is clicked
func (cc *CounterControl) onAddButtonClick(ms engine.MouseState) {
	if !ms.LeftButtonClick {
		return
	}

	for i := 0; i < len(cc.Items); i++ {
		if !cc.Items[i].enabled {
			cc.Items[i].Enable()
			cc.Items[i].nameTextField.Text = "COUNTER" + strconv.Itoa(i+1)
			cc.Items[i].valueTextField.Text = "0"
			return
		}
	}
}

// when the 'X' button is clicked on the counter
func (cc *CounterControl) onItemDisabled(arg CounterControlItemArgs) {
	// check if there is an enabled item below
	if i := arg.Index + 1; i < len(cc.Items) && cc.Items[i].enabled {
		cc.Items[arg.Index].nameTextField.Text = cc.Items[i].nameTextField.Text
		cc.Items[arg.Index].valueTextField.Text = cc.Items[i].valueTextField.Text
		cc.Items[i].Disable()
		cc.Items[arg.Index].Enable()
	}
}

func (cc *CounterControl) setCounters(counters data.Counters) {
	for i := 0; i < len(counters.Counters) && i < max_counters_items; i++ {
		counter := counters.Counters[i]
		cc.Items[i].Enable()
		cc.Items[i].nameTextField.Text = counter.Name.String()
		cc.Items[i].valueTextField.Text = strconv.Itoa(counter.Value)
	}
}

func (cc *CounterControl) GetCounters() data.Counters {
	counters := data.Counters{}
	for i := 0; i < len(cc.Items); i++ {
		if cc.Items[i].enabled {
			value, _ := strconv.Atoi(cc.Items[i].valueTextField.Text)
			counters.SetCounter(
				data.NewName(cc.Items[i].nameTextField.Text),
				value,
			)
		} else {
			break
		}
	}
	return counters
}

// CounterControlItem
var counter_control_item_background = data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorDKGRAY)}

type CounterControlItem struct {
	Index     int
	Position  data.Vector
	container engine.Container

	Disabled engine.Signal[CounterControlItemArgs]

	nameTextField  engine.TextField
	valueTextField engine.TextField
	removeButton   engine.Button
	enabled        bool
	rect           data.Rect
}

type CounterControlItemArgs struct {
	Index int
}

func (item *CounterControlItem) Init() {
	item.rect = data.Rect{
		Position: item.Position,
		Size:     data.Vector{X: 32, Y: 1},
	}
	item.container.Rectangle = data.Rect{
		Position: data.VectorZero,
		Size:     item.rect.Size,
	}
	item.container.Init()

	item.nameTextField = engine.TextField{
		Position:   data.Vector{X: 0, Y: 0},
		Background: data.Cell{Char: 0, Color: data.ColorWHITE},
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 20,
			TextType:  engine.TextCaps,
		},
	}

	item.valueTextField = engine.TextField{
		Position:   data.Vector{X: 21, Y: 0},
		Background: data.Cell{Char: 0, Color: data.ColorWHITE},
		Cursor:     data.Cell{Char: '_', Color: data.ColorGRAY},
		TextInput: engine.TextInput{
			MaxLength: 9,
			TextType:  engine.TextNumbers,
		},
	}

	// for tabs
	item.nameTextField.TextInput.Next = &item.valueTextField.TextInput
	item.valueTextField.TextInput.Next = &item.nameTextField.TextInput

	item.removeButton = engine.Button{
		Position: data.Vector{X: 31, Y: 0},
		Style:    red_button_style,
		Text:     "X",
	}
	item.removeButton.Click.Listen(item.onRemoveButtonClick)

	item.container.Init()
	item.container.AddChild(&item.nameTextField)
	item.container.AddChild(&item.valueTextField)
	item.container.AddChild(&item.removeButton)
}

func (item *CounterControlItem) Rect() data.Rect {
	return item.rect
}

func (item *CounterControlItem) Render(buffer *render.Buffer, offset data.Vector) {
	offset = offset.Add(item.Position)
	if item.enabled {
		item.container.Render(buffer, offset)
		buffer.SetCell(buffer.Index(offset.X+20, offset.Y), data.Cell{Char: ':', Color: data.NewColor(data.ColorWHITE, data.ColorDKGRAY)})
	} else {
		buffer.FillRectangle(buffer.Index(offset.X, offset.Y), item.rect.Size, counter_control_item_background)
	}
}

func (item *CounterControlItem) ProcessMouse(ms engine.MouseState) {
	ms.Position = ms.Position.Sub(item.Position)
	item.container.ProcessMouse(ms)
}

func (item *CounterControlItem) Enable() {
	item.enabled = true
	engine.Dirty()
}

func (item *CounterControlItem) onRemoveButtonClick(ms engine.MouseState) {
	if ms.LeftButtonClick {
		item.Disable()
	}
}

func (item *CounterControlItem) Disable() {
	item.enabled = false
	item.nameTextField.Text = ""
	item.valueTextField.Text = ""
	item.nameTextField.LoseFocus()
	item.valueTextField.LoseFocus()
	item.Disabled.Emit(CounterControlItemArgs{Index: item.Index})
	engine.Dirty()
}

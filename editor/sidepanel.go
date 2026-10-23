package main

import (
	"fmt"
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

const sidepanel_x = 62
const sidepanel_width = 18
const sidepanel_height = 60

var sidepanel_background = data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorDKGRAY)}
var sidepanel_button_style = engine.LabelButtonStyle{
	LeftBackground:       data.NewColor(data.ColorBLACK, data.ColorGRAY),
	LeftHoverBackground:  data.NewColor(data.ColorBLACK, data.ColorYELLOW),
	RightBackground:      data.NewColor(data.ColorWHITE, data.ColorDKGRAY),
	RightHoverBackground: data.NewColor(data.ColorYELLOW, data.ColorDKGRAY),
}
var sidepanel_box_style = data.BoxStyle{
	BorderTopRight:    data.Cell{Char: 191, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderTopLeft:     data.Cell{Char: 218, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottomRight: data.Cell{Char: 217, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottomLeft:  data.Cell{Char: 192, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderTop:         data.Cell{Char: 196, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottom:      data.Cell{Char: 196, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderLeft:        data.Cell{Char: 179, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderRight:       data.Cell{Char: 179, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	Background:        sidepanel_background,
}
var sidepanel_boxpanel_color = data.NewColor(data.ColorGRAY, data.ColorDKGRAY)
var sidepanel_label_color = data.NewColor(data.ColorWHITE, data.ColorDKGRAY)

type sidePanel struct {
	engine.Container

	// label boxes
	toolLabelBox   engine.LabelBox
	charLabelBox   engine.LabelBox
	colorLabelBox  engine.LabelBox
	bufferLabelBox engine.LabelBox
	infoLabelBox   engine.LabelBox

	// labels
	worldLabel     engine.Label
	boardLabel     engine.Label
	infoCoordLabel engine.Label

	// buttons
	loadButton      engine.LabelButton
	saveButton      engine.LabelButton
	quitButton      engine.LabelButton
	worldinfoButton engine.LabelButton
	boardsButton    engine.LabelButton
	boardinfoButton engine.LabelButton
	viewButton      engine.LabelButton
	objectButton    engine.LabelButton
	prefabButton    engine.LabelButton
	elementButton   engine.LabelButton
	drawButton      engine.LabelButton
	fillButton      engine.LabelButton
	textButton      engine.LabelButton
	remobjButton    engine.LabelButton

	// misc
	charSelect    CharSelect
	colorSelect   ColorSelect
	bufferTool    BufferTool
	toolBrushView BrushView
	infoBrushView BrushView
}

func (sp *sidePanel) Init() {

	// container
	//sp.Container = engine.Container{
	//	Rectangle: data.NewRect(sidepanel_x, 0, sidepanel_width, sidepanel_height),
	//}
	sp.Container.Rectangle = data.NewRect(sidepanel_x, 0, sidepanel_width, sidepanel_height)
	sp.Container.Init()

	// label boxes
	sp.toolLabelBox = engine.LabelBox{
		Rectangle: data.NewRect(0, 18, sidepanel_width, 5),
		Text:      " Tool:Draw ",
		Color:     sidepanel_boxpanel_color,
		Style:     sidepanel_box_style,
	}

	sp.charLabelBox = engine.LabelBox{
		Rectangle: data.NewRect(0, 28, sidepanel_width, 18),
		Text:      " Char ",
		Color:     sidepanel_boxpanel_color,
		Style:     sidepanel_box_style,
	}

	sp.colorLabelBox = engine.LabelBox{
		Rectangle: data.NewRect(0, 46, sidepanel_width, 4),
		Text:      " Color ",
		Color:     sidepanel_boxpanel_color,
		Style:     sidepanel_box_style,
	}

	sp.bufferLabelBox = engine.LabelBox{
		Rectangle: data.NewRect(0, 50, sidepanel_width, 4),
		Text:      " Buffer ",
		Color:     sidepanel_boxpanel_color,
		Style:     sidepanel_box_style,
	}

	sp.infoLabelBox = engine.LabelBox{
		Rectangle: data.NewRect(0, 54, sidepanel_width, 6),
		Text:      " Info ",
		Color:     sidepanel_boxpanel_color,
		Style:     sidepanel_box_style,
	}

	// labels
	sp.worldLabel = engine.Label{
		Position: data.Vector{X: 1, Y: 1},
		Text:     "",
		Color:    sidepanel_label_color,
		Length:   16,
	}
	sp.boardLabel = engine.Label{
		Position: data.Vector{X: 1, Y: 7},
		Text:     "Untitled Board",
		Color:    sidepanel_label_color,
		Length:   16,
	}
	sp.infoCoordLabel = engine.Label{
		Position: data.Vector{X: 1, Y: 58},
		Text:     "26,14",
		Color:    sidepanel_label_color,
		Length:   16,
	}

	// buttons
	sp.loadButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 4},
		Style:     sidepanel_button_style,
		LeftText:  " L ",
		RightText: " Load World ",
	}
	sp.quitButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 5},
		Style:     sidepanel_button_style,
		LeftText:  " Q ",
		RightText: " Quit ",
	}
	sp.worldinfoButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 3},
		Style:     sidepanel_button_style,
		LeftText:  " W ",
		RightText: " World Info ",
	}
	sp.saveButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 9},
		Style:     sidepanel_button_style,
		LeftText:  " S ",
		RightText: " Save ",
	}
	sp.boardsButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 10},
		Style:     sidepanel_button_style,
		LeftText:  " B ",
		RightText: " Boards ",
	}
	sp.boardinfoButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 11},
		Style:     sidepanel_button_style,
		LeftText:  " I ",
		RightText: " Board Info ",
	}
	sp.viewButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 13},
		Style:     sidepanel_button_style,
		LeftText:  " V ",
		RightText: " View: All ",
	}
	sp.objectButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 14},
		Style:     sidepanel_button_style,
		LeftText:  " O ",
		RightText: " Object ",
	}
	sp.prefabButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 15},
		Style:     sidepanel_button_style,
		LeftText:  " P ",
		RightText: " Prefab ",
	}
	sp.elementButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 16},
		Style:     sidepanel_button_style,
		LeftText:  " E ",
		RightText: " Element ",
	}
	sp.drawButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 23},
		Style:     sidepanel_button_style,
		LeftText:  " D ",
		RightText: " Draw ",
	}
	sp.fillButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 24},
		Style:     sidepanel_button_style,
		LeftText:  " F ",
		RightText: " Fill ",
	}
	sp.textButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 25},
		Style:     sidepanel_button_style,
		LeftText:  " T ",
		RightText: " Text ",
	}
	sp.remobjButton = engine.LabelButton{
		Position:  data.Vector{X: 1, Y: 26},
		Style:     sidepanel_button_style,
		LeftText:  " R ",
		RightText: " Remove Obj. ",
	}

	// misc
	sp.toolBrushView = BrushView{
		Position: data.Vector{X: 1, Y: 19},
	}

	sp.charSelect = CharSelect{
		Position: data.Vector{X: 1, Y: 29},
		Char:     1,
	}
	sp.colorSelect = ColorSelect{
		Position: data.Vector{X: 1, Y: 47},
		Color:    data.ColorWHITE,
	}
	sp.bufferTool = BufferTool{
		Position: data.Vector{X: 1, Y: 51},
	}
	sp.infoBrushView = BrushView{
		Position: data.Vector{X: 1, Y: 55},
	}

	// add to container for rendering and mouse processing

	// label boxes
	sp.AddChild(&sp.toolLabelBox)
	sp.AddChild(&sp.charLabelBox)
	sp.AddChild(&sp.colorLabelBox)
	sp.AddChild(&sp.bufferLabelBox)
	sp.AddChild(&sp.infoLabelBox)

	// labels
	sp.AddChild(&sp.worldLabel)
	sp.AddChild(&sp.boardLabel)
	sp.AddChild(&sp.infoCoordLabel)

	// buttons
	sp.AddChild(&sp.loadButton)
	sp.AddChild(&sp.quitButton)
	sp.AddChild(&sp.worldinfoButton)
	sp.AddChild(&sp.saveButton)
	sp.AddChild(&sp.boardsButton)
	sp.AddChild(&sp.boardinfoButton)
	sp.AddChild(&sp.viewButton)
	sp.AddChild(&sp.objectButton)
	sp.AddChild(&sp.prefabButton)
	sp.AddChild(&sp.elementButton)
	sp.AddChild(&sp.drawButton)
	sp.AddChild(&sp.fillButton)
	sp.AddChild(&sp.textButton)
	sp.AddChild(&sp.remobjButton)

	//misc
	sp.AddChild(&sp.toolBrushView)
	sp.AddChild(&sp.charSelect)
	sp.AddChild(&sp.colorSelect)
	sp.AddChild(&sp.bufferTool)
	sp.AddChild(&sp.infoBrushView)
}

func (sp *sidePanel) Rect() data.Rect {
	return sp.Container.Rectangle
}

func (sp *sidePanel) Render(buffer *render.Buffer, offset data.Vector) {
	rectPos := offset.Add(sp.Container.Rectangle.Position)
	i := buffer.Index(rectPos.X, rectPos.Y)
	buffer.FillRectangle(i, sp.Container.Rectangle.Size, sidepanel_background)
	sp.Container.Render(buffer, offset)
}

// signal implementations
func (sp *sidePanel) setInfoCoords(value data.Vector) {
	if value.X == -1 && value.Y == -1 {
		sp.infoCoordLabel.SetText("           ")
		return
	}
	sp.infoCoordLabel.SetText(fmt.Sprintf("%v,%v    ", value.X, value.Y))
}

func (sp *sidePanel) SetBrush(brush Brush) {
	sp.toolBrushView.SetBrush(brush)
	sp.charSelect.SetChar(brush.Cell.Char)
	sp.colorSelect.SetColor(brush.Cell.Color)
	sp.charLabelBox.Text = fmt.Sprintf(" Char %v ", brush.Cell.Char)
	sp.colorLabelBox.Text = fmt.Sprintf(" Color x%02X ", brush.Cell.Color)
}

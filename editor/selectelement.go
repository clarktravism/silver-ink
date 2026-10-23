package main

import (
	"log"
	"silver-ink/data"
	"silver-ink/engine"
)

const elements_file = "elements.json"

type SelectElement struct {
	Callback func(SelectElementCallback)
	Color    data.Color
}

type SelectElementCallback struct {
	Accepted bool
	Element  data.Element
}

func (se *SelectElement) Push() {

	if engine.ElementsFile.HasChanged() {
		if err := engine.ElementsFile.Reload(); err != nil {
			log.Printf("error while reading elements: %v\r\n", err)
		}
	}

	// create scroll box
	size := data.Vector{X: 38, Y: 41}
	sb := engine.ScrollBox{
		Rectangle: data.Rect{
			Position: engine.CenterPos(size),
			Size:     size,
		},
		Title:       "Select Element",
		Style:       scroll_box_style,
		Callback:    se.callbackScrollBox,
		ForceSelect: false,
	}

	for _, element := range engine.ElementsFile.Elements {
		if !element.Name.Empty() {
			sb.AppendLine(engine.ScrollBoxLine{
				Text:     element.Name.String(),
				Link:     SELECT_ELEMENT,
				Ext:      element,
				Centered: false,
				Icon:     data.Cell{Char: element.Char, Color: se.Color},
			})
		}
	}

	sb.AppendText("", false)
	sb.AppendLink("Cancel", CANCEL, nil)
	sb.Push()
}

var SELECT_ELEMENT = data.NewName("SELECT_ELEMENT")
var CANCEL = data.NewName("CANCEL")

func (se *SelectElement) callbackScrollBox(callback engine.ScrollBoxCallback) {
	if !callback.Accepted || callback.Link == CANCEL {
		se.Callback(SelectElementCallback{
			Accepted: false,
		})
		return
	}

	if callback.Link == SELECT_ELEMENT {
		element, _ := callback.Ext.(data.Element)

		se.Callback(SelectElementCallback{
			Accepted: true,
			Element:  element,
		})
	}
}

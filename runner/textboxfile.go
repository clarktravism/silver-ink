package main

import (
	"encoding/json"
	"log"
	"os"
	"silver-ink/data"
	"silver-ink/engine"
)

const textboxes_file = "textboxes.json"

var TextBoxFile textBoxFile

type textBoxFile struct {
	DialogID int `json:"DialogID"`
	SaveGameID int `json:"SaveGameID"`
	EndGameID int `json:"EndGameID"`
	TextBoxes []TextBoxStyle `json:"TextBoxes"`
}

type TextBoxStyle struct {
	ID int `json:"ID"`
	Width int `json:"Width"`
	Height int `json:"Height"`
	data.ScrollBoxStyle
}

func LoadTextboxes(wp *data.WorldPath) error {
	path := wp.Join(textboxes_file)
	// read and decode
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &TextBoxFile); err != nil {
		return err
	}
	return nil
}

func GetTextBoxStyle(id int) *TextBoxStyle {
	for i := range TextBoxFile.TextBoxes {
		if id == TextBoxFile.TextBoxes[i].ID {
			return &TextBoxFile.TextBoxes[i]
		}
	}
	return nil
}

func CreateTextBox(id int, title string) *engine.ScrollBox {
	style := GetTextBoxStyle(id)
	if style == nil {
		log.Printf("textbox id '%v' not found\r\n", id)
		return nil
	}
	size := data.Vector{X: style.Width, Y: style.Height}
	sb := engine.ScrollBox{
		Rectangle:  data.Rect{
			Position: engine.CenterPos(size),
			Size: size,
		},
		Title: title,
		Style: style.ScrollBoxStyle,
	}
	return &sb
}


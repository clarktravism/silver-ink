package main

import (
	"log"
	"silver-ink/data"
	"silver-ink/engine"
)

var LOAD = data.NewName("LOAD")
var TITLE = data.NewName("TITLE")

func PushGameOver() {
	scrollbox := CreateTextBox(TextBoxFile.EndGameID, "Game Over")
	scrollbox.Callback = callbackGameOver
	for _, line := range Text.lines {
		scrollbox.AppendLine(line)
	}
	// blank line
	if len(Text.lines) > 0 {
		scrollbox.AppendLine(engine.ScrollBoxLine{})
	}

	ResetText()

	scrollbox.AppendLink("Load Save File", LOAD, nil)
	scrollbox.AppendLink("Return to Title", TITLE, nil)

	TextBoxIsOpen = true
	scrollbox.Push()
}

func callbackGameOver(callback engine.ScrollBoxCallback) {
	TextBoxIsOpen = false
	if callback.Accepted && callback.Link == LOAD {
		textBox := CreateTextBox(TextBoxFile.SaveGameID, "Load Game")
		textBox.Callback = callbackLoadGameOver
		filenames := getSaveFiles()
		for _, filename := range filenames {
			textBox.AppendLink(filename, FILE, filename)
		}
		textBox.Push()
	} else {
		StartNewGame("")
	}
}

func callbackLoadGameOver(callback engine.ScrollBoxCallback) {
	if !callback.Accepted {
		StartNewGame("")
	}
	if callback.Link == FILE {
		name := callback.Ext.(string)
		if err := LoadSave(name); err != nil {
			log.Println(err)
			StartNewGame("")
		}
	}
}
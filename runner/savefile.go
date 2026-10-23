package main

import (
	"log"
	"os"
	"path/filepath"
	"silver-ink/data"
	"silver-ink/engine"
)

var SaveFile *data.SaveFile

func SaveBoard(board *data.Board) error {
	if SaveFile == nil {
		SaveFile = &data.SaveFile{}
	}
	if engine.Board == nil {
		return nil
	}
	if engine.Board.Name == engine.WorldFile.TitleBoard {
		return nil
	}
	return SaveFile.SaveBoard(board)
}

func SaveGame(name string) error {
	if SaveFile == nil {
		SaveFile = &data.SaveFile{}
	}
	if err := engine.WorldPath.WriteSaveFile(name, SaveFile, engine.Board); err != nil {
		return err
	}
	return nil
}

func LoadSave(name string) error {
	sf, err := engine.WorldPath.ReadSaveFile(name)
	if err != nil {
		return err
	}
	SaveFile = sf
	engine.FPS(SaveFile.FPS)

	board, err := SaveFile.ReadBoard(SaveFile.CurrentBoard)
	if err != nil {
		return err
	}

	if err := SetBoard(board); err != nil {
		return err
	}
	return nil
}

const saves_folder = "saves"
const save_ext = ".sav"

func getSaveFiles() []string {
	savesPath := engine.WorldPath.Join(saves_folder)
	entries, err := os.ReadDir(savesPath)
	if err != nil {
		log.Println(err)
		return nil
	}
	result := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != save_ext {
			continue
		}
		result = append(result, entry.Name())
	}
	return result
}

var FILE = data.NewName("FILE")
var NEW = data.NewName("NEW")

func OpenSavePrompt() {
	textBox := CreateTextBox(TextBoxFile.SaveGameID, "Save Game")
	textBox.Callback = callbackSavePrompt
	filenames := getSaveFiles()
	for _, filename := range filenames {
		textBox.AppendLink(filename, FILE, filename)
	}
	textBox.AppendLine(engine.ScrollBoxLine{})
	textBox.AppendLink("New Save", NEW, nil)
	textBox.Push()
}

func callbackSavePrompt(callback engine.ScrollBoxCallback) {
	if !callback.Accepted {
		return
	}
	if callback.Link == FILE {
		name := callback.Ext.(string)
		if err := SaveGame(name); err != nil {
			log.Println(err)
			return
		}
		return
	}
	if callback.Link == NEW {
		tp := engine.TextPrompt{
			Title: "New Save",
			Length: 20,
			AllowCancel: true,
			Callback: callbackNewSave,
			Style: text_prompt_style,
		}
		tp.Push()
		return
	}
}

func callbackNewSave(callback engine.TextPromptCallback) {
	if !callback.Accepted {
		return
	}
	if err := SaveGame(callback.Text + save_ext); err != nil {
		log.Println(err)
		return
	}
}

func OpenLoadPrompt() {
	textBox := CreateTextBox(TextBoxFile.SaveGameID, "Load Game")
	textBox.Callback = callbackLoadPrompt
	filenames := getSaveFiles()
	for _, filename := range filenames {
		textBox.AppendLink(filename, FILE, filename)
	}
	textBox.Push()
}

func callbackLoadPrompt(callback engine.ScrollBoxCallback) {
	if !callback.Accepted {
		return
	}
	if callback.Link == FILE {
		name := callback.Ext.(string)
		if err := LoadSave(name); err != nil {
			log.Println(err)
			return
		}
	}
}

var text_prompt_style = data.BoxStyle{
	BorderTopRight:        data.Cell{Char: 191, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderTopLeft:         data.Cell{Char: 218, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottomRight:     data.Cell{Char: 217, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottomLeft:      data.Cell{Char: 192, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderTop:             data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottom:          data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderLeft:            data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderRight:           data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivide:      data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivideLeft:  data.Cell{Char: 195, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivideRight: data.Cell{Char: 180, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	Background:            data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorBLACK)},
	TitleBackground:       data.Cell{Char: 0, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
}
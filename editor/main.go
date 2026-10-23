package main

import (
	"fmt"
	"log"
	"silver-ink/engine"
	"silver-ink/data"
)

const editor_width = 80
const editor_height = 60

var CHR_DEFAULT = data.NewName("CHR_DEFAULT")
var PAL_DEFAULT = data.NewName("PAL_DEFAULT")

func main() {
	// log setup
	engine.SetLog("editor.log")
	defer engine.CloseLog()

	engine.Setup("Silver-Ink Editor", editor_width, editor_height, 30)
	defer engine.Close()

	// load default charset
	if err := engine.LoadCharset(CHR_DEFAULT); err != nil {
		log.Println(err)
		return
	}
	if err := engine.LoadPalette(PAL_DEFAULT); err != nil {
		log.Println(err)
		return
	}
	engine.Dirty()

	engine.LoadKeyBinds("./keybinds.json")

	engine.SetQuitCallback(onQuit)

	sw := SelectWorld{
		Force: true,
		Callback: func(callback SelectWorldCallback) {
			LoadWorld(callback.WorldPath, callback.WorldName)
			eb := EditBoard{}
			eb.Push()
		},
	}
	sw.Push()
	engine.Run()
}

func LoadWorld(worldPath string, worldName string) {
	if err := engine.SetWorldPath(worldPath); err != nil {
		log.Printf("error while loading world: %v\r\n", err)
		return
	}
	engine.Title(fmt.Sprintf("Silver Ink Editor - %v", worldName))

	// resolve board
	if engine.WorldFile.StartBoard.Empty() {
		log.Println("error while loading world: no StartBoard")
		return
	}

	// load board
	board, err := engine.WorldPath.LoadBoard(engine.WorldFile.StartBoard)
	if err != nil {
		log.Printf("error while loading board: %v\r\n", err)
		return
	}
	engine.SetBoard(board)
	engine.Dirty()
}

func onQuit() {
	SaveBoard()
	engine.Quit()
}

func SaveBoard() {
	if engine.Board != nil {
		if err := engine.WorldPath.SaveBoard(engine.Board); err != nil {
			log.Printf("error while saving board: %v\r\n", err)
		}
	}
}

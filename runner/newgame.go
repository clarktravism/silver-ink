package main

import (
	"log"
	"silver-ink/data"
	"silver-ink/engine"
)

func StartNewGame(boardName string) {
	name := engine.WorldFile.TitleBoard
	if boardName != "" {
		name = data.NewName(boardName)
	}

	SaveFile = &data.SaveFile{}
	SaveFile.CopyWorldData(engine.WorldFile)

	board, err := ReadBoardFile(name)
	board.Enter = data.Vector{X: -1, Y: -1}
	if err != nil {
		log.Println(err)
		return
	}
	if err := SetBoard(board); err != nil {
		log.Println(err)
		return
	}
}
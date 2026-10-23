package main

import (
	"strings"
	"os"
	"log"
	"path/filepath"
	"silver-ink/data"
	"silver-ink/engine"
)

type SelectBoard struct {
	Callback func(SelectBoardCallback)
}

type SelectBoardCallback struct {
	Accepted bool
	Name data.Name
}

var OPEN_BOARD = data.NewName("OPEN BOARD")
var NEW_BOARD = data.NewName("NEW_BOARD")

func (sb *SelectBoard) Push() {
	sb.pushScrollBox()
}

func (sb *SelectBoard) pushScrollBox() {
	size := data.Vector{ X:38, Y: 25 }

	sbox := engine.ScrollBox{
		Rectangle:  data.Rect{
			Position: engine.CenterPos(size),
			Size: size,
		},
		Title: "Select Board",
		Style: scroll_box_style,
		Callback: sb.callbackSelectBoard,
	}

	names := getBoardNames()
	for _, name := range names {
		sbox.AppendLink(name.String(), OPEN_BOARD, name)
	}

	sbox.AppendText("", false)
	sbox.AppendLink("Create New Board", NEW_BOARD, nil)
	sbox.AppendLink("Cancel", CANCEL, nil)
	sbox.Push()
}

func getBoardNames() []data.Name {
	path := engine.WorldPath.Join(data.BoardsFolder)
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	result := make([]data.Name, 0)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != data.BoardsEXT {
			continue
		}
		name := data.NewName(strings.TrimSuffix(entry.Name(), data.BoardsEXT))
		result = append(result, name)
	}
	return result
}

func (sb *SelectBoard) callbackSelectBoard(callback engine.ScrollBoxCallback) {
	
	if !callback.Accepted || callback.Link == CANCEL {
		sb.Callback(SelectBoardCallback{
			Accepted: false,
		})
	}

	if callback.Link == NEW_BOARD {
		sb.pushNewBoard()
		return
	}

	if callback.Link == OPEN_BOARD {
		name, _ := callback.Ext.(data.Name)
		sb.Callback(SelectBoardCallback{
			Accepted: true,
			Name: name,
		})
	}
}

func (sb *SelectBoard) pushNewBoard() {
	tp := engine.TextPrompt{
		Title: "New Board Name",
		Length: 20,
		AllowCancel: true,
		Callback: sb.callbackNewBoard,
		Style: text_prompt_style,
	}
	tp.Push()
}

func (sb *SelectBoard) callbackNewBoard(callback engine.TextPromptCallback) {
	if !callback.Accepted || callback.Text == "" {
		sb.pushScrollBox()
		return
	}

	name := data.NewName(callback.Text)
	board := data.NewBoard()
	board.Name = name

	if err := engine.WorldPath.SaveBoard(&board); err != nil {
		log.Println(err)
	}

	sb.Callback(SelectBoardCallback{
		Accepted: true,
		Name: name,
	})
}
package main

import (
	"log"
	"silver-ink/data"
	"silver-ink/engine"
)

func ReadBoardFile(name data.Name) (*data.Board, error) {
	board, err := SaveFile.ReadBoard(name)
	if err != nil {
		return nil, err
	}
	if board == nil {
		board, err = engine.WorldPath.LoadBoard(name)
		if err != nil {
			return nil, err
		}
	}

	return board, nil
}

var queuedBoard *data.Board

func QueueSwapBoard(board *data.Board, enter data.Vector) {
	if board.Player != -1 && enter.X != -1 && enter.Y != -1 {
		board.Objects[board.Player].Position = enter
	}
	queuedBoard = board
}

func swapBoard() error {
	if queuedBoard == nil {
		return nil
	}
	board := queuedBoard
	queuedBoard = nil
	if err := SaveBoard(engine.Board); err != nil {
		return err
	}
	SetBoard(board)
	return nil
}

func SetBoard(board *data.Board) error {
	engine.ClearActionBuffer()
	ResetText()
	SetView(board.ViewID)
	UpdateCamera(board.Follow, board)
	engine.SetBoard(board)
	if err := engine.LoadScriptsToBoard(board); err != nil {
		return err
	}
	Event(ENTER, 0)
	return nil
}

func (g *Game) processBoard() {
	if engine.Board == nil {
		return
	}
	engine.Board.Tick++
	g.processAction()
	if err := swapBoard(); err != nil {
		log.Println(err)
	}
	processObjects()
}

const camera_width = 29
const camera_height = 25

var camera_size = data.Vector{X: camera_width, Y: camera_height}

func (g *Game) processAction() {
	if engine.Board.Tick%SaveFile.InputCycle == 0 {
		if action := engine.GetAction(); action != data.DefaultName {
			switch action {
			case HK_PLAY:
				g.playAction()
			case HK_QUIT:
				g.quitAction()
			case HK_LOAD:
				OpenLoadPrompt()
			case HK_SAVE:
				OpenSavePrompt()
			default:
				if engine.Board.Player != -1 {
					ObjectEvent(engine.Board.Player, action, 0)
				}
			}
		}
	}
}

func (g *Game) isTitleBoard() bool {
	return engine.WorldFile.TitleBoard == engine.Board.Name
}

func (g *Game) quitAction() {
	engine.Quit()
}

func (g *Game) playAction() {
	if g.isTitleBoard() {
		board, err := ReadBoardFile(engine.WorldFile.StartBoard)
		if err != nil {
			log.Println(err)
			return
		}
		err = SetBoard(board)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

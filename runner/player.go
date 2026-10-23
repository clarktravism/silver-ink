package main

import (
	"silver-ink/data"
	"silver-ink/engine"
)

// return true if player crossed the border
func (ctx *ScriptContext) PlayerBorderCheck(pos data.Vector) bool {
	if !ctx.IsPlayer() {
		return false
	}

	name := data.Name{}
	enter := data.Vector{}

	if pos.Y < 0 {
		// north
		if engine.Board.North.Empty() {
			return false
		}
		name = engine.Board.North
		enter = data.Vector{X: pos.X, Y: data.BoardHeight - 1}
	} else if pos.Y >= data.BoardHeight {
		// south
		if engine.Board.South.Empty() {
			return false
		}
		name = engine.Board.South
		enter = data.Vector{X: pos.X, Y: 0}
	} else if pos.X < 0 {
		// west
		if engine.Board.West.Empty() {
			return false
		}
		name = engine.Board.West
		enter = data.Vector{X: data.BoardWidth - 1, Y: pos.Y}
	} else if pos.X >= data.BoardWidth {
		// east
		if engine.Board.East.Empty() {
			return false
		}
		name = engine.Board.East
		enter = data.Vector{X: 0, Y: pos.Y}
	}

	if !name.Empty() {
		board, err := ReadBoardFile(name)
		if err != nil {
			ctx.Error(err)
			return false
		}
		QueueSwapBoard(board, enter)
		return true
	}
	return false
}

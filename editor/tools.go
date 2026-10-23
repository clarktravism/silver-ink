package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	//"silver-ink/render"
	"strconv"
)

type Tool interface {
	SetBrush(value Brush)
	Apply(boardPos data.Vector, state engine.MouseState)
	Init()
	Close()
}

// Draw Tool
type DrawTool struct {
	Brush       Brush
	ScrollBoard *ScrollBoard
	lineStart   *data.Vector
}

func (dt *DrawTool) Init() {

}

func (dt *DrawTool) Close() {
	dt.lineStart = nil
	dt.ScrollBoard.ClearOverlay()
}

func (dt *DrawTool) SetBrush(value Brush) {
	dt.Brush = value
}

func (dt *DrawTool) Apply(boardPos data.Vector, state engine.MouseState) {
	if boardPos.X == -1 || boardPos.Y == -1 {
		dt.lineStart = nil
		dt.ScrollBoard.ClearOverlay()
		return
	}
	if state.LeftButtonDown {
		// line draw
		if engine.IsShift(state.Modifier) {
			if dt.lineStart == nil {
				dt.lineStart = &boardPos
			}
			dt.ScrollBoard.ClearOverlay()
			data.Line(*dt.lineStart, boardPos, func(v data.Vector) {
				dt.ScrollBoard.SetOverlay(v, dt.Brush.Cell)
			})
			return
		}
		// regular draw
		if dt.lineStart == nil {
			dt.Brush.Apply(engine.Board, boardPos)
		} else {
			dt.lineStart = nil
			dt.ScrollBoard.ClearOverlay()
		}
	} else if dt.lineStart != nil { // mouse released
		data.Line(*dt.lineStart, boardPos, func(v data.Vector) {
			dt.Brush.Apply(engine.Board, v)
		})
		dt.lineStart = nil
		dt.ScrollBoard.ClearOverlay()
	} else { // regular hover
		dt.ScrollBoard.ClearOverlay()
		dt.ScrollBoard.SetOverlay(boardPos, dt.Brush.Cell)
	}
}

// Remove Obj Tool
type RemoveObjTool struct{
	ScrollBoard *ScrollBoard
}

func (rot *RemoveObjTool) Init()  {}
func (rot *RemoveObjTool) Close() {
	rot.ScrollBoard.ClearOverlay()
}

func (rot *RemoveObjTool) SetBrush(Brush) {
}

var removeCell = data.Cell{ Char: 'X', Color: data.ColorDKRED }

func (rot *RemoveObjTool) Apply(boardPos data.Vector, state engine.MouseState) {
	if state.LeftButtonDown {
		rot.ScrollBoard.ClearOverlay()
		objects := engine.Board.Objects
		for i := range objects {
			if objects[i].Position == boardPos {
				engine.Board.RemoveObject(i)
				return
			}
		}
	} else {
		rot.ScrollBoard.ClearOverlay()
		cell := data.Cell{ Char: 'X', Color: data.ColorGRAY }
		objects := engine.Board.Objects
		for i := range objects {
			if objects[i].Position == boardPos {
				cell.Color = data.ColorRED
				break
			}
		}
		rot.ScrollBoard.SetOverlay(boardPos, cell)
	}
}

// Text Tool
type TextTool struct {
	Brush       Brush
	ScrollBoard *ScrollBoard
	textInput   engine.TextInput
	startPos    data.Vector
	cursorPos   data.Vector
}

func (tt *TextTool) Init() {
	tt.textInput.OnKeyPress = tt.onKeyPress
	tt.textInput.OnAccept = tt.onAccept
	tt.textInput.OnCancel = tt.Close
	tt.textInput.OnBackspace = tt.onBackspace
}

func (tt *TextTool) Close() {
	tt.ScrollBoard.ClearOverlay()
	tt.textInput.LoseFocus()
}

func (tt *TextTool) SetBrush(brush Brush) {
	tt.Brush = brush
}

func (tt *TextTool) Apply(boardPos data.Vector, state engine.MouseState) {
	if state.LeftButtonClick {
		tt.textInput.TakeFocus()
		tt.startPos = boardPos
		tt.setCursor(boardPos)
	}
}

func (tt *TextTool) setCursor(pos data.Vector) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= data.BoardWidth || pos.Y >= data.BoardHeight {
		tt.Close()
		return
	}
	tt.cursorPos = pos
	tt.ScrollBoard.ClearOverlay()
	tt.ScrollBoard.SetOverlay(tt.cursorPos, data.Cell{Char: '_', Color: data.ColorGRAY})
}

func (tt *TextTool) onKeyPress(text string) {
	if len(text) == 0 {
		return
	}
	tt.Brush.Cell.Char = text[0]
	tt.Brush.Apply(engine.Board, tt.cursorPos)
	tt.setCursor(tt.cursorPos.Add(data.Vector{X: 1, Y: 0}))
}

func (tt *TextTool) onAccept() {
	pos := tt.startPos.Add(data.Vector{X: 0, Y: 1})
	tt.startPos = pos
	tt.setCursor(pos)
}

func (tt *TextTool) onBackspace() {
	pos := tt.cursorPos.Add(data.Vector{X: -1, Y: 0})
	tt.setCursor(pos)
	if tt.textInput.Focused {
		tt.Brush.Cell.Char = 32
		tt.Brush.Apply(engine.Board, tt.cursorPos)
	}
}

// Fill Tool
type FillTool struct{
	Brush Brush
	ScrollBoard *ScrollBoard
}

func (ft *FillTool) Init()  {
	ft.ScrollBoard.ClearOverlay()
}
func (ft *FillTool) Close() {}

func (ft *FillTool) SetBrush(value Brush) {
	ft.Brush = value
}

var fillCell = data.Cell{ Char: 'F', Color: data.ColorWHITE }

func (ft *FillTool) Apply(boardPos data.Vector, state engine.MouseState) {
	if ft.Brush.Prefab != nil {
		return
	}
	if state.LeftButtonClick {
		//capture tile at pos
		tile := engine.Board.Tiles[data.TileID(boardPos)]
		newTile := data.Tile{ID: ft.Brush.ElementID, Cell: ft.Brush.Cell}
		if tile == newTile {
			return
		}

		data.Fill(boardPos,
			// inside
			func(v data.Vector) bool {
				if data.OutsideBoard(v) { return false }
				return engine.Board.Tiles[data.TileID(v)] == tile
			},
			// set
			func (v data.Vector) {
				engine.Board.Tiles[data.TileID(v)] = newTile
			})
	} else {
		ft.ScrollBoard.ClearOverlay()
		ft.ScrollBoard.SetOverlay(boardPos, fillCell)
	}
}

// Brush
type Brush struct {
	ElementID byte
	Cell      data.Cell
	Prefab    *data.Prefab
}

func (b *Brush) Type() string {
	if b.Prefab != nil {
		return "Object"
	}
	return "Element " + strconv.Itoa(int(b.ElementID))
}

func (b *Brush) Name() data.Name {
	if b.Prefab != nil {
		return b.Prefab.Name
	}
	if element := engine.ElementsFile.GetElement(b.ElementID); element != nil {
		return element.Name
	}
	return data.DefaultName
}

func (b *Brush) Apply(board *data.Board, pos data.Vector) {
	engine.Dirty()
	if b.Prefab != nil {
		// remove any object at position
		objects := engine.Board.Objects
		for i := range objects {
			if objects[i].Position == pos {
				engine.Board.RemoveObject(i)
				break
			}
		}

		obj := b.Prefab.Object()
		obj.Position = pos
		obj.Cell = b.Cell
		board.AddObject(obj, b.Prefab.Player, b.Prefab.Follow)
		return
	}
	tid := data.TileID(pos)
	board.Tiles[tid] = data.Tile{
		ID:   b.ElementID,
		Cell: b.Cell,
	}
}

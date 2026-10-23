package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

type ScrollBoard struct {
	Rectangle     data.Rect
	BoardViewMode BoardViewMode

	boardHover   data.Vector
	scrollOffset data.Vector

	//signals
	BoardMouseChange engine.Signal[BoardMouseArgs]

	hover   bool
	overlay []OverlayCell
}

type OverlayCell struct {
	Position data.Vector
	Cell     data.Cell
}

// signal args
type BoardMouseArgs struct {
	engine.MouseState
	BoardPosition data.Vector
}

type BoardViewMode int

const (
	BoardViewModeAll     BoardViewMode = 0
	BoardViewModeTiles   BoardViewMode = 1
	BoardViewModeObjects BoardViewMode = 2
)

func (sb *ScrollBoard) SetViewMode(mode BoardViewMode) {
	sb.BoardViewMode = mode
	engine.Dirty()
}

func (sb *ScrollBoard) Move(vector data.Vector) {
	dst := vector.Add(sb.scrollOffset)

	maxX := data.BoardWidth - sb.Rectangle.Size.X
	maxY := data.BoardHeight - sb.Rectangle.Size.Y

	if dst.X < 0 {
		dst.X = 0
	} else if dst.X > maxX {
		dst.X = maxX
	}

	if dst.Y < 0 {
		dst.Y = 0
	} else if dst.Y > maxY {
		dst.Y = maxY
	}

	sb.scrollOffset = dst
	engine.Dirty()
}

// control override
func (sb *ScrollBoard) Init() {
	sb.boardHover = data.Vector{X: -1, Y: -1}
	sb.overlay = make([]OverlayCell, 0, 32)
}

func (sb *ScrollBoard) Rect() data.Rect {
	return sb.Rectangle
}

func (sb *ScrollBoard) ProcessMouse(ms engine.MouseState) {
	inside := sb.Rectangle.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(sb.Rectangle.Position)
	if inside != sb.hover {
		sb.hover = inside
		//hover stuff
		engine.Dirty()
		if !sb.hover {
			sb.boardHover = data.Vector{X: -1, Y: -1}
		}
	}
	bma := BoardMouseArgs{
		MouseState:    ms,
		BoardPosition: sb.BoardPos(ms.Position),
	}
	sb.BoardMouseChange.Emit(bma)
}

func (sb *ScrollBoard) IsInside(pos data.Vector) bool {
	return sb.Rectangle.IsInside(pos)
}

func (sb *ScrollBoard) Render(buffer *render.Buffer, offset data.Vector) {
	if engine.Board == nil {
		return
	}

	offset = offset.Add(sb.Rectangle.Position)

	// draw tiles
	if sb.BoardViewMode != BoardViewModeObjects {
		tiles := engine.Board.Tiles
		tidy := data.TileID(sb.scrollOffset)
		biy := buffer.Index(offset.X, offset.Y)
		for iy := 0; iy < sb.Rectangle.Size.Y; iy++ {
			// reset x
			tidx := tidy
			bix := biy
			for ix := 0; ix < sb.Rectangle.Size.X; ix++ {
				buffer.SetCell(bix, tiles[tidx].Cell)
				tidx++
				bix++
			}
			// shift y + 1
			tidy += data.BoardWidth
			biy += buffer.Stride()
		}
	} else { // clear tiles
		biy := buffer.Index(offset.X, offset.Y)
		for iy := 0; iy < sb.Rectangle.Size.Y; iy++ {
			// reset x
			bix := biy
			for ix := 0; ix < sb.Rectangle.Size.X; ix++ {
				buffer.SetCell(bix, data.Cell{})
				bix++
			}
			// shift y + 1
			biy += buffer.Stride()
		}
	}

	if sb.BoardViewMode != BoardViewModeTiles {
		for _, o := range engine.Board.Objects {
			p := o.Position.Sub(sb.scrollOffset)
			if sb.Rectangle.IsInside(p) {
				bi := buffer.Index(p.X, p.Y)
				buffer.SetCell(bi, o.Cell)
			}
		}
	}

	for _, o := range sb.overlay {
		p := o.Position.Sub(sb.scrollOffset)
		bi := buffer.Index(p.X, p.Y)
		buffer.SetCell(bi, o.Cell)
	}

}

// adjusts the screen pos with the scroll to get board position
func (sb *ScrollBoard) BoardPos(pos data.Vector) data.Vector {
	if sb.Rectangle.IsInside(pos) {
		pos = pos.Add(sb.scrollOffset)
		if pos.X < 0 || pos.Y < 0 || pos.X >= data.BoardWidth || pos.Y >= data.BoardHeight {
			return data.Vector{X: -1, Y: -1}
		}
		return pos
	}
	return data.Vector{X: -1, Y: -1}
}

func (sb *ScrollBoard) ClearOverlay() {
	sb.overlay = sb.overlay[:0]
}

func (sb *ScrollBoard) SetOverlay(pos data.Vector, cell data.Cell) {
	sb.overlay = append(sb.overlay, OverlayCell{
		Position: pos,
		Cell:     cell,
	})
}

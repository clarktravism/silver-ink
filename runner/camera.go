package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

var bufferDrawFlags [buffer_length]data.BitField

func UpdateCamera(i int, board *data.Board) {
	if board == nil || board.Follow == -1 || board.Follow != i {
		return
	}
	obj := board.Objects[i]
	board.Camera = obj.Position.Add(view.Camera.Offset)
	if board.Camera.X < 0 {
		board.Camera.X = 0
	}
	if board.Camera.Y < 0 {
		board.Camera.Y = 0
	}
	if board.Camera.X + view.Camera.Size.X > data.BoardWidth {
		board.Camera.X = data.BoardWidth - view.Camera.Size.X
	}
	if board.Camera.Y + view.Camera.Size.Y > data.BoardHeight {
		board.Camera.Y = data.BoardHeight - view.Camera.Size.Y
	}
}

func CameraRender(buffer *render.Buffer) {
	camRect := data.Rect{
		Position: engine.Board.Camera,
		Size:     view.Camera.Size,
	}

	cells := buffer.GetCells()

	// draw tiles to buffer
	tiles := engine.Board.Tiles
	tY := data.TileID(camRect.Position)
	bY := buffer.VectorIndex(view.Camera.Position)
	for iy := 0; iy < camRect.Size.Y; iy++ {
		// reset x
		tX := tY
		bX := bY
		for ix := 0; ix < camRect.Size.X; ix++ {
			TileEvent(tX, TICK, 0)
			cells[bX] = tiles[tX].Cell
			tX++
			bX++
		}
		// move y
		tY += data.BoardWidth
		bY += buffer.Stride()
	}

	// draw objects to buffer
	bufferDrawFlags = [buffer_length]data.BitField{}
	for _, o := range engine.Board.Objects {
		if o.Free || !camRect.IsInside(o.Position) {
			continue
		}
		p := o.Position.Add(view.Camera.Position).Sub(camRect.Position)
		bi := buffer.Index(p.X, p.Y)
		if !bufferDrawFlags[bi].Has(data.DrawFlagPriority) {
			
			if o.DrawFlags.Has(data.DrawFlagShowBG) {
				cells[bi] = o.Cell
			} else {
				cells[bi].Char = o.Cell.Char
				cells[bi].Color = (cells[bi].Color & 0xF0) | (o.Cell.Color & 0x0F)
			}
			
			bufferDrawFlags[bi] = o.DrawFlags
		}
	}
}

package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

type BufferTool struct {
	Position data.Vector

	Items [32]Brush

	// signals
	ItemClicked engine.Signal[BufferItemArgs]
	ItemHover   engine.Signal[BufferItemArgs]

	rect data.Rect
}

// signal args
type BufferItemArgs struct {
	RightButtonClick bool
	LeftButtonClick  bool
	Brush            *Brush
}

// control overrides

func (bt *BufferTool) Init() {
	bt.rect = data.Rect{
		Position: bt.Position,
		Size:     data.Vector{X: 16, Y: 2},
	}
}

func (bt *BufferTool) Rect() data.Rect {
	return bt.rect
}

func (bt *BufferTool) ProcessMouse(ms engine.MouseState) {
	inside := bt.rect.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(bt.rect.Position)
	if inside {
		i := (ms.Position.Y * 16) + ms.Position.X
		bt.ItemHover.Emit(BufferItemArgs{
			RightButtonClick: false,
			LeftButtonClick:  false,
			Brush:            &bt.Items[i],
		})
		if ms.LeftButtonClick != ms.RightButtonClick {
			bt.ItemClicked.Emit(BufferItemArgs{
				RightButtonClick: ms.RightButtonClick,
				LeftButtonClick:  ms.LeftButtonClick,
				Brush:            &bt.Items[i],
			})
		}
	}
}

func (bt *BufferTool) Render(buffer *render.Buffer, offset data.Vector) {
	pos := bt.Position.Add(offset)

	sY := buffer.Index(pos.X, pos.Y)
	i := 0
	for iy := 0; iy < 2; iy++ {
		sX := sY
		for ix := 0; ix < 16; ix++ {
			buffer.SetCell(sX, bt.Items[i].Cell)
			sX++
			i++
		}
		sY += buffer.Stride()
	}

}

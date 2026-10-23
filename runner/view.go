package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
	"strconv"
)

var view data.View

func SetView(id int) {
	view = engine.ViewFile.GetView(id)
	if engine.Board != nil {
		UpdateCamera(engine.Board.Follow, engine.Board)
	}
}

func ViewRender(buffer *render.Buffer) {
	for _, p := range view.Panels {
		bY := buffer.VectorIndex(p.Position)
		for iy := 0; iy < p.Size.Y; iy++ {
			bX := bY
			for ix := 0; ix < p.Size.X; ix++ {
				buffer.Cells[bX] = p.Cell
				bX++
			}
			bY += buffer.Stride()
		}
	}

	for _, pi := range view.Icons {
		buffer.SetCell(buffer.VectorIndex(pi.Position), pi.Cell)
	}

	for _, pl := range view.Labels {
		buffer.WriteText(buffer.VectorIndex(pl.Position), pl.Color, pl.Text)
	}

	for _, pf := range view.Flags {
		counter := SaveFile.GetCounter(pf.Counter)
		if counter == nil || counter.Value <= 0 {
			continue
		}
		buffer.SetCell(buffer.VectorIndex(pf.Position), pf.Cell)
	}

	for _, pc := range view.Counters {
		str := "0"
		counter := SaveFile.GetCounter(pc.Counter)
		if counter != nil {
			str = strconv.Itoa(counter.Value)
		}
		buffer.WriteText(buffer.VectorIndex(pc.Position), pc.Color, str)
	}
}

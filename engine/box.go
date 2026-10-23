package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

// component to help controls render frames/boxes
type Box struct {
	Rectangle   data.Rect
	Style       data.BoxStyle
	Title       string
	dirty       bool
}

func NewBox(style data.BoxStyle, title string) Box {
	box := Box{
		Style: style,
		Title: title,
	}
	return box
}

func (box *Box) InsidePostion() data.Vector {
	pos := data.Vector{
		X: box.Rectangle.Position.X + box.Style.PaddingLeft + 1,
		Y: box.Rectangle.Position.Y + box.Style.PaddingTop + 1,
	}
	if box.Title != "" {
		pos.Y += 2
	}
	return pos
}

func (box *Box) InsideSize() data.Vector {
	size := data.Vector{
		X: box.Rectangle.Size.X - box.Style.PaddingLeft - box.Style.PaddingRight - 2,
		Y: box.Rectangle.Size.Y - box.Style.PaddingTop - box.Style.PaddingBottom - 2,
	}
	if box.Title != "" {
		size.Y -= 2
	}
	return size
}

func (box *Box) Render(buffer *render.Buffer, offset data.Vector) {
	rect := box.Rectangle
	rect.Position = rect.Position.Add(offset)
	p2 := rect.Position2()

	// top border
	index := buffer.Index(rect.Position.X, rect.Position.Y)
	buffer.SetCell(index, box.Style.BorderTopLeft)
	for i := 0; i < rect.Size.X-2; i++ {
		index++
		buffer.SetCell(index, box.Style.BorderTop)
	}
	buffer.SetCell(index+1, box.Style.BorderTopRight)

	//bottom border
	index = buffer.Index(rect.Position.X, p2.Y)
	buffer.SetCell(index, box.Style.BorderBottomLeft)
	for i := 0; i < rect.Size.X-2; i++ {
		index++
		buffer.SetCell(index, box.Style.BorderBottom)
	}
	buffer.SetCell(index+1, box.Style.BorderBottomRight)

	//left border
	index = buffer.Index(rect.Position.X, rect.Position.Y+1)
	for i := 0; i < rect.Size.Y-2; i++ {
		buffer.SetCell(index, box.Style.BorderLeft)
		index += buffer.Stride()
	}

	//right border
	index = buffer.Index(p2.X, rect.Position.Y+1)
	for i := 0; i < rect.Size.Y-2; i++ {
		buffer.SetCell(index, box.Style.BorderRight)
		index += buffer.Stride()
	}

	// Title
	if box.Title != "" {
		// Text
		index = buffer.Index(rect.Position.X+1, rect.Position.Y+1)
		index = buffer.WriteText(index, box.Style.TitleBackground.Color, box.Title)
		for i := len(box.Title); i < rect.Size.X-2; i++ {
			buffer.SetCell(index, box.Style.TitleBackground)
			index++
		}

		// Horizontal Divide
		index = buffer.Index(rect.Position.X, rect.Position.Y+2)
		buffer.SetCell(index, box.Style.HorizontalDivideLeft)
		index++
		for i := 0; i < rect.Size.X-2; i++ {
			buffer.SetCell(index, box.Style.HorizontalDivide)
			index++
		}
		buffer.SetCell(index, box.Style.HorizontalDivideRight)

		// Background
		indexY := buffer.Index(rect.Position.X+1, rect.Position.Y+3)
		for iy := 0; iy < rect.Size.Y-4; iy++ {
			index = indexY
			for ix := 0; ix < rect.Size.X-2; ix++ {
				buffer.SetCell(index, box.Style.Background)
				index++
			}
			indexY += buffer.Stride()
		}

	} else {
		// Background - Without Title
		indexY := buffer.Index(rect.Position.X+1, rect.Position.Y+1)
		for iy := 0; iy < rect.Size.Y-2; iy++ {
			index = indexY
			for ix := 0; ix < rect.Size.X-2; ix++ {
				buffer.SetCell(index, box.Style.Background)
				index++
			}
			indexY += buffer.Stride()
		}
	}
}
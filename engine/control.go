package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type Control interface {
	Init()
	Rect() data.Rect
	Render(buffer *render.Buffer, offset data.Vector)
	ProcessMouse(ms MouseState)
}

func CenterPos(size data.Vector) data.Vector {
	return data.Vector{
		X: (BufferWidth - size.X) / 2,
		Y: (BufferHeight - size.Y) / 2,
	}
}
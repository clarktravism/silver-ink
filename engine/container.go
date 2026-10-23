package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

type Container struct {
	Rectangle  data.Rect
	Children   []Control
	hoverChild int
}

func (c *Container) Init() {
	c.hoverChild = -1
}

func (c *Container) AddChild(child Control) {
	c.Children = append(c.Children, child)
	child.Init()
}

func (c *Container) Render(buffer *render.Buffer, offset data.Vector) {
	offset = offset.Add(c.Rectangle.Position)
	for _, child := range c.Children {
		child.Render(buffer, offset)
	}
}

func (c *Container) ProcessMouse(ms MouseState) {
	inside := c.Rectangle.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(c.Rectangle.Position)

	// check if mouse is within the same child since last mouse state
	if c.hoverChild != -1 {
		child := c.Children[c.hoverChild]
		child.ProcessMouse(ms)
		if child.Rect().IsInside(ms.Position) {
			return // dont need to check other children
		}
		c.hoverChild = -1
	}

	if inside {
		for i, child := range c.Children {
			if child.Rect().IsInside(ms.Position) {
				child.ProcessMouse(ms)
				c.hoverChild = i
				return
			}
		}
	}
}

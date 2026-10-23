package data

var RectZero = Rect{Position: VectorZero, Size: VectorZero}

type Rect struct {
	Position Vector `json:"Position"`
	Size Vector `json:"Size"`
}

func (r Rect) Position2() Vector {
	return Vector{X: r.Position.X + r.Size.X - 1, Y: r.Position.Y + r.Size.Y - 1}
}

func NewRect(x int, y int, width int, height int) Rect {
	return Rect{
		Position: Vector{X: x, Y: y},
		Size: Vector{X: width, Y: height},
	}
}

func (r Rect) IsInside(pos Vector) bool {
	pos2 := r.Position2()
	return pos.X >= r.Position.X && pos.X <= pos2.X && pos.Y >= r.Position.Y && pos.Y <= pos2.Y
}

func (r Rect) Intersect(value Rect) Rect {
	pos2 := r.Position2().Min(value.Position2())
	pos := r.Position.Max(value.Position)
	return Rect{Position: pos, Size: pos2.Sub(pos).Add(Vector{X: 1, Y: 1})}
}
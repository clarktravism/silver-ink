package data

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
)

var VectorZero = Vector{}
var VectorMinus = Vector{X: -1, Y: -1}

type Vector struct {
	X, Y int
}

func (v Vector) Add(value Vector) Vector {
	return Vector{X: v.X + value.X, Y: v.Y + value.Y}
}

func (v Vector) Sub(value Vector) Vector {
	return Vector{X: v.X - value.X, Y: v.Y - value.Y}
}

func (v Vector) Min(value Vector) Vector {
	if v.X > value.X {
		v.X = value.X
	}
	if v.Y > value.Y {
		v.Y = value.Y
	}
	return v
}

func (v Vector) Max(value Vector) Vector {
	if value.X > v.X {
		v.X = value.X
	}
	if value.Y > v.Y {
		v.Y = value.Y
	}
	return v
}

var ErrVectorBadFormat = errors.New("vector bad format")

func (v *Vector) UnmarshalJSON(data []byte) error {
	str := ""
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	xStr, yStr, ok := strings.Cut(str, ",")
	if !ok {
		return ErrVectorBadFormat
	}
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return ErrVectorBadFormat
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return ErrVectorBadFormat
	}
	*v = Vector{X: x, Y: y}
	return nil
}

func (v Vector) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(v.X) + "," + strconv.Itoa(v.Y)), nil
}

func ReadVector(r io.Reader, vector *Vector) error {
	buf := [4]byte{}
	_, err := io.ReadFull(r, buf[:])
	vector.X = int(binary.LittleEndian.Uint16(buf[0:2]))
	vector.Y = int(binary.LittleEndian.Uint16(buf[2:4]))
	return err
}

func WriteVector(w io.Writer, vector Vector) error {
	if err := WriteInt(w, vector.X); err != nil {
		return err
	}
	if err := WriteInt(w, vector.Y); err != nil {
		return err
	}
	return nil
}

// Math Stuff

func Line(start Vector, end Vector, f func(Vector)) {
	if Abs(end.Y-start.Y) < Abs(end.X-start.X) {
		if start.X > end.X {
			lineLow(end, start, f)
		} else {
			lineLow(start, end, f)
		}
	} else {
		if start.Y > end.Y {
			lineHigh(end, start, f)
		} else {
			lineHigh(start, end, f)
		}
	}
}

func lineLow(start, end Vector, f func(Vector)) {
	d := end.Sub(start)
	yi := 1
	if d.Y < 0 {
		yi = -1
		d.Y = -d.Y
	}
	D := (2 * d.Y) - d.X
	y := start.Y
	for x := start.X; x <= end.X; x++ {
		f(Vector{X: x, Y: y})
		if D > 0 {
			y = y + yi
			D = D + (2 * (d.Y - d.X))
		} else {
			D = D + 2*d.Y
		}
	}
}

func lineHigh(start, end Vector, f func(Vector)) {
	d := end.Sub(start)
	xi := 1
	if d.X < 0 {
		xi = -1
		d.X = -d.X
	}
	D := (2 * d.X) - d.Y
	x := start.X
	for y := start.Y; y <= end.Y; y++ {
		f(Vector{X: x, Y: y})
		if D > 0 {
			x = x + xi
			D = D + (2 * (d.X - d.Y))
		} else {
			D = D + 2*d.X
		}
	}
}

func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func Fill(at Vector, inside func(Vector) bool, set func(Vector)) {
	if !inside(at) {
		return
	}
	a := make([]fillseed, 0, 8)
	a = append(a, fillseed{at.X, at.X, at.Y, 1})
	a = append(a, fillseed{at.X, at.X, at.Y - 1, -1})
	for len(a) > 0 {
		s := a[len(a)-1]
		a = a[:len(a)-1]
		x := s.x1
		if inside(Vector{X: x, Y: s.y}) {
			for inside(Vector{X: x - 1, Y: s.y}) {
				set(Vector{X: x - 1, Y: s.y})
				x = x - 1
			}
			if x < s.x1 {
				a = append(a, fillseed{x1: x, x2: s.x1 - 1, y: s.y - s.dy, dy: -s.dy})
			}
		}
		for s.x1 <= s.x2 {
			for inside(Vector{X: s.x1, Y: s.y}) {
				set(Vector{X: s.x1, Y: s.y})
				s.x1 = s.x1 + 1
			}
			if s.x1 > x {
				a = append(a, fillseed{x1: x, x2: s.x1 - 1, y: s.y + s.dy, dy: s.dy})
			}
			if s.x1-1 > s.x2 {
				a = append(a, fillseed{x1: s.x2 + 1, x2: s.x1 - 1, y: s.y - s.dy, dy: -s.dy})
			}
			s.x1 = s.x1 + 1
			for s.x1 <= s.x2 && !inside(Vector{X: s.x1, Y: s.y}) {
				s.x1 = s.x1 + 1
			}
			x = s.x1
		}
	}
}

type fillseed struct {
	x1 int
	x2 int
	y  int
	dy int
}

func (v1 Vector) Distance(v2 Vector) int {
	if v1.X == -1 || v1.Y == -1 || v2.X == -1 || v2.Y == -1 {
		return -1
	}
	return Abs(v1.X-v2.X) + Abs(v1.Y-v2.Y)
}

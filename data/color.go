package data

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

type Color byte

func (c Color) BG() Color {
	return (c & 0xF0) >> 4
}

func (c Color) FG() Color {
	return (c & 0x0F)
}

func NewColor(fg Color, bg Color) Color {
	return (bg << 4) | fg
}

func (c *Color) UnmarshalJSON(data []byte) error {
	// color name unmarshal
	if data[0] == '"' {
		str := ""
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}

		// hex unmarshal
		if len(str) > 0 && str[0] == 'x' {
			v, err := strconv.ParseInt(str[1:], 16, 8)
			if err != nil {
				return err
			}
			*c = Color(v)
			return nil
		}

		fg, bg, _ := strings.Cut(str, "/")
		cFG, _ := GetColor(NewName(fg))
		cBG, _ := GetColor(NewName(bg))
		*c = NewColor(cFG, cBG)
		return nil
	}

	// number unmarshal
	b := byte(0)
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	*c = Color(b)
	return nil
}

func (c Color) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(c))), nil
}

func ReadColor(r io.Reader, color *Color) error {
	buf := [1]byte{}
	_, err := io.ReadFull(r, buf[:])
	*color = Color(buf[0])
	return err
}
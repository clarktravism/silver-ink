package data

import (
	"io"
)

type Cell struct {
	Color Color `json:"Color"`
	Char  byte `json:"Char"`
}

func ReadCell(r io.Reader, cell *Cell) error {
	buf := [2]byte{}
	_, err := io.ReadFull(r, buf[:])
	*cell = Cell{Color: Color(buf[0]), Char: buf[1]}
	return err
}

func WriteCell(w io.Writer, cell Cell) error {
	if err := WriteByte(w, byte(cell.Color)); err != nil {
		return err
	}
	if err := WriteByte(w, byte(cell.Char)); err != nil {
		return err
	}
	return nil
}
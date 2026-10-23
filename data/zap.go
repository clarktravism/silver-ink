package data

import (
	"io"
)

type Zap struct {
	Name  Name
	Value uint8
}

func (o *Object) IncZap(name Name, value uint8) {
	zap := o.getZap(name)
	if zap != nil {
		zap.Value += value
		return
	}
	o.Zaps = append(o.Zaps, Zap{Name: name, Value: 1})
}

func (o *Object) DecZap(name Name, value uint8) {
	zap := o.getZap(name)
	if zap != nil {
		if value > zap.Value {
			zap.Value = 0
		} else {
			zap.Value -= value
		}
		return
	}
}

func (o *Object) ClearZap(name Name) {
	zap := o.getZap(name)
	if zap != nil {
		zap.Value = 0
	}
}

func (o *Object) ClearAllZaps() {
	o.Zaps = o.Zaps[:0]
}

func (o *Object) getZap(name Name) *Zap {
	for i := range o.Zaps {
		if o.Zaps[i].Name == name {
			return &o.Zaps[i]
		}
	}
	return nil
}

func (o *Object) GetZap(name Name) uint8 {
	if zap := o.getZap(name); zap != nil {
		return zap.Value
	}
	return 0
}

func readZaps(r io.Reader, zaps *[]Zap) error {
	length := uint8(0)
	if err := ReadUint8(r, &length); err != nil {
		return err
	}
	slice := make([]Zap, length)
	for i := uint8(0); i < length; i++ {
		if err := ReadName(r, &slice[i].Name); err != nil {
			return err
		}
		if err := ReadUint8(r, &slice[i].Value); err != nil {
			return err
		}
	}
	*zaps = slice
	return nil
}

func writeZaps(w io.Writer, zaps []Zap) error {
	length := uint8(len(zaps))
	if err := WriteUint8(w, length); err != nil {
		return err
	}
	for _, zap := range zaps {
		if err := WriteName(w, zap.Name); err != nil {
			return err
		}
		if err := WriteUint8(w, zap.Value); err != nil {
			return err
		}
	}
	return nil
}
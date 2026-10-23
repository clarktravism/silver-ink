package data

import (
	"io"
	"encoding/binary"
)

type Counters struct {
	Counters []Counter `json:"Counters"`
}

type Counter struct {
	Name  Name `json:"Name"`
	Value int  `json:"Value"`
}

func (c *Counters) ClearCounter(name Name) {
	for i := range c.Counters {
		if c.Counters[i].Name == name {
			c.Counters[i] = c.Counters[len(c.Counters)-1]
			c.Counters[len(c.Counters)-1] = Counter{}
			c.Counters = c.Counters[:len(c.Counters)-1]
			break
		}
	}
}

func (c *Counters) SetCounter(name Name, value int) {
	counter := c.GetCounter(name)
	if counter != nil {
		counter.Value = value
		return
	}
	c.Counters = append(c.Counters, Counter{Name: name, Value: value})
}

func (c *Counters) IncCounter(name Name, value int) {
	counter := c.GetCounter(name)
	if counter != nil {
		counter.Value += value
		return
	}
	c.Counters = append(c.Counters, Counter{Name: name, Value: value})
}

func (c *Counters) GetCounter(name Name) *Counter {
	for i := range c.Counters {
		if c.Counters[i].Name == name {
			return &c.Counters[i]
		}
	}
	return nil
}

func (c *Counters) GetOrAddCounter(name Name) *Counter {
	for i := range c.Counters {
		if c.Counters[i].Name == name {
			return &c.Counters[i]
		}
	}
	c.Counters = append(c.Counters, Counter{Name: name, Value: 0})
	return &c.Counters[len(c.Counters)-1]
}

func ReadCounters(r io.Reader, counters *Counters) error {
	length := uint8(0)
	if err := ReadUint8(r, &length); err != nil {
		return err
	}
	slice := make([]Counter, length)
	for i := uint8(0); i < length; i++ {
		if err := ReadName(r, &slice[i].Name); err != nil {
			return err
		}
		if err := ReadInt(r, &slice[i].Value); err != nil {
			return err
		}
	}
	counters.Counters = slice
	return nil
}

func WriteCounters(w io.Writer, counters Counters) error {
	c := counters.Counters
	binary.Write(w, binary.LittleEndian, uint8(len(c)))
	for i := range c {
		if err := WriteName(w, c[i].Name); err != nil {
			return err
		}
		if err := WriteInt(w, c[i].Value); err != nil {
			return err
		}
	}
	return nil
}

func (c *Counters) Copy() Counters {
	cc := Counters{}
	cc.Counters = make([]Counter, len(c.Counters))
	copy(cc.Counters, c.Counters)
	return cc
}
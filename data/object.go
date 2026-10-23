package data

import (
	"errors"
	"io"
)

type Object struct {
	Name     Name
	Group    Name
	Cycle    uint8
	Position Vector
	Step     Vector
	Cell
	DrawFlags      BitField
	Collision      BitField
	PushFlags      BitField
	ScriptName     Name
	ScriptPosition int
	Lock           byte
	Target         Name
	Counters
	Zaps   []Zap
	Repeat int

	// not persisted
	Script *Script
	Free   bool
}

const ( //DrawFlags
	DrawFlagPriority BitField = 0x1
	DrawFlagShowBG   BitField = 0x2
	DrawFlagLit      BitField = 0x4
)

const (
	CanPush       BitField = 0x1
	PushableNorth BitField = 0x2
	PushableSouth BitField = 0x4
	PushableEast  BitField = 0x8
	PushableWest  BitField = 0x10
)

const (
	PushableNorthSouth = PushableNorth | PushableSouth
	PushableEastWest   = PushableEast | PushableWest
	PushableAll        = PushableNorthSouth | PushableEastWest
)

const (
	Collision1 BitField = 0x1
	Collision2 BitField = 0x2
	Collision3 BitField = 0x4
	Collision4 BitField = 0x8
	Collision5 BitField = 0x10
	Collision6 BitField = 0x20
	Collision7 BitField = 0x40
	Collision8 BitField = 0x80
)

func (obj *Object) SetFree() {
	obj.Free = true
}

func (obj *Object) LabelPosition(label Name, skip int) int {
	if obj.Lock != 0 || obj.Script == nil {
		return -1
	}
	skip += int(obj.GetZap(label))
	return obj.Script.LabelPosition(label, skip)
}

const (
	// object byte codes
	objBasic    byte = 1 // name, group, cycle, position, step, cell
	objFlags    byte = 2 // collision, drawflags, pushflags
	objCounters byte = 3
	objScript   byte = 4 // script name, script position, lock, repeat, zaps
)

func readObject(r io.Reader, obj *Object) error {

	segments := uint8(0)
	if err := ReadUint8(r, &segments); err != nil {
		return ErrBadBoardData
	}

	for i := 0; i < int(segments); i++ {
		// read next code
		code := byte(0)
		if err := ReadByte(r, &code); err != nil {
			// finished reading
			if errors.Is(err, io.EOF) {
				break
			}
			return ErrBadBoardData
		}

		var err error

		switch code {
		case objBasic:
			err = objBasicRead(r, obj)
		case objFlags:
			err = objFlagsRead(r, obj)
		case objCounters:
			err = objCountersRead(r, obj)
		case objScript:
			err = objScriptRead(r, obj)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func objBasicRead(r io.Reader, obj *Object) error {
	if err := ReadName(r, &obj.Name); err != nil {
		return err
	}
	if err := ReadName(r, &obj.Group); err != nil {
		return err
	}
	if err := ReadByte(r, &obj.Cycle); err != nil {
		return err
	}
	if err := ReadVector(r, &obj.Position); err != nil {
		return err
	}
	if err := ReadVector(r, &obj.Step); err != nil {
		return err
	}
	if err := ReadCell(r, &obj.Cell); err != nil {
		return err
	}
	return nil
}

func objFlagsRead(r io.Reader, obj *Object) error {
	if err := ReadBitField(r, &obj.Collision); err != nil {
		return err
	}
	if err := ReadBitField(r, &obj.DrawFlags); err != nil {
		return err
	}
	if err := ReadBitField(r, &obj.PushFlags); err != nil {
		return err
	}
	return nil
}

func objCountersRead(r io.Reader, obj *Object) error {
	if err := ReadCounters(r, &obj.Counters); err != nil {
		return err
	}
	if err := ReadName(r, &obj.Target); err != nil {
		return err
	}
	return nil
}

func objScriptRead(r io.Reader, obj *Object) error {
	if err := ReadName(r, &obj.ScriptName); err != nil {
		return err
	}
	if err := ReadInt(r, &obj.ScriptPosition); err != nil {
		return err
	}
	if err := ReadByte(r, &obj.Lock); err != nil {
		return err
	}
	if err := ReadInt(r, &obj.Repeat); err != nil {
		return err
	}
	if err := readZaps(r, &obj.Zaps); err != nil {
		return err
	}
	return nil
}

func writeObject(w io.Writer, obj Object) error {

	WriteUint8(w, 4)

	WriteByte(w, objBasic)
	if err := WriteName(w, obj.Name); err != nil {
		return err
	}
	if err := WriteName(w, obj.Group); err != nil {
		return err
	}
	if err := WriteByte(w, obj.Cycle); err != nil {
		return err
	}
	if err := WriteVector(w, obj.Position); err != nil {
		return err
	}
	if err := WriteVector(w, obj.Step); err != nil {
		return err
	}
	if err := WriteCell(w, obj.Cell); err != nil {
		return err
	}

	WriteByte(w, objFlags)
	if err := WriteBitField(w, obj.Collision); err != nil {
		return err
	}
	if err := WriteBitField(w, obj.DrawFlags); err != nil {
		return err
	}
	if err := WriteBitField(w, obj.PushFlags); err != nil {
		return err
	}

	WriteByte(w, objCounters)
	if err := WriteCounters(w, obj.Counters); err != nil {
		return err
	}
	if err := WriteName(w, obj.Target); err != nil {
		return err
	}

	WriteByte(w, objScript)
	if err := WriteName(w, obj.ScriptName); err != nil {
		return err
	}
	if err := WriteInt(w, obj.ScriptPosition); err != nil {
		return err
	}
	if err := WriteByte(w, obj.Lock); err != nil {
		return err
	}
	if err := WriteInt(w, obj.Repeat); err != nil {
		return err
	}
	if err := writeZaps(w, obj.Zaps); err != nil {
		return err
	}
	return nil
}

func (object *Object) Prefab() Prefab {
	return Prefab{
		Name:       object.Name,
		Group:      object.Group,
		Cycle:      object.Cycle,
		Cell:       object.Cell,
		DrawFlags:  object.DrawFlags,
		Collision:  object.Collision,
		PushFlags:  object.PushFlags,
		Counters:   object.Counters.Copy(), // need deep copy
		ScriptName: object.ScriptName,
		Script:     object.Script,
		Target:     object.Target,
	}
}
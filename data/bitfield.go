package data

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type BitField byte

func (bf *BitField) Set(value BitField) {
	*bf |= value
}

func (bf *BitField) Clear(value BitField) {
	*bf &^= value
}

func (bf BitField) Has(value BitField) bool {
	return bf&value == value
}

func (bf BitField) Collide(value BitField) bool {
	return bf&value != 0
}

func (bf *BitField) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		str := ""
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}

		// number format "00000011"
		if val, err := strconv.ParseUint(str, 2, 8); err == nil {
			*bf = BitField(val)
			return nil
		}

		// text format
		val := BitField(0)
		strs := strings.Split(strings.ToUpper(str), ",")
		for _, s := range strs {
			if v, ok := bitfieldValues[s]; ok {
				val |= v
			}
		}
		*bf = val
		return nil
	}
	val := byte(0)
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*bf = BitField(val)
	return nil
}

func (bf *BitField) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%b", byte(*bf))), nil
}

func (bf BitField) String() string {
	return fmt.Sprintf("%b", byte(bf))
}

func ReadBitField(r io.Reader, value *BitField) error {
	buf := [1]byte{}
	_, err := io.ReadFull(r, buf[:])
	*value = BitField(buf[0])
	return err
}

func WriteBitField(w io.Writer, value BitField) error {
	binary.Write(w, binary.LittleEndian, byte(value))
	return nil
}

var bitfieldValues = map[string]BitField{
	// PUSH FLAGS
	"CAN_PUSH" : CanPush,
	"PUSHABLE_N" : PushableNorth,
	"PUSHABLE_S" : PushableSouth,
	"PUSHABLE_E" : PushableEast,
	"PUSHABLE_W" : PushableWest,
	"PUSHABLE_ALL" : PushableNorth | PushableSouth | PushableEast | PushableWest,
	"PUSHABLE_NS" : PushableNorth | PushableSouth,
	"PUSHABLE_EW" : PushableEast | PushableWest,

	// DRAW FLAGS
	"PRIORITY" : DrawFlagPriority,
	"SHOW_BG" : DrawFlagShowBG,
	"SEEN" : DrawFlagLit,

	// COLLISON
	"COLLISION_1" : Collision1,
	"COLLISION_2" : Collision2,
	"COLLISION_3" : Collision3,
	"COLLISION_4" : Collision4,
	"COLLISION_5" : Collision5,
	"COLLISION_6" : Collision6,
	"COLLISION_7" : Collision7,
	"COLLISION_8" : Collision8,

	"WALLS" : Collision1,
	"WATER" : Collision2,
	"FLOOR" : Collision3,
}
package data

import (
	"io"
	"strings"
	"unique"
	"encoding/json"
	"encoding/binary"
)

type Name unique.Handle[string]

var DefaultName Name

func NewName(value string) Name {
	return Name(unique.Make(strings.ToUpper(value)))
}

func (n Name) String() string {
	if n.Empty() { 
		return ""
	}
	return unique.Handle[string](n).Value()
}

func (n Name) Empty() bool {
	return n == DefaultName
}

func (n Name) Equal(value string) bool {
	return n.String() == value
}

func (n *Name) UnmarshalJSON(data []byte) error {
	str := ""
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*n = NewName(str)
	return nil
}

func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

func ReadName(r io.Reader, value *Name) error {
	length := uint8(0)
	if err := ReadUint8(r, &length); err != nil {
		return err
	}
	if length == 0 {
		*value = DefaultName
		return nil
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}
	*value = NewName(string(buf))
	return nil
}

func WriteName(w io.Writer, value Name) error {
	if value.Empty() {
		return binary.Write(w, binary.LittleEndian, uint8(0))
	}
	str := value.String()
	if err := binary.Write(w, binary.LittleEndian, uint8(len(str))); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, []byte(str))
}
package data

import (
	"encoding/binary"
	"io"
)

func ReadSlice[T any](r io.Reader) ([]T, error) {
	length := uint8(0)
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	slice := make([]T, length)
	for i := 0; i < len(slice); i++ {
		if err := binary.Read(r, binary.LittleEndian, &slice[i]); err != nil {
			return nil, err
		}
	}
	return slice, nil
}

func WriteSlice[T any](w io.Writer, slice []T) error {
	if err := binary.Write(w, binary.LittleEndian, uint8(len(slice))); err != nil {
		return err
	}
	for i := 0; i < len(slice); i++ {
		if err := binary.Write(w, binary.LittleEndian, slice[i]); err != nil {
			return nil
		}
	}
	return nil
}

func ReadUint8(r io.Reader, value *uint8) error {
	buf := [1]byte{}
	_, err := io.ReadFull(r, buf[:])
	*value = buf[0]
	return err
}

func ReadInt16(r io.Reader, value *int16) error {
	buf := [2]byte{}
	_, err := io.ReadFull(r, buf[:])
	*value = int16(binary.LittleEndian.Uint16(buf[:]))
	return err
}

func ReadInt(r io.Reader, value *int) error {
	buf := [2]byte{}
	_, err := io.ReadFull(r, buf[:])
	*value = int(int16(binary.LittleEndian.Uint16(buf[:])))
	return err
}

func ReadUInt16(r io.Reader, value *uint16) error {
	buf := [2]byte{}
	_, err := io.ReadFull(r, buf[:])
	*value = binary.LittleEndian.Uint16(buf[:])
	return err
}

func ReadByte(r io.Reader, value *byte) error {
	buf := [1]byte{}
	_, err := io.ReadFull(r, buf[:])
	*value = buf[0]
	return err
}

func ReadString(r io.Reader, value *string) error {
	length := uint8(0)
	if err := ReadUint8(r, &length); err != nil {
		return err
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}
	*value = string(buf)
	return nil
}

func WriteUint8(w io.Writer, value uint8) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func WriteInt16(w io.Writer, value int16) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func WriteUInt16(w io.Writer, value uint16) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func WriteInt(w io.Writer, value int) error {
	return binary.Write(w, binary.LittleEndian, int16(value))
}

func WriteByte(w io.Writer, value byte) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func WriteString(w io.Writer, value string) error {
	if value == "" {
		return binary.Write(w, binary.LittleEndian, uint8(0))
	}
	if err := binary.Write(w, binary.LittleEndian, uint8(len(value))); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, []byte(value))
}


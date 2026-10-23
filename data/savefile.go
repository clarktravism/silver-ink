package data

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

const saves_folder = "saves"

type SaveFile struct {
	CurrentBoard Name
	FPS          int
	InputCycle   uint8
	Counters
	Boards       map[Name][]byte
}

func (sf *SaveFile) SaveBoard(board *Board) error {
	if board.Flags&PersistFlag != PersistFlag {
		return nil
	}

	buf := bytes.Buffer{}
	buf.Grow(1024)

	if err := WriteBoard(&buf, board); err != nil {
		return err
	}

	if sf.Boards == nil {
		sf.Boards = make(map[Name][]byte)
	}

	sf.Boards[board.Name] = buf.Bytes()
	return nil
}

func (sf *SaveFile) ReadBoard(name Name) (*Board, error) {
	if sf == nil {
		return nil, nil
	}

	d, ok := sf.Boards[name]
	if !ok {
		return nil, nil
	}

	board := NewBoard()
	if err := ReadBoard(bytes.NewReader(d), &board); err != nil {
		return nil, err
	}

	if board.Flags&PersistFlag != PersistFlag {
		delete(sf.Boards, name)
	}

	return &board, nil
}

const (
	// world byte code
	saveCurrentBoard byte = 1
	saveFPS          byte = 2
	saveInputCycle   byte = 3
	saveCounters     byte = 5
	saveBoards       byte = 6
)

func (sf *SaveFile) CopyWorldData(wf *WorldFile) {
	sf.FPS = wf.FPS
	sf.InputCycle = wf.InputCycle
	sf.Counters = wf.Counters.Copy()
}

func (wp WorldPath) WriteSaveFile(path string, sf *SaveFile, current *Board) error {
	fullpath := wp.Join(saves_folder, path)
	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()

	sf.Write(f, current)
	return nil
}

func (sf *SaveFile) Write(w io.Writer, current *Board) error {

	// buffer to use
	buf := bytes.Buffer{}
	buf.Grow(1024)

	// number of segments
	WriteUint8(w, 5)

	// write world data
	WriteByte(w, saveCurrentBoard)
	WriteName(w, current.Name)

	WriteByte(w, saveFPS)
	WriteInt(w, sf.FPS)

	WriteByte(w, saveInputCycle)
	WriteUint8(w, sf.InputCycle)

	WriteByte(w, saveCounters)
	WriteCounters(w, sf.Counters)

	// write board data
	WriteByte(w, saveBoards)
	// write board count
	count := len(sf.Boards)
	if _, ok := sf.Boards[current.Name]; !ok {
		count++
	}
	WriteInt(w, count)

	// write current board - id, len, data
	WriteName(w, current.Name)
	if err := WriteBoard(&buf, current); err != nil {
		return err
	}
	WriteInt(w, buf.Len())
	buf.WriteTo(w)
	buf.Reset()

	// write all other boards
	for k, v := range sf.Boards {
		// skip current board
		if k == current.Name {
			continue
		}
		WriteName(w, k)    // write id
		WriteInt(w, len(v)) // len

		buf2 := bytes.NewBuffer(v)
		_, err := buf2.WriteTo(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wp WorldPath) ReadSaveFile(path string) (*SaveFile, error) {
	fullpath := wp.Join(saves_folder, path)
	f, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sf := SaveFile{}
	if err := sf.Read(f); err != nil {
		return nil, err
	}
	return &sf, nil
}

var ErrBadSaveData = errors.New("bad save data")

func (sf *SaveFile) Read(r io.Reader) error {
	br := bufio.NewReader(r)
	
	segments := uint8(0)
	if err := ReadUint8(br, &segments); err != nil {
		return ErrBadBoardData
	}

	for i := uint8(0); i < segments; i++ {
		// segment and length
		code := byte(0)
		if err := ReadByte(br, &code); err != nil {
			return ErrBadSaveData
		}

		switch code {
		case saveCurrentBoard:
			if err := ReadName(br, &sf.CurrentBoard); err != nil {
				return err
			}
		case saveFPS:
			if err := ReadInt(br, &sf.FPS); err != nil {
				return err
			}
		case saveInputCycle:
			if err := ReadUint8(br, &sf.InputCycle); err != nil {
				return err
			}
		case saveCounters:
			if err := ReadCounters(br, &sf.Counters); err != nil {
				return err
			}
		case saveBoards:
			len := 0
			if err := ReadInt(br, &len); err != nil {
				return ErrBadSaveData
			}

			sf.Boards = make(map[Name][]byte, len)
			for j := 0; j < len; j++ {
				// board id
				name := Name{}
				ReadName(br, &name)

				// board length
				blen := 0
				ReadInt(br, &blen)

				// board data
				b := make([]byte, blen)
				if _, err := io.ReadFull(br, b); err != nil {
					return ErrBadSaveData
				}

				sf.Boards[name] = b
			}
		default:
			return ErrBadSaveData
		}
	}
	return nil
}

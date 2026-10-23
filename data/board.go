package data

import (
	"bytes"
	"errors"
	"io"
	"iter"
	"os"
)

const BoardsFolder = "boards"
const BoardsEXT = ".sib"
const board_version = uint16(1)
const BoardWidth = 128
const BoardHeight = 96
const BoardLength = BoardWidth * BoardHeight

var BoardSize = Vector{X: BoardWidth, Y: BoardHeight}

const object_limit = 256

type Board struct {
	North   Name
	South   Name
	East    Name
	West    Name
	Charset Name
	Palette Name
	Flags   byte
	Counters

	// tiles
	Tiles [BoardLength]Tile

	// objects
	objects [object_limit]Object
	Objects []Object

	// camera
	Camera Vector
	Follow int
	ViewID int

	Player int

	// runner
	Tick     uint8
	ObjectID int
	Enter    Vector

	// not saved
	Name Name
}

func (b *Board) ObjectCount() int {
	return len(b.Objects)
}

type Tile struct {
	Cell
	ID uint8
}

const (
	SavedBoardFlag = 0x1
	PersistFlag    = 0x2
	DarknessFlag   = 0x4
)

func (b *Board) AddObject(obj Object, isPlayer bool, isFollow bool) int {
	if len(b.Objects) < object_limit {
		b.Objects = append(b.Objects, obj)
		i := len(b.Objects) - 1
		if isPlayer {
			b.Player = i
		}
		if isFollow {
			b.Follow = i
		}
		return i
	}
	return -1
}

func (b *Board) FreeObjects() {
	for i := 0; i < len(b.Objects); i++ {
		if b.Objects[i].Free {
			b.RemoveObject(i)
			i--
		}
	}
}

func (b *Board) RemoveObject(i int) {
	copy(b.objects[i:], b.Objects[i+1:])
	b.Objects[len(b.Objects)-1] = Object{}
	b.Objects = b.Objects[:len(b.Objects)-1]
	if b.ObjectID >= i && b.ObjectID != 0 {
		b.ObjectID--
	}
	// update id refs
	if b.Follow == i {
		b.Follow = -1
	} else if b.Follow > i {
		b.Follow--
	}
	if b.Player == i {
		b.Player = -1
	} else if b.Player > i {
		b.Player--
	}
}

func (b *Board) ObjectAt(pos Vector) int {
	for i := range b.Objects {
		if b.Objects[i].Position == pos && !b.Objects[i].Free {
			return i
		}
	}
	return -1
}

func (wp *WorldPath) LoadBoard(name Name) (*Board, error) {
	filename := name.String() + BoardsEXT
	fullpath := wp.Join(BoardsFolder, filename)

	board := NewBoard()
	board.Name = name

	f, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := ReadBoard(f, &board); err != nil {
		return nil, err
	}
	return &board, nil
}

var charsetDefault = NewName("CHR_DEFAULT")
var paletteDefault = NewName("PAL_DEFAULT")

func NewBoard() Board {
	return Board{
		Enter:   Vector{X: -1, Y: -1},
		Follow:  -1,
		Player:  -1,
		Charset: charsetDefault,
		Palette: paletteDefault,
		ViewID:  1,
		Flags:   PersistFlag,
	}
}

// temp should only be false for editor
func (wp *WorldPath) SaveBoard(board *Board) error {

	// get the fullpath, temp file might not exist yet
	fullpath := ""
	filename := board.Name.String() + BoardsEXT
	fullpath = wp.Join(BoardsFolder, filename)

	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()

	return WriteBoard(f, board)
}

var ErrBadBoardData = errors.New("bad board data")

const (
	// board byte codes
	brdConnections byte = 1
	brdCharset     byte = 2
	brdPalette     byte = 3
	brdFlags       byte = 4
	brdCounters    byte = 5
	brdTiles       byte = 6
	brdObjects     byte = 7 // count and object data
	brdCamera      byte = 8 // pos and follow
	brdPlayer      byte = 9
	brdRunState    byte = 10 // tick and objectid
	brdEnter       byte = 12
)

func ReadBoard(r io.Reader, board *Board) error {

	// segments
	segments := uint8(0)
	if err := ReadUint8(r, &segments); err != nil {
		return ErrBadBoardData
	}

	for i := uint8(0); i < segments; i++ {

		// read next code
		code := byte(0)
		if err := ReadByte(r, &code); err != nil {
			// finished reading
			if errors.Is(err, io.EOF) {
				break
			}
			return ErrBadBoardData
		}

		// read data by code
		var err error
		switch code {
		case brdConnections:
			err = readBrdConnections(r, board)
		case brdCharset:
			err = readBrdCharset(r, board)
		case brdPalette:
			err = readBrdPalette(r, board)
		case brdFlags:
			err = readBrdFlags(r, board)
		case brdCounters:
			err = readBrdCounters(r, board)
		case brdTiles:
			err = readBrdTiles(r, board)
		case brdObjects:
			err = readBrdObjects(r, board)
		case brdCamera:
			err = readBrdCamera(r, board)
		case brdPlayer:
			err = readBrdPlayer(r, board)
		case brdRunState:
			err = readBrdRunState(r, board)
		case brdEnter:
			err = readBrdEnter(r, board)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func readBrdConnections(r io.Reader, board *Board) error {
	// connected board ids
	if err := ReadName(r, &board.North); err != nil {
		return ErrBadBoardData
	}
	if err := ReadName(r, &board.South); err != nil {
		return ErrBadBoardData
	}
	if err := ReadName(r, &board.East); err != nil {
		return ErrBadBoardData
	}
	if err := ReadName(r, &board.West); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdCharset(r io.Reader, board *Board) error {
	if err := ReadName(r, &board.Charset); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdPalette(r io.Reader, board *Board) error {
	if err := ReadName(r, &board.Palette); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdFlags(r io.Reader, board *Board) error {
	// board flags
	if err := ReadByte(r, &board.Flags); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdCounters(r io.Reader, board *Board) error {
	// counters
	if err := ReadCounters(r, &board.Counters); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdTiles(r io.Reader, board *Board) error {
	// tiles
	if err := readTiles(r, &board.Tiles); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdObjects(r io.Reader, board *Board) error {
	// objects
	count := 0
	if err := ReadInt(r, &count); err != nil {
		return ErrBadBoardData
	}
	board.Objects = board.objects[:count]
	if err := readObjects(r, board.Objects); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdCamera(r io.Reader, board *Board) error {
	// camera
	if err := ReadVector(r, &board.Camera); err != nil {
		return ErrBadBoardData
	}
	if err := ReadInt(r, &board.Follow); err != nil {
		return ErrBadBoardData
	}
	if err := ReadInt(r, &board.ViewID); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdPlayer(r io.Reader, board *Board) error {
	if err := ReadInt(r, &board.Player); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdRunState(r io.Reader, board *Board) error {
	// runner
	if err := ReadUint8(r, &board.Tick); err != nil {
		return ErrBadBoardData
	}
	if err := ReadInt(r, &board.ObjectID); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readBrdEnter(r io.Reader, board *Board) error {
	if err := ReadVector(r, &board.Enter); err != nil {
		return ErrBadBoardData
	}
	return nil
}

func readTiles(r io.Reader, tiles *[BoardLength]Tile) error {
	// number of tile data to read
	length := uint16(0)
	if err := ReadUInt16(r, &length); err != nil {
		return err
	}

	pos := 0
	for i := uint16(0); i < length; i++ {
		// repeat tile count
		count := uint8(0)
		if err := ReadUint8(r, &count); err != nil {
			return err
		}

		// tile data
		cell := Cell{}
		if err := ReadCell(r, &cell); err != nil {
			return err
		}
		id := uint8(0)
		if err := ReadUint8(r, &id); err != nil {
			return err
		}

		// place tiles
		for j := uint8(0); j < count; j++ {
			tiles[pos] = Tile{Cell: cell, ID: id}
			pos++
		}
	}
	return nil
}

func readObjects(r io.Reader, objects []Object) error {
	// number of objects to read
	for i := range objects {
		if err := readObject(r, &objects[i]); err != nil {
			return err
		}
	}
	return nil
}

func WriteBoard(w io.Writer, board *Board) error {

	// number of segments
	WriteUint8(w, 11)

	// connected board ids
	WriteByte(w, brdConnections)
	if err := WriteName(w, board.North); err != nil {
		return err
	}
	if err := WriteName(w, board.South); err != nil {
		return err
	}
	if err := WriteName(w, board.East); err != nil {
		return err
	}
	if err := WriteName(w, board.West); err != nil {
		return err
	}

	// charset and palette
	WriteByte(w, brdCharset)
	if err := WriteName(w, board.Charset); err != nil {
		return err
	}

	WriteByte(w, brdPalette)
	if err := WriteName(w, board.Palette); err != nil {
		return err
	}

	// board flags
	WriteByte(w, brdFlags)
	if err := WriteByte(w, board.Flags); err != nil {
		return err
	}

	// counters
	WriteByte(w, brdCounters)
	if err := WriteCounters(w, board.Counters); err != nil {
		return err
	}

	// tiles
	WriteByte(w, brdTiles)
	if err := writeTiles(w, &board.Tiles); err != nil {
		return err
	}

	// objects
	WriteByte(w, brdObjects)
	if err := WriteInt(w, len(board.Objects)); err != nil {
		return err
	}
	if err := writeObjects(w, board.Objects); err != nil {
		return err
	}

	// camera
	WriteByte(w, brdCamera)
	if err := WriteVector(w, board.Camera); err != nil {
		return err
	}
	if err := WriteInt(w, board.Follow); err != nil {
		return err
	}
	if err := WriteInt(w, board.ViewID); err != nil {
		return err
	}

	WriteByte(w, brdPlayer)
	if err := WriteInt(w, board.Player); err != nil {
		return err
	}

	// run state
	WriteByte(w, brdRunState)
	if err := WriteUint8(w, board.Tick); err != nil {
		return err
	}
	if err := WriteInt(w, board.ObjectID); err != nil {
		return err
	}

	WriteByte(w, brdEnter)
	if err := WriteVector(w, board.Enter); err != nil {
		return err
	}

	return nil
}

func writeTiles(w io.Writer, tiles *[BoardLength]Tile) error {

	length := uint16(0)

	// write to buffer
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	for i := 0; i < BoardLength; i++ {
		tile := tiles[i]

		// count repeated tiles
		count := uint8(0)
		for j := i + 1; j < BoardLength; j++ {
			if tiles[j] != tile || count == 254 {
				break
			}
			count++
		}

		if err := WriteUint8(buf, count+1); err != nil {
			return err
		}
		if err := WriteCell(buf, tile.Cell); err != nil {
			return err
		}
		if err := WriteUint8(buf, tile.ID); err != nil {
			return err
		}

		i += int(count)
		length++
	}

	// write the tile data from buffer
	if err := WriteUInt16(w, length); err != nil {
		return err
	}
	if _, err := buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}

func writeObjects(w io.Writer, objs []Object) error {
	for i := range objs {
		if err := writeObject(w, objs[i]); err != nil {
			return err
		}
	}
	return nil
}

func TileID(value Vector) int {
	if value.X < 0 || value.Y < 0 || value.X >= BoardWidth || value.Y >= BoardHeight {
		return -1
	}
	return value.X + (value.Y << 7)
}

func TilePos(tid int) Vector {
	return Vector{
		X: tid & 127,
		Y: tid >> 7,
	}
}

func OutsideBoard(pos Vector) bool {
	return pos.X < 0 || pos.Y < 0 || pos.X >= BoardWidth || pos.Y >= BoardHeight
}

func (b *Board) ObjectsAll() iter.Seq2[int, *Object] {
	return func(yield func(int, *Object)bool) {
		for i := range b.Objects {
			obj := &b.Objects[i]
			if !obj.Free {
				if !yield (i, obj) {
					return
				}
			}
		}
	}
}

func (b *Board) ObjectsAt(pos Vector) iter.Seq2[int, *Object] {
	return func(yield func(int, *Object)bool) {
		for i := range b.Objects {
			obj := &b.Objects[i]
			if !obj.Free && obj.Position == pos {
				if !yield (i, obj) {
					return
				}
			}
		}
	}
}
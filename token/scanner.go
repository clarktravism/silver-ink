package token

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"silver-ink/data"
)

const line_size = 256

type NameRepo interface {
	GetName(data.Name) int
	IsElement(data.Name) bool
	IsPrefab(data.Name) bool
}

type Labeler interface {
	AddLabel(name data.Name, pos int)
}

type ByteCodeScanner struct {
	r  io.Reader // reader of raw code
	l  Labeler   // where to add labels
	nr NameRepo

	scanner    *bufio.Scanner // reads from r
	lineBuffer [line_size]byte
	buf        [line_size]byte
	wordBuffer [32][]byte
	pos        int
	line       int
}

func NewByteCodeScanner(r io.Reader, l Labeler, nr NameRepo) *ByteCodeScanner {
	bcs := ByteCodeScanner{
		r:  r,
		l:  l,
		nr: nr,
	}
	bcs.scanner = bufio.NewScanner(r)
	bcs.scanner.Buffer(bcs.lineBuffer[:], line_size)
	bcs.scanner.Split(bcs.ByteCodeLineScan)
	return &bcs
}

func (bcs *ByteCodeScanner) Reset() {
	bcs.pos = 0
	bcs.line = 0
}

func (bcs *ByteCodeScanner) Scan() bool {
	return bcs.scanner.Scan()
}

func (bcs *ByteCodeScanner) Err() error {
	return bcs.scanner.Err()
}

func (bcs *ByteCodeScanner) Bytes() []byte {
	return bcs.scanner.Bytes()
}

// assumes that ElementFile, PrefabFile, BoardsFile are available
func (bcs *ByteCodeScanner) ByteCodeLineScan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		token, err = bcs.byteCodeLine(data[0:i])
		bcs.line++
		bcs.pos += len(token)
		return i + 1, token, err
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		token, err = bcs.byteCodeLine(data)
		bcs.line++
		bcs.pos += len(token)
		return len(data), token, err
	}
	// Request more data.
	return 0, nil, nil
}

func (bcs *ByteCodeScanner) byteCodeLine(b []byte) (d []byte, err error) {
	b = dropCR(b)

	// blank new line
	if len(b) == 0 {
		// empty text
		bcs.buf[0] = CmdText
		bcs.buf[1] = 3
		bcs.buf[2] = 0
		return bcs.buf[0:3], nil
	}

	// comment
	if b[0] == '`' {
		return bcs.buf[:0], nil
	}

	// label
	if b[0] == ':' {
		return bcs.buf[:0], bcs.parseLabelLine(b[1:])
	}

	w := 0 // bytes written / length of token

	switch b[0] {
	case '#':
		_, seg := bcs.scanSegment(b[1:])
		_, wN, err := bcs.putSegment(seg, bcs.buf[w:])
		if err != nil {
			return nil, fmt.Errorf("error line %v: %v", bcs.line+1, err)
		}
		w += wN
	case '!', '$':
		_, seg := bcs.scanSegment(b)
		_, wN, err := bcs.putSegment(seg, bcs.buf[w:])
		if err != nil {
			return nil, fmt.Errorf("error line %v: %v", bcs.line+1, err)
		}
		w += wN
	case '/', '?':
		for len(b) > 0 {
			n, seg := bcs.scanSegment(b)
			b = b[n:]
			_, wN, err := bcs.putSegment(seg, bcs.buf[w:])
			if err != nil {
				return nil, fmt.Errorf("error line %v: %v", bcs.line+1, err)
			}
			w += wN
		}
	default:
		w, err = bcs.putTextLine(b, bcs.buf[:])
	}
	if err != nil {
		return nil, fmt.Errorf("error line %v: %v", bcs.line+1, err)
	}
	if w > 255 {
		return nil, fmt.Errorf("error line %v: line too long", bcs.line+1)
	}

	return bcs.buf[:w], nil
}

func (bcs *ByteCodeScanner) putSegment(b []byte, buf []byte) (r int, w int, err error) {

	// switch on first character
	switch b[0] {
	case '#':
		r, w, err = bcs.putCmd(b[1:], buf)
	case '!':
		w, err = bcs.putLinkLine(b[1:], buf)
		r = len(b)
	case '$':
		w, err = bcs.putCenterTextLine(b[1:], buf)
		r = len(b)
	case '/', '?':
		r, w, err = bcs.putShortCmd(b, buf)
	default:
		r, w, err = bcs.putCmd(b, buf)
		if r < len(b) {
			return 0, 0, ErrBadCommmand
		}
	}
	if err != nil {
		return 0, 0, err
	}

	return r, w, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func (bcs *ByteCodeScanner) putInt(buf []byte, v int) {
	binary.LittleEndian.PutUint16(buf, uint16(v))
}

func (bcs *ByteCodeScanner) putName(buf []byte, name data.Name) int {
	i := bcs.nr.GetName(name)
	binary.LittleEndian.PutUint16(buf, uint16(i))
	return 2
}

// line parses

var ErrSpaceInLabel = errors.New("unexpected character in label")

func (bcs *ByteCodeScanner) parseLabelLine(b []byte) error {
	if bytes.ContainsAny(b, " ") {
		return ErrSpaceInLabel
	}
	name := data.NewName(string(b))
	bcs.l.AddLabel(name, bcs.pos)
	return nil
}

// need to add command length
func (bcs *ByteCodeScanner) putTextLine(b []byte, buf []byte) (int, error) {
	bcs.buf[0] = CmdText
	n := copy(bcs.buf[3:], b)
	buf[1] = byte(n + 3)
	buf[2] = byte(n)
	return n + 3, nil
}

// need to add command length
func (bcs *ByteCodeScanner) putCenterTextLine(b []byte, buf []byte) (int, error) {
	bcs.buf[0] = CmdTextCenter
	n := copy(bcs.buf[3:], b)
	buf[1] = byte(n + 3)
	buf[2] = byte(n)
	return n + 3, nil
}

var ErrMissingLinkSep = errors.New("missing ';' on link")

// need to add command length
// cmd, len, labelname, len text, text
func (bcs *ByteCodeScanner) putLinkLine(b []byte, buf []byte) (int, error) {
	buf[0] = CmdTextLink
	sep := bytes.IndexAny(b, ";")
	if sep == -1 {
		return 0, ErrMissingLinkSep
	}
	name := data.NewName(string(b[:sep]))
	bcs.putName(buf[2:4], name)
	n := copy(buf[5:], b[sep+1:])
	buf[1] = byte(n + 5)
	buf[4] = byte(n)
	return n + 5, nil
}

func (bcs *ByteCodeScanner) putShortCmd(b []byte, buf []byte) (r int, w int, err error) {
	// get cmd byte
	cmd, ok := shortCmdCodes[b[0]]
	if !ok {
		return 0, 0, ErrBadParam
	}
	buf[w] = cmd
	w += 2
	r = 1

	if params, ok := cmdParams[cmd]; ok {
		rN, wN, err := bcs.putParams(b[r:], buf[w:], params)
		if err != nil {
			return 0, 0, err
		}
		r += rN
		w += wN
	}

	buf[1] = byte(w) // length of cmd + params

	return r, w, nil
}

func (bcs *ByteCodeScanner) putCmd(b []byte, buf []byte) (r int, w int, err error) {
	// get cmd byte
	r, word := bcs.scanWord(b)
	str := strings.ToUpper(string(word))
	cmd, ok := cmdCodes[str]
	if !ok {
		cmd = CmdGoto
		r = 0
	}
	buf[w] = cmd
	w += 2

	if params, ok := cmdParams[cmd]; ok {
		rN, wN, err := bcs.putParams(b[r:], buf[w:], params)
		if err != nil {
			return 0, 0, err
		}
		r += rN
		w += wN
	}

	buf[1] = byte(w) // length of cmd + params

	return r, w, nil
}

func (bcs *ByteCodeScanner) scanWord(b []byte) (int, []byte) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(b); start += width {
		var r rune
		r, width = utf8.DecodeRune(b[start:])
		if !isSpace(r) {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(b); i += width {
		var r rune
		r, width = utf8.DecodeRune(b[i:])
		if isSpace(r) {
			return i + width, b[start:i]
		}
	}
	// final word
	return len(b), b[start:]
}

func (bcs *ByteCodeScanner) scanSegment(b []byte) (int, []byte) {
	// Scan until segment end: '/' or '?'
	for width, i := 0, 0; i < len(b); i += width {
		var r rune
		r, width = utf8.DecodeRune(b[i:])
		if i > 0 && isSegmentStop(r) {
			return i, b[:i]
		}
	}
	// final word
	return len(b), b
}

func isSpace(r rune) bool {
	return r == ' '
}

func isSegmentStop(r rune) bool {
	return r == '/' || r == '?'
}

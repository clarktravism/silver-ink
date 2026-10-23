package token

import (
	"strings"
)

const (
	DirIdle byte = iota + 1
	DirNorth
	DirSouth
	DirEast
	DirWest
	DirFlow
	DirSeek
	DirRnd
	DirRndns
	DirRndne
	DirCW
	DirCCW
	DirOPP
	DirRnddp
)

var dirCodes = map[string]byte{
	"I":     DirIdle,
	"IDLE":  DirIdle,
	"N":     DirNorth,
	"NORTH": DirNorth,
	"S":     DirSouth,
	"SOUTH": DirSouth,
	"E":     DirEast,
	"EAST":  DirEast,
	"W":     DirWest,
	"WEST":  DirWest,
	"CW":    DirCW,
	"CCW":   DirCCW,
	"OPP":   DirOPP,
	"FLOW":  DirFlow,
	"SEEK":  DirSeek,
	"RND":   DirRnd,
	"RNDNS": DirRndns,
	"RNDNE": DirRndne,
	"RNDDP": DirRnddp,
}

var dirParams = map[byte][]int{
	DirCW:    {ParamDir},
	DirCCW:   {ParamDir},
	DirOPP:   {ParamDir},
	DirRnddp: {ParamDir},
}

func (bcs *ByteCodeScanner) putDirParam(b []byte, buf []byte) (r int, w int, err error) {
	r, word := bcs.scanWord(b)
	str := strings.ToUpper(string(word))
	code := dirCodes[str]
	buf[0] = ArgDir
	buf[1] = code
	w = 2
	if params, ok := dirParams[code]; ok {
		rN, wN, err := bcs.putParams(b[r:], buf[w:], params)
		if err != nil {
			return 0, 0, err
		}
		return r + rN, w + wN, nil
	}
	return r, w, nil
}
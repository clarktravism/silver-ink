package token

import (
	"strings"
)

const (
	PosSelf byte = iota + 1
	PosAt
	PosWhere
	PosBy
	PosTo
	PosPlayer
	PosFollow
)

var posCodes = map[string]byte{
	"SELF":   PosSelf,
	"AT":     PosAt,
	"WHERE":  PosWhere,
	"BY":     PosBy,
	"TO":     PosTo,
	"PLAYER": PosPlayer,
	"FOLLOW": PosFollow,
}

var posParams = map[byte][]int{
	PosAt:    {ParamInteger, ParamInteger},
	PosWhere: {ParamInteger, ParamInteger},
	PosBy:    {ParamObject},
	PosTo:    {ParamPos | ParamDir, ParamPos | ParamDir},
}

func (bcs *ByteCodeScanner) putPosParam(b []byte, buf []byte) (r int, w int, err error) {
	r, word := bcs.scanWord(b)
	str := strings.ToUpper(string(word))
	code := posCodes[str]
	buf[0] = ArgPos
	buf[1] = code
	w = 2
	if params, ok := posParams[code]; ok {
		rN, wN, err := bcs.putParams(b[r:], buf[w:], params)
		if err != nil {
			return 0, 0, err
		}
		return r + rN, w + wN, nil
	}
	return r, w, nil
}

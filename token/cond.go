package token

import (
	"strings"
)

const (
	CondBlocked byte = iota + 1
	CondContact
	CondNot
	CondAligned
	CondAny
	CondDetect
	CondAccepted

	CondIsTrue // for flags
	CondEqual
	CondNotEqual
	CondGT
	CondLT
	CondGTE
	CondLTE
	CondBetween
	CondAnd
	CondOr
	CondXor
)

var condCodes = map[string]byte{
	"BLOCKED": CondBlocked,
	"CONTACT": CondContact,
	"NOT":     CondNot,
	"ALLIGNED": CondAligned,
	"ALIGNED":  CondAligned,
	"ANY":      CondAny,
	"DETECT":   CondDetect,
	"ACCEPTED": CondAccepted,
}

var intCondCodes = map[string]byte{
	"IS":      CondEqual,
	"EQUAL":   CondEqual,
	"==":      CondEqual,
	"!=":      CondNotEqual,
	"<>":      CondNotEqual,
	">":       CondGT,
	">=":      CondGTE,
	"<":       CondLT,
	"<=":      CondLTE,
	"BETWEEN": CondBetween,
}

var condParams = map[byte][]int{
	CondBlocked:  {ParamDir | ParamPos},
	CondNot:      {ParamCond},
	CondAny:      {ParamInteger | ParamOptional, ParamElement | ParamObject},
	CondDetect:   {ParamDir | ParamPos, ParamInteger | ParamOptional, ParamElement | ParamObject},
	CondEqual:    {ParamInteger},
	CondNotEqual: {ParamInteger},
	CondGT:       {ParamInteger},
	CondGTE:      {ParamInteger},
	CondLT:       {ParamInteger},
	CondLTE:      {ParamInteger},
	CondBetween:  {ParamInteger, ParamInteger},
}

func (bcs *ByteCodeScanner) putCondParam(b []byte, buf []byte) (r int, w int, err error) {
	r, word := bcs.scanWord(b)
	str := strings.ToUpper(string(word))
	code := condCodes[str]
	buf[0] = ArgCond
	buf[1] = code
	w = 2
	if params, ok := condParams[code]; ok {
		rN, wN, err := bcs.putParams(b[r:], buf[w:], params)
		if err != nil {
			return 0, 0, err
		}
		return r + rN, w + wN, nil
	}
	return r, w, nil
}

// for cond like #IF COUNTER, #IF COUNTER >= INT
// gets tokenized to <cond> <params>
func (bcs *ByteCodeScanner) putIntCondParam(b []byte, buf []byte) (r int, w int, err error) {
	w = 2

	// read leading int param
	// skips two bytes the write, so that there is room for the <cond>
	nR, nW, err := bcs.putIntegerParam(b, buf[w:])
	if err != nil {
		return 0, 0, err
	}
	w += nW
	r += nR

	// get next word
	nR, word := bcs.scanWord(b[r:])
	if len(word) == 0 {
		return 0, 0, ErrBadParam
	}

	// get the cond byte
	code, ok := intCondCodes[string(word)]
	if !ok { // no match, so this is a simple flag check
		buf[0] = ArgCond
		buf[1] = CondIsTrue
		return r, w, nil
	} else {
		r += nR
	}

	buf[0] = ArgCond
	buf[1] = code

	// get rest of the params for this code
	// remember that we already got the leading int param
	params, ok := condParams[code]
	if !ok {
		return 0, 0, ErrBadParam
	}

	nR, nW, err = bcs.putParams(b[r:], buf[w:], params)
	if err != nil {
		return 0, 0, err
	}
	w += nW
	r += nR

	return r, w, nil
}

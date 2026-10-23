package token

import (
	"strings"
)

const (
	MathAdd byte = iota + 1
	MathSub
	MathMul
	MathDiv
	MathMod
	MathRandom
	MathAbs
	MathDist
	MathDiff
)

var mathCodes = map[string]byte{
	"ADD":    MathAdd,
	"SUB":    MathSub,
	"MUL":    MathMul,
	"DIV":    MathDiv,
	"MOD":    MathMod,
	"RANDOM": MathRandom,
	"ABS":    MathAbs,
	"DIST":   MathDist,
	"DIFF":   MathDiff,
}

var mathParams = map[byte][]int{
	MathAdd:    {ParamInteger, ParamInteger},
	MathSub:    {ParamInteger, ParamInteger},
	MathDiv:    {ParamInteger, ParamInteger},
	MathMod:    {ParamInteger, ParamInteger},
	MathRandom: {ParamInteger, ParamInteger},
	MathAbs:    {ParamInteger},
	MathDist:   {ParamPos, ParamPos},
	MathDiff:   {ParamInteger, ParamInteger},
}

func (bcs *ByteCodeScanner) putMathParam(b []byte, buf []byte) (r int, w int, err error) {
	r, word := bcs.scanWord(b)
	str := strings.ToUpper(string(word))
	code := mathCodes[str]
	buf[0] = ArgMath
	buf[1] = code
	w = 2
	if params, ok := mathParams[code]; ok {
		rN, wN, err := bcs.putParams(b[r:], buf[w:], params)
		if err != nil {
			return 0, 0, err
		}
		return r + rN, w + wN, nil
	}
	return r, w, nil
}
package token

import (
	"bytes"
	"errors"
	"silver-ink/data"
	"strconv"
	"strings"
	"unicode"
)

const (
	ParamOptional int = 1 << iota
	ParamInteger
	ParamCounter
	ParamScopedLabel
	ParamLabel
	ParamPos
	ParamDir
	ParamFlag
	ParamElement
	ParamPrefab
	ParamTarget
	ParamCond
	ParamCmd
	ParamSound
	ParamObject
	ParamObjectName
	ParamName
	ParamText
)

var ErrMissingParam = errors.New("missing param")
var ErrBadCommmand = errors.New("bad command or too many params")

func (bcs *ByteCodeScanner) putParams(b []byte, buf []byte, params []int) (r int, w int, err error) {
	for _, param := range params {
		nR, nW, err := bcs.putParam(b[r:], buf[w:], param)
		if err != nil {
			return 0, 0, err
		}
		r += nR
		w += nW
	}
	return r, w, nil
}

// return read bytes, wrote bytes
var ErrBadParam = errors.New("bad param")

func (bcs *ByteCodeScanner) putParam(b []byte, buf []byte, param int) (r int, w int, err error) {

	optional := param&ParamOptional == ParamOptional

	// end of line
	if len(b) == 0 {
		if optional {
			return 0, 0, nil
		}
		return 0, 0, ErrMissingParam
	}

	// text - always to the end of the line
	if param&ParamText == ParamText {
		buf[0] = ArgText
		buf[1] = byte(len(b))
		w = 2
		w += copy(buf[2:], b)
		return len(b), w, nil
	}

	r, word := bcs.scanWord(b)
	str := strings.ToUpper(string(word))
	name := data.NewName(str)

	// element
	if param&ParamElement == ParamElement {
		// check element
		if bcs.nr.IsElement(name) {
			buf[0] = ArgElement
			bcs.putName(buf[1:], name)
			return r, 3, nil
		}
	}

	// prefab
	if param&ParamPrefab == ParamPrefab {
		// check prefab
		if bcs.nr.IsPrefab(name) {
			buf[0] = ArgPrefab
			bcs.putName(buf[1:], name)
			return r, 3, nil
		}
	}

	// command
	if param&ParamCmd == ParamCmd {
		if str == "THEN" {
			b = b[r:]
			if len(b) == 0 {
				return 0, 0, ErrMissingParam
			}
		} else {
			r = 0
		}
		if b[0] == '#' {
			b = b[1:]
			r = 1
			if len(b) == 0 {
				return 0, 0, ErrMissingParam
			}
		}

		buf[0] = ArgCmd
		rN, w, err := bcs.putCmd(b, buf[1:])
		if err != nil {
			return 0, 0, ErrBadParam
		}
		return r + rN, w + 1, nil
	}

	if param&ParamFlag == ParamFlag {
		if code, ok := flagCodes[str]; ok {
			buf[0] = ArgFlag
			buf[1] = code
			return r, 2, nil
		}
	}

	// scope:label
	if param&ParamScopedLabel == ParamScopedLabel {
		if i := bytes.IndexAny(word, ":"); i != -1 {
			w, err = bcs.splitScopeLabel(i, word, buf)
			return r, w, err
		}
	}

	// target
	if param&ParamTarget == ParamTarget {
		buf[0] = ArgTarget
		return r, 1, nil
	}

	// label
	if param&ParamLabel == ParamLabel {
		buf[0] = ArgLabel
		bcs.putName(buf[1:], name)
		return r, 3, nil
	}

	// object name
	if param&ParamObjectName == ParamObjectName {
		buf[0] = ArgObjectName
		bcs.putName(buf[1:], name)
		return r, 3, nil
	}

	// generic name
	if param&ParamName == ParamName {
		buf[0] = ArgName
		bcs.putName(buf[1:], name)
		return r, 3, nil
	}

	// dir
	if param&ParamDir == ParamDir {
		if _, ok := dirCodes[str]; ok {
			return bcs.putDirParam(b, buf)
		}
	}

	// flag
	if param&ParamPos == ParamPos {
		if _, ok := posCodes[str]; ok {
			return bcs.putPosParam(b, buf)
		}
	}

	// cond
	if param&ParamCond == ParamCond {
		// normal cond
		if _, ok := condCodes[str]; ok {
			return bcs.putCondParam(b, buf)
		}
		// int cond
		return bcs.putIntCondParam(b, buf)
	}

	// object
	if param&ParamObject == ParamObject {
		if code, ok := objectCodes[str]; ok {
			buf[0] = ArgObject
			buf[1] = code
			return r, 2, nil
		}
		buf[0] = ArgObjectName
		w = bcs.putName(buf[1:], name)
		return r, w + 1, nil
	}

	// counter
	if param&ParamCounter == ParamCounter {
		return bcs.putCounterParam(b, buf)
	}

	// to handle optional color before element|prefab
	// so, element and prefab can never be optional
	if bcs.isElementOrPrefab(name) {
		if optional { // move to next param
			return 0, 0, nil
		}
		return 0, 0, ErrBadParam
	}

	// integer
	if param&ParamInteger == ParamInteger {
		return bcs.putIntegerParam(b, buf)
	}

	if optional {
		return 0, 0, nil
	}

	return 0, 0, ErrBadParam
}

func (bcs *ByteCodeScanner) isElementOrPrefab(name data.Name) bool {
	return bcs.nr.IsElement(name) || bcs.nr.IsPrefab(name)
}

func (bcs *ByteCodeScanner) putIntegerParam(b []byte, buf []byte) (r int, w int, err error) {
	r, word := bcs.scanWord(b)
	str := string(word)

	// math op
	if _, ok := mathCodes[str]; ok {
		return bcs.putMathParam(b, buf)
	}

	// number literal
	if unicode.IsDigit(rune(word[0])) || word[0] == '-' {
		if v, err := strconv.Atoi(str); err == nil {
			buf[0] = ArgInteger
			bcs.putInt(buf[1:], v)
			return r, 3, nil
		} else {
			return 0, 0, ErrBadParam
		}
	}
	// hex number literal
	if word[0] == 'x' {
		if v, err := strconv.ParseUint(str[1:], 16, 16); err == nil {
			buf[0] = ArgInteger
			bcs.putInt(buf[1:], int(v))
			return r, 3, nil
		} else {
			return 0, 0, ErrBadParam
		}
	}

	// scope:counter
	if i := bytes.IndexAny(word, ":"); i != -1 {
		w, err = bcs.splitCounter(i, word, buf)
		return r, w, err
	}

	name := data.NewName(str)
	// color
	if color, ok := data.ColorNames[name]; ok {
		buf[0] = ArgInteger
		bcs.putInt(buf[1:], int(color))
		return r, 3, nil
	}

	// counter
	buf[0] = ArgCounter
	bcs.putName(buf[1:], name)
	return r, 3, nil
}

package token

import (
	"bytes"
	"strings"
	"silver-ink/data"
)

var scopeCounterCodes = map[string]byte{
	"SELF":   ScopeSelf,
	"BOARD":  ScopeBoard,
}

func (bcs *ByteCodeScanner) putCounterParam(b []byte, buf []byte) (r int, w int, err error) {
	r, word := bcs.scanWord(b)
	str := string(word)

	// scope:counter
	if i := bytes.IndexAny(word, ":"); i != -1 {
		w, err = bcs.splitCounter(i, word, buf)
		return r, w, err
	}

	name := data.NewName(str)

	// counter
	buf[0] = ArgCounter
	bcs.putName(buf[1:], name)
	return r, 3, nil
}

func (bcs *ByteCodeScanner) splitCounter(i int, word []byte, buf []byte) (w int, err error) {
	if i == 0 {
		return 0, ErrBadParam
	}

	// scope
	scope := strings.ToUpper(string(word[:i]))
	code, ok := scopeCounterCodes[scope]
	if !ok {
		buf[0] = ArgScopeName
		bcs.putName(buf[1:], data.NewName(scope))
		w += 3
	} else {
		buf[0] = ArgScope
		buf[1] = code
		w += 2
	}

	// name
	name := data.NewName(string(word[i+1:]))
	buf[w] = ArgCounter
	bcs.putName(buf[w+1:], name)
	w += 3
	return w, nil
}
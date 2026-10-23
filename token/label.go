package token

import (
	"strings"
	"silver-ink/data"
)

var scopeLabelCodes = map[string]byte{
	"SELF":   ScopeSelf,
	"THIS":   ScopeSelf,
	"ALL":    ScopeAll,
	"OTHERS": ScopeOthers,
	"PLAYER": ScopePlayer,
	"GROUP" : ScopeGroup,
}

func (bcs *ByteCodeScanner) splitScopeLabel(i int, word []byte, buf []byte) (w int, err error) {
	if i == 0 {
		return 0, ErrBadParam
	}
	// scope
	scope := strings.ToUpper(string(word[:i]))
	code, ok := scopeLabelCodes[scope]
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
	buf[w] = ArgLabel
	bcs.putName(buf[w+1:], name)
	w += 3
	return w, nil
}
package token

const (
	ObjNone byte = iota + 1
	ObjSelf
	ObjPlayer
	ObjGroup
)

var objectCodes = map[string]byte{
	"NONE":   ObjNone,
	"SELF":   ObjSelf,
	"PLAYER": ObjPlayer,
	"GROUP": ObjGroup,
}

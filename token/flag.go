package token

const (
	FlagCollision1 byte = iota + 1
	FlagCollision2
	FlagCollision3
	FlagCollision4
	FlagCollision5
	FlagCollision6
	FlagCollision7
	FlagCollision8

	FlagCanPush
	FlagPushableN
	FlagPushableS
	FlagPushableE
	FlagPushableW
	FlagPushableAll
	FlagPushableNS
	FlagPushableEW

	FlagPriority
	FlagShowBG
)

var flagCodes = map[string]byte{
	"COLLISION_1" : FlagCollision1,
	"COLLISION_2" : FlagCollision2,
	"COLLISION_3" : FlagCollision3,
	"COLLISION_4" : FlagCollision4,
	"COLLISION_5" : FlagCollision5,
	"COLLISION_6" : FlagCollision6,
	"COLLISION_7" : FlagCollision7,
	"COLLISION_8" : FlagCollision8,
	"COLLISION_WALLS" : FlagCollision1,
	"COLLISION_WATER" : FlagCollision2,
	"COLLISION_FLOOR" : FlagCollision3,
	
	"CAN_PUSH" : FlagCanPush,
	"PUSHABLE_N" : FlagPushableN,
	"PUSHABLE_S" : FlagPushableS,
	"PUSHABLE_E" : FlagPushableE,
	"PUSHABLE_W" : FlagPushableW,
	"PUSHABLE_NS" : FlagPushableNS,
	"PUSHABLE_EW" : FlagPushableEW,
	"PUSHABLE_ALL" : FlagPushableAll,

	"PRIORITY" : FlagPriority,
	"SHOW_BG" : FlagShowBG,
}
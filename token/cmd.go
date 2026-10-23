package token

// Command ByteCodes
const (
	ByteZero byte = iota // used for unavailable/invalid values

	// text
	CmdText       // len(8) text..
	CmdTextCenter // len(8) text..
	CmdTextLink   // int(16 name) len(8) text..

	// commands
	CmdAnd
	CmdBecome
	CmdBoard
	CmdCamera
	CmdFollow
	CmdChange
	CmdChar
	CmdClear
	CmdCollision
	CmdColor
	CmdCycle
	CmdDie
	CmdDuplicate
	CmdElse
	CmdEnd
	CmdEndGame
	CmdEvent
	CmdGive
	CmdGo
	CmdGoShort
	CmdGoto
	CmdGroup
	CmdIdle
	CmdIf
	CmdLock
	CmdLog
	CmdName
	CmdPut
	CmdRepeat
	CmdRestart
	CmdRestore
	CmdScript
	CmdSend
	CmdSet
	CmdView
	CmdShoot
	CmdTake
	CmdTarget
	CmdTextBox
	CmdTile
	CmdTry
	CmdTryShort
	CmdUnlock
	CmdWalk
	CmdWrite
	CmdZap
	CmdDebug
)

var cmdCodes = map[string]byte{
	"AND":       CmdAnd,
	"BECOME":    CmdBecome,
	"BOARD":     CmdBoard,
	"CAMERA":    CmdCamera,
	"CHANGE":    CmdChange,
	"CHAR":      CmdChar,
	"CLEAR":     CmdClear,
	"COLLISION": CmdCollision,
	"COLOR":     CmdColor,
	"CYCLE":     CmdCycle,
	"DIE":       CmdDie,
	"DUPLICATE": CmdDuplicate,
	"ELSE":      CmdElse,
	"END":       CmdEnd,
	"ENDGAME":   CmdEndGame,
	"EVENT":     CmdEvent,
	"FOLLOW":    CmdFollow,
	"GIVE":      CmdGive,
	"GO":        CmdGo,
	"GOTO":      CmdGoto,
	"GROUP":     CmdGroup,
	"IDLE":      CmdIdle,
	"IF":        CmdIf,
	"LOCK":      CmdLock,
	"LOG":       CmdLog,
	"NAME":      CmdName,
	"PUT":       CmdPut,
	"REPEAT":    CmdRepeat,
	"RESTART":   CmdRestart,
	"RESTORE":   CmdRestore,
	"SCRIPT":    CmdScript,
	"SEND":      CmdSend,
	"SET":       CmdSet,
	"VIEW":      CmdView,
	"SHOOT":     CmdShoot,
	"TAKE":      CmdTake,
	"TARGET":    CmdTarget,
	"TEXTBOX":   CmdTextBox,
	"TRY":       CmdTry,
	"UNLOCK":    CmdUnlock,
	"WALK":      CmdWalk,
	"ZAP":       CmdZap,
	"DEBUG":     CmdDebug,
}

var shortCmdCodes = map[byte]byte{
	'/': CmdGoShort,
	'?': CmdTryShort,
}

var cmdParams = map[byte][]int{
	CmdAnd:       {ParamCmd},
	CmdBecome:    {ParamInteger | ParamOptional, ParamElement | ParamPrefab},
	CmdBoard:     {ParamName | ParamTarget, ParamInteger | ParamOptional, ParamInteger | ParamOptional},
	CmdCamera:    {ParamInteger, ParamInteger},
	CmdChange:    {ParamInteger | ParamOptional, ParamElement, ParamInteger | ParamOptional, ParamElement},
	CmdChar:      {ParamInteger},
	CmdClear:     {ParamCounter | ParamFlag},
	CmdCollision: {ParamInteger, ParamInteger},
	CmdColor:     {ParamInteger},
	CmdCycle:     {ParamInteger},
	CmdDie:       {},
	CmdDuplicate: {ParamDir | ParamPos | ParamObject, ParamDir | ParamPos},
	CmdElse:      {ParamCmd},
	CmdEnd:       {},
	CmdEvent:     {ParamDir | ParamPos | ParamOptional, ParamScopedLabel | ParamLabel, ParamInteger | ParamOptional},
	CmdFollow:    {ParamObject},
	CmdGive:      {ParamCounter, ParamInteger},
	CmdGo:        {ParamDir | ParamPos},
	CmdGoShort:   {ParamDir},
	CmdGoto:      {ParamLabel, ParamInteger | ParamOptional},
	CmdGroup:     {ParamObjectName},
	CmdIdle:      {},
	CmdIf:        {ParamCond, ParamCmd},
	CmdLock:      {},
	CmdLog:       {ParamText},
	CmdName:      {ParamObjectName},
	CmdPut:       {ParamDir | ParamPos, ParamInteger | ParamOptional, ParamElement | ParamPrefab},
	CmdRepeat:    {ParamInteger, ParamCmd},
	CmdRestart:   {},
	CmdRestore:   {ParamLabel},
	CmdScript:    {ParamName},
	CmdSend:      {ParamDir | ParamPos | ParamOptional, ParamScopedLabel | ParamLabel, ParamInteger | ParamOptional},
	CmdSet:       {ParamCounter | ParamFlag, ParamInteger | ParamOptional},
	CmdView:      {ParamInteger},
	CmdShoot:     {ParamDir, ParamInteger | ParamOptional, ParamPrefab},
	CmdTake:      {ParamCounter, ParamInteger, ParamCmd | ParamOptional},
	CmdTarget:    {ParamName},
	CmdTextBox:   {ParamInteger, ParamText | ParamOptional },
	CmdTile:      {ParamDir | ParamPos, ParamInteger | ParamOptional, ParamInteger},
	CmdTry:       {ParamDir | ParamPos, ParamCmd | ParamOptional},
	CmdTryShort:  {ParamDir | ParamPos},
	CmdUnlock:    {},
	CmdWalk:      {ParamDir},
	CmdWrite:     {ParamDir | ParamPos, ParamInteger, ParamText},
	CmdZap:       {ParamLabel},
}

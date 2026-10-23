package main

import (
	"errors"
	"log"
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/token"
)

var ErrMissingElementOrPrefab = errors.New("missing element or prefab")
var ErrBoardNotFound = errors.New("board not found")
var ErrBadParam = errors.New("bad or missing param")
var ErrBadToken = errors.New("bad token")

const command_limit = 100

// sends event to all
func Event(event data.Name, skip int) bool {
	result := false
	for i, o := range engine.Board.Objects {
		if !o.Free {
			if ctx := NewObjectEventContext(i, event, skip); ctx != nil {
				RunScript(ctx)
				result = true
			}
		}
	}
	for i := range data.BoardLength {
		result = TileEvent(i, event, skip) || result
	}
	return result
}

func EventAt(pos data.Vector, event data.Name, skip int) bool {
	result := false
	for i, o := range engine.Board.Objects {
		if !o.Free && o.Position == pos {
			if ctx := NewObjectEventContext(i, event, skip); ctx != nil {
				RunScript(ctx)
				result = true
			}
		}
	}
	result = TileEvent(data.TileID(pos), event, skip) || result
	return result
}

func TileEvent(tid int, event data.Name, skip int) bool {
	if ctx := NewTileContext(tid, event, skip); ctx != nil {
		RunScript(ctx)
		return true
	}
	return false
}

func ObjectEvent(id int, event data.Name, skip int) bool {
	if ctx := NewObjectEventContext(id, event, skip); ctx != nil {
		RunScript(ctx)
		return true
	}
	return false
}

func SetObjectScriptPosition(i int, label data.Name, skip int) int {
	obj := &engine.Board.Objects[i]
	if obj.Lock != 0 {
		return obj.ScriptPosition
	}
	switch label {
	case RESTART:
		obj.ScriptPosition = 0
	case END:
		obj.ScriptPosition = -1
	default:
		obj.ScriptPosition = obj.LabelPosition(label, skip)
	}
	return obj.ScriptPosition
}

func RunScript(ctx *ScriptContext) {
	defer ctx.Close()
	for !ctx.Stop && !TextBoxIsOpen {
		ctx.runNext()
	}
	if ctx.Err != nil {
		log.Println(ctx.Err)
		return
	}
	if !ctx.Event {
		ctx.SetScriptPosition(ctx.CurrentPos)
	}
}

func (ctx *ScriptContext) runNext() {
	ctx.CommandsRun++
	if ctx.CommandsRun >= command_limit {
		ctx.Stop = true
		return
	}
	if ctx.CurrentPos >= len(ctx.Script.Data) {
		ctx.CurrentPos = -1
		ctx.Stop = true
		return
	}
	if ctx.CurrentPos < 0 {
		ctx.Stop = true
		return
	}
	ctx.CmdStartPos = ctx.CurrentPos
	ctx.CmdEndPos = ctx.CurrentPos + int(ctx.Script.Data[ctx.CurrentPos+1])
	cmd := ctx.Script.Data[ctx.CurrentPos]
	ctx.ArgPos = ctx.CurrentPos + 2
	switch cmd {
	case token.CmdText:
		ctx.text()
	case token.CmdTextCenter:
		ctx.textCenter()
	case token.CmdTextLink:
		ctx.textLink()
	case token.CmdAnd:
		ctx.cmdAnd()
	case token.CmdBecome:
		ctx.cmdBecome()
	case token.CmdBoard:
		ctx.cmdBoard()
	case token.CmdCamera:
		ctx.cmdCamera()
	case token.CmdChar:
		ctx.cmdChar()
	case token.CmdClear:
		ctx.cmdClear()
	case token.CmdCollision:
		ctx.cmdCollision()
	case token.CmdColor:
		ctx.cmdColor()
	case token.CmdCycle:
		ctx.cmdCycle()
	case token.CmdDie:
		ctx.cmdDie()
	case token.CmdDuplicate:
		ctx.cmdDuplicate()
	case token.CmdElse:
		ctx.cmdElse()
	case token.CmdEnd:
		ctx.cmdEnd()
	case token.CmdEndGame:
		ctx.cmdEndGame()
	case token.CmdEvent:
		ctx.cmdEvent()
	case token.CmdGive:
		ctx.cmdGive()
	case token.CmdFollow:
		ctx.cmdFollow()
	case token.CmdGo:
		ctx.cmdGo()
	case token.CmdGoShort:
		ctx.cmdGo()
	case token.CmdGoto:
		ctx.cmdGoto()
	case token.CmdGroup:
		ctx.cmdGroup()
	case token.CmdIdle:
		ctx.cmdIdle()
	case token.CmdIf:
		ctx.cmdIf()
	case token.CmdLock:
		ctx.cmdLock()
	case token.CmdLog:
		ctx.cmdLog()
	case token.CmdName:
		ctx.cmdName()
	case token.CmdPut:
		ctx.cmdPut()
	case token.CmdRepeat:
		ctx.cmdRepeat()
	case token.CmdRestart:
		ctx.cmdRestart()
	case token.CmdRestore:
		ctx.cmdRestore()
	case token.CmdScript:
		ctx.cmdScript()
	case token.CmdSend:
		ctx.cmdSend()
	case token.CmdSet:
		ctx.cmdSet()
	case token.CmdView:
		ctx.cmdView()
	case token.CmdShoot:
		ctx.cmdShoot()
	case token.CmdTake:
		ctx.cmdTake()
	case token.CmdTarget:
		ctx.cmdTarget()
	case token.CmdTextBox:
		ctx.cmdTextBox()
	case token.CmdTile:
		ctx.cmdTile()
	case token.CmdTry:
		ctx.cmdTry()
	case token.CmdTryShort:
		ctx.cmdTryShort()
	case token.CmdUnlock:
		ctx.cmdUnlock()
	case token.CmdWalk:
		ctx.cmdWalk()
	case token.CmdZap:
		ctx.cmdZap()
	case token.CmdDebug:
		ctx.advance()
	default:
		ctx.Stop = true
		ctx.Error(ErrBadToken)
	}
}

func (ctx *ScriptContext) advance() {
	if ctx.PopRepeat() {
		ctx.rollback()
		return
	}
	ctx.CurrentPos = ctx.CmdEndPos
	if ctx.CurrentPos >= len(ctx.Script.Data)-1 {
		ctx.cmdEnd()
	}
}

func (ctx *ScriptContext) rollback() {
	ctx.CurrentPos = ctx.CmdStartPos
}

// commands
func (ctx *ScriptContext) text() {
	ctx.advance()
	b := ctx.NextText()
	Text.AddLine(engine.ScrollBoxLine{
		Text: string(b),
	})
}

func (ctx *ScriptContext) textCenter() {
	ctx.advance()
	b := ctx.NextText()
	Text.AddLine(engine.ScrollBoxLine{
		Text:     string(b),
		Centered: true,
	})
}

func (ctx *ScriptContext) textLink() {
	ctx.advance()
	name := ctx.NextName()
	b := ctx.NextText()
	Text.AddLine(engine.ScrollBoxLine{
		Text:     string(b),
		Link:     name,
		Centered: true,
	})
}

func (ctx *ScriptContext) cmdAnd() {
	ctx.advance()

	if !ctx.LastResult {
		return
	}

	arg1 := ctx.NextArg()
	if arg1.Code == token.ArgCmd {
		ctx.CurrentPos = ctx.ArgPos
		ctx.runNext()
	} else {
		ctx.Error(ErrBadParam)
		return
	}
}

func (ctx *ScriptContext) cmdBecome() {
	ctx.advance()

	// optional color
	color := ctx.Cell().Color
	arg := ctx.NextArg()
	if arg.CanInt() {
		if v, ok := ctx.GetArgInt(arg); ok {
			color = data.Color(v)
		}
		arg = ctx.NextArg()
	}

	if arg.Code == token.ArgElement {
		ctx.Die()
		element := arg.GetElement()
		// create tile
		tile := data.Tile{
			Cell: element.Cell,
			ID:   element.ID,
		}
		tile.Color = color

		// set tile
		tid := data.TileID(ctx.Position())
		engine.Board.Tiles[tid] = tile
		return
	}
	if arg.Code == token.ArgPrefab {
		ctx.Die()
		prefab := arg.GetPrefab()

		// create object
		obj := prefab.Object()
		obj.Position = ctx.Position()
		obj.Color = color

		id := engine.Board.AddObject(obj, prefab.Player, prefab.Follow)
		ObjectEvent(id, READY, 0)
		return
	}

	ctx.Error(ErrMissingElementOrPrefab)
}

func (ctx *ScriptContext) cmdBoard() {
	ctx.advance()
	ctx.Stop = true

	enter := data.VectorMinus

	argBrd := ctx.NextArg()
	if !argBrd.CanName() {
		ctx.Error(ErrBadParam)
	}

	argX := ctx.NextArg()
	if argX.CanInt() {
		enter.X, _ = ctx.GetArgInt(argX)
	}

	argY := ctx.NextArg()
	if argY.CanInt() {
		enter.Y, _ = ctx.GetArgInt(argY)
	}

	board, err := ReadBoardFile(argBrd.Name)
	if err != nil {
		ctx.Error(err)
		return
	}

	if enter == data.VectorMinus {
		if id := ctx.ObjectID(); id != -1 {
			name := engine.Board.Objects[id].Name
			group := engine.Board.Objects[id].Group
			// find other object with name and group
			for _, obj := range board.Objects {
				if obj.Name == name && obj.Group == group {
					enter = obj.Position
					break
				}
			}
		}
	}

	QueueSwapBoard(board, enter)
}

func (ctx *ScriptContext) cmdCamera() {
	ctx.advance()
	// close in case it was following an object
	engine.Board.Follow = -1

	arg := ctx.NextArg()
	if pos, ok := ctx.GetArgPos(arg); ok {
		engine.Board.Camera = pos
	}
}

// [color] <element> [color] <element>
func (ctx *ScriptContext) cmdChange() {
	ctx.advance()
	color1, element1, ok := ctx.NextColorElement()
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}
	color2, element2, ok := ctx.NextColorElement()
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}
	for i := range engine.Board.Tiles {
		tile := engine.Board.Tiles[i]
		if tile.ID != element1.ID {
			continue
		}
		if color1 != nil && tile.Color != *color1 {
			continue
		}
		engine.Board.Tiles[i].ID = element2.ID
		engine.Board.Tiles[i].Char = element2.Char
		if color2 != nil {
			engine.Board.Tiles[i].Color = *color2
		}
	}
}

// <char>
func (ctx *ScriptContext) cmdChar() {
	ctx.advance()
	arg := ctx.NextArg()
	if v, ok := ctx.GetArgInt(arg); ok {
		cell := ctx.Cell()
		cell.Char = byte(v)
	}
}

func (ctx *ScriptContext) cmdClear() {
	ctx.advance()
	arg := ctx.NextArg()
	if arg.Code == token.ArgFlag {
		ctx.ClearFlag(arg.Byte)
		return
	}
	if c, ok := ctx.GetArgCounter(arg); ok && c != nil {
		c.Value = 0
	}
}

func (ctx *ScriptContext) cmdCollision() {
	ctx.advance()

	c := ctx.Collision()

	arg1 := ctx.NextArg()
	value1, ok := ctx.GetArgInt(arg1)
	if !ok {
		return
	}

	p := data.BitField(1 >> (value1 - 1))

	arg2 := ctx.NextArg()
	value2, ok := ctx.GetArgInt(arg2)
	if !ok {
		return
	}

	if value2 > 1 {
		c.Set(p)
	} else {
		c.Clear(p)
	}
}

func (ctx *ScriptContext) cmdColor() {
	ctx.advance()
	arg := ctx.NextArg()
	if v, ok := ctx.GetArgInt(arg); ok {
		cell := ctx.Cell()
		cell.Color = data.Color(v)
	}
}

func (ctx *ScriptContext) cmdCycle() {
	ctx.advance()
	arg := ctx.NextArg()
	if v, ok := ctx.GetArgInt(arg); ok {
		if id := ctx.ObjectID(); id != -1 {
			engine.Board.Objects[id].Cycle = uint8(v)
		}
	}
}

func (ctx *ScriptContext) cmdDie() {
	ctx.CurrentPos = -1
	ctx.Stop = true
	ctx.Die()
}

// <pos|dir|object> <pos|dir>
func (ctx *ScriptContext) cmdDuplicate() {
	ctx.advance()

	i := -1

	// get object to duplicate
	arg1 := ctx.NextArg()
	if arg1.CanPos() {
		pos, _ := ctx.GetArgPos(arg1)
		i = engine.Board.ObjectAt(pos)
	} else if arg1.CanObject() {
		i = ctx.GetArgObject(arg1)
	} else {
		ctx.Error(ErrBadParam)
		return
	}
	if i == -1 {
		return
	}

	// copy object
	obj := engine.Board.Objects[i]
	obj.ScriptPosition = 0
	obj.Counters = engine.Board.Objects[i].Counters.Copy()
	obj.ClearAllZaps()
	obj.Lock = 0
	obj.Repeat = 0

	// get position
	dst, ok := ctx.GetArgPos(ctx.NextArg())
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}
	obj.Position = dst

	// check collision
	if CollisionAt(dst, obj.Collision, -1) {
		return
	}

	engine.Board.AddObject(obj, false, false)
}

func (ctx *ScriptContext) cmdElse() {
	ctx.advance()

	if ctx.LastResult {
		return
	}

	arg1 := ctx.NextArg()
	if arg1.Code == token.ArgCmd {
		ctx.CurrentPos = ctx.ArgPos
		ctx.runNext()
	} else {
		ctx.Error(ErrBadParam)
		return
	}
}

func (ctx *ScriptContext) cmdEnd() {
	ctx.CurrentPos = -1
	ctx.Stop = true
}

func (ctx *ScriptContext) cmdEndGame() {
	ctx.CurrentPos = -1
	ctx.Stop = true
	PushGameOver()
}

func (ctx *ScriptContext) cmdEvent() {
	ctx.advance()
	arg1 := ctx.NextArg()
	ctx.LastEvent = false

	// <label> [skip]
	if arg1.Code == token.ArgLabel {
		skip, _ := ctx.GetArgInt(ctx.NextArg())
		ctx.LastEvent = ctx.RunEvent(arg1.Name, skip)
		return
	}
	// SELF:<label> [skip]
	if arg1.Code == token.ArgScope && arg1.Byte == token.ScopeSelf {
		label := ctx.NextArg().Name
		skip, _ := ctx.GetArgInt(ctx.NextArg())
		ctx.LastEvent = ctx.RunEvent(label, skip)
		return
	}

	// [pos] option
	pos := data.VectorMinus
	if arg1.CanPos() {
		pos, _ = ctx.GetArgPos(arg1)
		arg1 = ctx.NextArg()
	}

	// [name|group:]<label> [skip]
	if arg1.Code == token.ArgScopeName {
		label := ctx.NextArg().Name
		if label.Empty() {
			ctx.Error(ErrBadParam)
			return
		}

		skip, _ := ctx.GetArgInt(ctx.NextArg())

		iter := engine.Board.ObjectsAll()
		if pos != data.VectorMinus { // apply [pos]
			iter = engine.Board.ObjectsAt(pos)
		}

		for i, obj := range iter {
			if obj.Name == arg1.Name || obj.Group == arg1.Name {
				ctx.LastEvent = ObjectEvent(i, label, skip) || ctx.LastEvent
			}
		}
		return
	}

	// [scope:]<label> [skip]
	if arg1.Code == token.ArgScope {

		label := ctx.NextArg().Name
		if label.Empty() {
			ctx.Error(ErrBadParam)
			return
		}

		skip, _ := ctx.GetArgInt(ctx.NextArg())

		switch arg1.Byte {
		case token.ScopePlayer:
			i := engine.Board.Player
			if i == -1 {
				return
			}
			obj := &engine.Board.Objects[i]
			if obj.Free {
				return
			}
			if pos != data.VectorMinus && obj.Position != pos {
				return
			}
			ctx.LastEvent = ObjectEvent(i, label, skip)
		case token.ScopeOthers:
			iter := engine.Board.ObjectsAll()
			if pos != data.VectorMinus { // apply [pos]
				iter = engine.Board.ObjectsAt(pos)
				// if object, send to tile as well
				if i := ctx.ObjectID(); i != -1 {
					tid := data.TileID(pos)
					ctx.LastEvent = TileEvent(tid, label, skip)
				}
			}
			if group := ctx.GetGroup(); !group.Empty() {
				for i, obj := range iter {
					if !obj.Free && obj.Group != group {
						ctx.LastEvent = ObjectEvent(i, label, skip) || ctx.LastEvent
					}
				}
			} else {
				for i, obj := range iter {
					if !obj.Free && i != ctx.ObjectID() {
						ctx.LastEvent = ObjectEvent(i, label, skip) || ctx.LastEvent
					}
				}
			}

		case token.ScopeGroup:
			iter := engine.Board.ObjectsAll()
			if pos != data.VectorMinus { // apply [pos]
				iter = engine.Board.ObjectsAt(pos)
			}
			if group := ctx.GetGroup(); !group.Empty() {
				for i, obj := range iter {
					if !obj.Free && obj.Group == group {
						ctx.LastEvent = ObjectEvent(i, label, skip) || ctx.LastEvent
					}
				}
			}
		case token.ScopeAll:
			iter := engine.Board.ObjectsAll()
			if pos != data.VectorMinus { // apply [pos]
				iter = engine.Board.ObjectsAt(pos)
				tid := data.TileID(pos)
				ctx.LastEvent = TileEvent(tid, label, skip)
			} else { // send to all tiles
				for i := range data.BoardLength {
					ctx.LastEvent = TileEvent(i, label, skip) || ctx.LastEvent
				}
			}
			for i, obj := range iter {
				if !obj.Free {
					ctx.LastEvent = ObjectEvent(i, label, skip) || ctx.LastEvent
				}
			}
		}
		return
	}

	// should only happen with [dir|pos] <label>
	if arg1.Code == token.ArgLabel {
		label := arg1.Name
		if label.Empty() {
			ctx.Error(ErrBadParam)
			return
		}
		skip, _ := ctx.GetArgInt(ctx.NextArg())
		if pos != data.VectorMinus {
			ctx.LastEvent = EventAt(pos, label, skip)
		}
	}

}

// <object>
func (ctx *ScriptContext) cmdFollow() {
	ctx.advance()
	arg := ctx.NextArg()
	// close previous follow
	engine.Board.Follow = -1
	// read and follow object, if found
	id := -1
	if arg.CanObject() {
		id = ctx.GetArgObject(arg)
	}
	if id != -1 {
		engine.Board.Follow = id
	}
}

func (ctx *ScriptContext) cmdGive() {
	ctx.advance()
	arg1 := ctx.NextArg()
	var counter *data.Counter

	switch arg1.Code {
	case token.ArgScope: // [scope] <counter> <amount>
		switch arg1.Byte {
		case token.ScopeBoard:
			arg2 := ctx.NextArg()
			if arg2.Code != token.ArgCounter {
				ctx.Error(ErrBadParam)
				return
			}
			counter = engine.Board.GetOrAddCounter(arg2.Name)
		case token.ScopeSelf:
			arg2 := ctx.NextArg()
			if arg2.Code != token.ArgCounter {
				ctx.Error(ErrBadParam)
				return
			}
			counter = ctx.GetOrAddCounter(arg2.Name)
		}
	case token.ArgCounter: // <counter> <amount>
		counter = SaveFile.GetOrAddCounter(arg1.Name)
	}

	if counter == nil { // tile event
		return
	}

	argAmount := ctx.NextArg()
	if v, ok := ctx.GetArgInt(argAmount); ok {
		counter.Value += v
	} else {
		ctx.Error(ErrBadParam)
		return
	}
}

// <pos>
func (ctx *ScriptContext) cmdGo() {
	ctx.advance()
	ctx.Stop = true

	argDir := ctx.NextArg()
	if argDir.Byte == token.DirIdle {
		return
	}

	pos, ok := ctx.GetArgPos(argDir)
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	if ctx.PlayerBorderCheck(pos) {
		return
	}

	if CollisionAt(pos, ctx.Collision(), ctx.ObjectID()) {
		// attempt to push
		if ctx.PushFlags().Has(data.CanPush) {
			pushDir := stepTo(ctx.Position().Sub(pos))
			if Push(pos, pushDir, ctx.Collision(), 0) {
				ctx.Move(pos)
				return
			}
		}

		ctx.rollback()
		return
	}

	ctx.Move(pos)
}

// <label> [skip]
func (ctx *ScriptContext) cmdGoto() {
	ctx.advance()

	argLabel := ctx.NextArg()
	if argLabel.Code != token.ArgLabel {
		ctx.Error(ErrBadParam)
		return
	}

	skip := 0
	skipArg := ctx.NextArg()
	if skipArg.CanInt() {
		skip, _ = ctx.GetArgInt(skipArg)
	}

	p := ctx.LabelPosition(argLabel.Name, skip)
	if p == -1 {
		return
	}

	ctx.CurrentPos = p
}

func (ctx *ScriptContext) cmdGroup() {
	ctx.advance()

	argName := ctx.NextArg()
	if argName.Code != token.ArgObjectName {
		ctx.Error(ErrBadParam)
		return
	}

	if id := ctx.ObjectID(); id != -1 {
		engine.Board.Objects[id].Group = argName.Name
	}
}

func (ctx *ScriptContext) cmdIdle() {
	ctx.advance()
	ctx.Stop = true
}

func (ctx *ScriptContext) cmdIf() {
	ctx.advance()

	arg1 := ctx.NextArg()
	pass, ok := ctx.GetArgCond(arg1)
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	ctx.LastResult = pass

	if pass {
		arg2 := ctx.NextArg()
		if arg2.Code == token.ArgCmd {
			ctx.CurrentPos = ctx.ArgPos
			ctx.runNext()
		} else {
			ctx.Error(ErrBadParam)
			return
		}
	}
}

func (ctx *ScriptContext) cmdLock() {
	ctx.advance()
	if id := ctx.ObjectID(); id != -1 {
		engine.Board.Objects[id].Lock = 1
	}
}

func (ctx *ScriptContext) cmdLog() {
	ctx.advance()
	text, ok := ctx.GetArgText(ctx.NextArg())
	if !ok {
		ctx.Error(ErrBadParam)
	}
	log.Println(string(text))
}

func (ctx *ScriptContext) cmdName() {
	ctx.advance()

	argName := ctx.NextArg()
	if argName.Code != token.ArgObjectName {
		ctx.Error(ErrBadParam)
		return
	}

	if id := ctx.ObjectID(); id != -1 {
		engine.Board.Objects[id].Name = argName.Name
	}
}

// <pos> [color] <element|prefab>
func (ctx *ScriptContext) cmdPut() {
	ctx.advance()
	color := -1

	argPos := ctx.NextArg()
	if !argPos.CanPos() {
		ctx.Error(ErrBadParam)
		return
	}
	pos, _ := ctx.GetArgPos(argPos)

	arg := ctx.NextArg()
	if arg.Code == token.ArgInteger {
		if v, ok := ctx.GetArgInt(arg); ok {
			color = v
		}
		arg = ctx.NextArg()
	}

	if arg.Code == token.ArgElement {
		element := arg.GetElement()
		// create tile
		tile := data.Tile{
			Cell: element.Cell,
			ID:   element.ID,
		}
		if color >= 0 {
			tile.Color = data.Color(color)
		}
		// set tile
		tid := data.TileID(pos)
		engine.Board.Tiles[tid] = tile
		return
	}
	if arg.Code == token.ArgPrefab {
		prefab := arg.GetPrefab()
		if CollisionAt(pos, prefab.Collision, -1) {
			return
		}
		// create object
		obj := prefab.Object()
		obj.Position = pos
		if color >= 0 {
			obj.Color = data.Color(color)
		}
		engine.Board.AddObject(obj, prefab.Player, prefab.Follow)
		return
	}

	ctx.Error(ErrMissingElementOrPrefab)
}

func (ctx *ScriptContext) cmdRestart() {
	ctx.CurrentPos = 0
}

// <label> [amount]
func (ctx *ScriptContext) cmdRestore() {
	ctx.advance()

	argLabel := ctx.NextArg()
	if argLabel.Code != token.ArgLabel {
		ctx.Error(ErrBadParam)
		return
	}

	amount := 1

	arg2 := ctx.NextArg()
	if v, ok := ctx.GetArgInt(arg2); ok {
		amount = v
	}

	if id := ctx.ObjectID(); id != -1 {
		obj := engine.Board.Objects[id]
		obj.DecZap(argLabel.Name, uint8(amount))
	}
}

func (ctx *ScriptContext) cmdScript() {
	ctx.advance()

	arg1 := ctx.NextArg()
	if !arg1.CanName() {
		ctx.Error(ErrBadParam)
		return
	}

	if script := engine.GetScript(arg1.Name); script != nil {
		if id := ctx.ObjectID(); id != -1 {
			engine.Board.Objects[id].Script = script
			engine.Board.Objects[id].ScriptPosition = 0
			engine.Board.Objects[id].ClearAllZaps()
		}
	}
}

func (ctx *ScriptContext) cmdSend() {
	ctx.advance()
	arg1 := ctx.NextArg()

	// <label> [skip] or SELF:<label> [skip]
	// does nothing for element scripts
	if arg1.Code == token.ArgLabel || arg1.Code == token.ArgScope && arg1.Byte == token.ScopeSelf {
		i := ctx.ObjectID()
		if i == -1 {
			return
		}

		label := data.Name{}
		if arg1.Code == token.ArgScope {
			label = ctx.NextArg().Name
			if label.Empty() {
				ctx.Error(ErrBadParam)
				return
			}
		} else {
			label = arg1.Name
		}

		skip, _ := ctx.GetArgInt(ctx.NextArg())
		p := SetObjectScriptPosition(i, label, skip)
		if !ctx.Event {
			ctx.CurrentPos = p
		}
		return
	}

	// [pos] option
	pos := data.VectorMinus
	if arg1.CanPos() {
		pos, _ = ctx.GetArgPos(arg1)
		arg1 = ctx.NextArg()
	}

	// [pos|dir] <label> [skip]
	/*if arg1.CanPos() {
		pos, _ := ctx.GetArgPos(arg1)

		label := ctx.NextArg().Name
		if label.Empty() {
			ctx.Error(ErrBadParam)
			return
		}

		skip, _ := ctx.GetArgInt(ctx.NextArg())

		for i, obj := range engine.Board.Objects {
			if !obj.Free && obj.Position == pos {
				p := SetObjectScriptPosition(i, label, skip)
				if !ctx.Event && ctx.ObjectID() == i {
					ctx.CurrentPos = p
				}
			}
		}
		return
	}*/

	// [name|group:]<label> [skip]
	if arg1.Code == token.ArgScopeName {
		label := ctx.NextArg().Name
		if label.Empty() {
			ctx.Error(ErrBadParam)
			return
		}

		skip, _ := ctx.GetArgInt(ctx.NextArg())

		iter := engine.Board.ObjectsAll()
		if pos != data.VectorMinus { // apply [pos]
			iter = engine.Board.ObjectsAt(pos)
		}

		for i, obj := range iter {
			if obj.Name == arg1.Name || obj.Group == arg1.Name {
				p := SetObjectScriptPosition(i, label, skip)
				if !ctx.Event && ctx.ObjectID() == i {
					ctx.CurrentPos = p
				}
			}
		}
		return
	}

	// [scope:]<label> [skip]
	if arg1.Code == token.ArgScope {

		label := ctx.NextArg().Name
		if label.Empty() {
			ctx.Error(ErrBadParam)
			return
		}

		skip, _ := ctx.GetArgInt(ctx.NextArg())

		switch arg1.Byte {
		case token.ScopePlayer:
			if i := engine.Board.Player; i != -1 {
				if obj := engine.Board.Objects[i]; !obj.Free {
					if pos != data.VectorMinus && pos != obj.Position {
						return
					}
					p := SetObjectScriptPosition(i, label, skip)
					if !ctx.Event && ctx.ObjectID() == i {
						ctx.CurrentPos = p
					}
				}
			}
		case token.ScopeOthers:
			iter := engine.Board.ObjectsAll()
			if pos != data.VectorMinus { // apply [pos]
				iter = engine.Board.ObjectsAt(pos)
			}
			for i := range iter {
				if i != ctx.ObjectID() {
					SetObjectScriptPosition(i, label, skip)
				}
			}
		case token.ScopeGroup:
			group := ctx.GetGroup()
			if !group.Empty() {
				iter := engine.Board.ObjectsAll()
				if pos != data.VectorMinus { // apply [pos]
					iter = engine.Board.ObjectsAt(pos)
				}
				for i, obj := range iter {
					if obj.Group == group {
						p := SetObjectScriptPosition(i, label, skip)
						if !ctx.Event && ctx.ObjectID() == i {
							ctx.CurrentPos = p
						}
					}
				}
			}
		case token.ScopeAll:
			iter := engine.Board.ObjectsAll()
			if pos != data.VectorMinus { // apply [pos]
				iter = engine.Board.ObjectsAt(pos)
			}
			for i := range iter {
				p := SetObjectScriptPosition(i, label, skip)
				if !ctx.Event && ctx.ObjectID() == i {
					ctx.CurrentPos = p
				}
			}
		}
		return
	}

	// <label> [skip]
	if arg1.Code == token.ArgLabel {
		skip, _ := ctx.GetArgInt(ctx.NextArg())
		iter := engine.Board.ObjectsAll()
		if pos != data.VectorMinus { // apply [pos]
			iter = engine.Board.ObjectsAt(pos)
		}
		for i := range iter {
			p := SetObjectScriptPosition(i, arg1.Name, skip)
			if !ctx.Event && ctx.ObjectID() == i {
				ctx.CurrentPos = p
			}
		}
	}

}

func (ctx *ScriptContext) cmdSet() {
	ctx.advance()
	arg1 := ctx.NextArg()
	var counter *data.Counter

	switch arg1.Code {
	case token.ArgFlag:
		ctx.SetFlag(arg1.Byte)
		return
	case token.ArgScope: // [scope] <counter> <amount>
		switch arg1.Byte {
		case token.ScopeBoard:
			argCounter := ctx.NextArg()
			if argCounter.Code != token.ArgCounter {
				ctx.Error(ErrBadParam)
				return
			}
			counter = engine.Board.GetOrAddCounter(argCounter.Name)
		case token.ScopeSelf:
			argCounter := ctx.NextArg()
			if argCounter.Code != token.ArgCounter {
				ctx.Error(ErrBadParam)
				return
			}
			counter = ctx.GetOrAddCounter(argCounter.Name)
		}
	case token.ArgCounter: // <counter> <amount>
		counter = SaveFile.GetOrAddCounter(arg1.Name)
	}

	if counter == nil { // tile event
		return
	}

	argAmount := ctx.NextArg()
	if argAmount.Code == token.ByteZero {
		counter.Value = 1
		return
	}
	if v, ok := ctx.GetArgInt(argAmount); ok {
		counter.Value = v
	} else {
		ctx.Error(ErrBadParam)
		return
	}
}

// <view>
func (ctx *ScriptContext) cmdView() {
	ctx.advance()
	arg := ctx.NextArg()
	if v, ok := ctx.GetArgInt(arg); ok {
		engine.Board.ViewID = v
		SetView(v)
	}
}

// <value> <cmd>
func (ctx *ScriptContext) cmdRepeat() {
	ctx.advance()
	val, ok := ctx.GetArgInt(ctx.NextArg())
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}
	ctx.SetRepeat(val)

	cmdArg := ctx.NextArg()
	if cmdArg.Code == token.ArgCmd {
		ctx.CurrentPos = ctx.ArgPos
		ctx.runNext()
	} else {
		ctx.Error(ErrBadParam)
		return
	}
}

// <dir> [color] <prefab>
func (ctx *ScriptContext) cmdShoot() {
	ctx.advance()
	ctx.Stop = true
	color := -1

	argDir := ctx.NextArg()
	if argDir.Byte == token.DirIdle {
		return
	}

	if !argDir.CanPos() {
		ctx.Error(ErrBadParam)
		return
	}
	dir, _ := ctx.GetArgDir(argDir)

	arg := ctx.NextArg()
	if arg.Code == token.ArgInteger {
		if v, ok := ctx.GetArgInt(arg); ok {
			color = v
		}
		arg = ctx.NextArg()
	}

	if arg.Code == token.ArgPrefab {
		prefab := arg.GetPrefab()
		pos := ctx.Position().Add(dir)
		/*if CollisionAt(pos, prefab.Collision, -1) {
			return
		}*/
		// create object
		obj := prefab.Object()
		obj.Position = pos
		obj.Step = dir
		if color >= 0 {
			obj.Color = data.Color(color)
		}
		engine.Board.AddObject(obj, prefab.Player, prefab.Follow)
		return
	}

	ctx.Error(ErrMissingElementOrPrefab)
}

func (ctx *ScriptContext) cmdTake() {
	ctx.advance()
	arg1 := ctx.NextArg()
	var counter *data.Counter

	switch arg1.Code {
	case token.ArgScope: // [scope] <counter> <amount>
		switch arg1.Byte {
		case token.ScopeBoard:
			argCounter := ctx.NextArg()
			if argCounter.Code != token.ArgCounter {
				ctx.Error(ErrBadParam)
				return
			}
			counter = engine.Board.GetCounter(argCounter.Name)
		case token.ScopeSelf:
			argCounter := ctx.NextArg()
			if argCounter.Code != token.ArgCounter {
				ctx.Error(ErrBadParam)
				return
			}
			counter = ctx.GetCounter(argCounter.Name)
		}
	case token.ArgCounter: // <counter> <amount>
		counter = SaveFile.GetCounter(arg1.Name)
	}

	argAmount := ctx.NextArg()
	v, ok := ctx.GetArgInt(argAmount)
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	if counter != nil && counter.Value >= v {
		counter.Value -= v
		return
	}

	// optional cmd if not enough
	arg2 := ctx.NextArg()
	if arg2.Code == token.ArgCmd {
		ctx.CurrentPos = ctx.ArgPos
		ctx.runNext()
	}
}

// <name>
func (ctx *ScriptContext) cmdTarget() {
	ctx.advance()
	arg := ctx.NextArg()
	if arg.CanName() {
		ctx.SetTarget(arg.Name)
	}
}

// <id> [title]
func (ctx *ScriptContext) cmdTextBox() {
	ctx.advance()

	id, ok := ctx.GetArgInt(ctx.NextArg())
	if !ok {
		log.Println(ErrBadParam)
		return
	}

	Text.textBoxID = id

	if title, ok := ctx.GetArgText(ctx.NextArg()); ok {
		Text.title = string(title)
	}
}

// <pos|dir> [color] <char>
func (ctx *ScriptContext) cmdTile() {
	ctx.advance()

	pos, ok := ctx.GetArgPos(ctx.NextArg())
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	if data.OutsideBoard(pos) {
		return
	}
	i := data.TileID(pos)

	arg1, ok := ctx.GetArgInt(ctx.NextArg())
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	arg2, ok := ctx.GetArgInt(ctx.NextArg())
	if ok {
		engine.Board.Tiles[i].Color = data.Color(arg1)
		engine.Board.Tiles[i].Char = byte(arg2)
	} else {
		engine.Board.Tiles[i].Char = byte(arg1)
	}
}

// <pos> [cmd]
func (ctx *ScriptContext) cmdTry() {
	ctx.advance()
	ctx.Stop = true

	argDir := ctx.NextArg()
	if argDir.Byte == token.DirIdle {
		return
	}

	pos, ok := ctx.GetArgPos(argDir)
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	if ctx.PlayerBorderCheck(pos) {
		return
	}

	if CollisionAt(pos, ctx.Collision(), ctx.ObjectID()) {

		// attempt to push
		if ctx.PushFlags().Has(data.CanPush) {
			pushDir := stepTo(ctx.Position().Sub(pos))
			if Push(pos, pushDir, ctx.Collision(), 0) {
				ctx.Move(pos)
				return
			}
		}

		// optional cmd
		arg2 := ctx.NextArg()
		if arg2.Code == token.ArgCmd {
			ctx.CurrentPos = ctx.ArgPos
			ctx.runNext()
		}
		return
	}

	ctx.Move(pos)
}

// <pos> [cmd]
func (ctx *ScriptContext) cmdTryShort() {
	ctx.advance()
	ctx.Stop = true

	argDir := ctx.NextArg()
	if argDir.Byte == token.DirIdle {
		return
	}

	pos, ok := ctx.GetArgPos(argDir)
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	if ctx.PlayerBorderCheck(pos) {
		return
	}

	if CollisionAt(pos, ctx.Collision(), ctx.ObjectID()) {
		return
	}

	ctx.Move(pos)
}

func (ctx *ScriptContext) cmdUnlock() {
	ctx.advance()
	if id := ctx.ObjectID(); id != -1 {
		engine.Board.Objects[id].Lock = 0
	}
}

// <dir>
func (ctx *ScriptContext) cmdWalk() {
	ctx.advance()
	arg1 := ctx.NextArg()
	if dir, ok := ctx.GetArgDir(arg1); ok {
		if id := ctx.ObjectID(); id != -1 {
			engine.Board.Objects[id].Step = dir
		}
	} else {
		ctx.Error(ErrBadParam)
		return
	}
}

// <pos|dir> <color> <text>
func (ctx *ScriptContext) cmdWrite() {
	ctx.advance()

	pos, ok := ctx.GetArgPos(ctx.NextArg())
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	if data.OutsideBoard(pos) {
		return
	}
	start := data.TileID(pos)

	argColor, ok := ctx.GetArgInt(ctx.NextArg())
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}
	color := data.Color(argColor)

	argText, ok := ctx.GetArgText(ctx.NextArg())
	if !ok {
		ctx.Error(ErrBadParam)
		return
	}

	for i := 0; i < len(argText); i++ {
		ti := i + start
		if ti >= engine.BufferLength {
			return
		}
		engine.Board.Tiles[ti].Color = color
		engine.Board.Tiles[ti].Char = argText[i]
	}
}

// <label> [amount]
func (ctx *ScriptContext) cmdZap() {
	ctx.advance()

	argLabel := ctx.NextArg()
	if argLabel.Code != token.ArgLabel {
		ctx.Error(ErrBadParam)
		return
	}

	amount := 1

	arg2 := ctx.NextArg()
	if v, ok := ctx.GetArgInt(arg2); ok {
		amount = v
	}

	if id := ctx.ObjectID(); id != -1 {
		obj := engine.Board.Objects[id]
		obj.IncZap(argLabel.Name, uint8(amount))
	}
}

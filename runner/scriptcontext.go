package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/token"
)

type ScriptContext struct {
	Script      *data.Script
	CurrentPos  int
	ArgPos      int
	CmdStartPos int
	CmdEndPos   int
	Event       bool
	Stop        bool
	Err         error
	LastResult  bool
	LastEvent   bool
	EntityContext
	CommandsRun int
}

type EntityContext interface {
	GetInteger(byte, data.Name) int
	Cell() *data.Cell
	Position() data.Vector
	Move(data.Vector)
	Collision() data.BitField
	PushFlags() data.BitField
	ObjectID() int
	SetScriptPosition(int)
	GetCounter(data.Name) *data.Counter
	GetOrAddCounter(data.Name) *data.Counter
	SetCounter(data.Name, int)
	Die()
	RunEvent(data.Name,int) bool
	LabelPosition(data.Name, int) int
	IsPlayer() bool
	SetRepeat(int)
	PopRepeat() bool
	SetFlag(byte)
	ClearFlag(byte)
	GetTarget() data.Name
	SetTarget(data.Name)
}

func (ctx *ScriptContext) Error(err error) {
	ctx.Err = err
	ctx.cmdEnd()
}

func (ctx *ScriptContext) Close() {
	Text.ShowText(ctx.ObjectID())
}

func (ctx *ScriptContext) GetInteger(scope byte, name data.Name) int {
	if scope == token.ScopeNone {
		switch name {
		case THISX:
			return ctx.Position().X
		case THISY:
			return ctx.Position().Y
		case STEPX:
			if obj := ctx.GetObject(); obj != nil {
				return obj.Step.X
			}
			return 0
		case STEPY:
			if obj := ctx.GetObject(); obj != nil {
				return obj.Step.Y
			}
			return 0
		case ENTERX:
			return engine.Board.Enter.X
		case ENTERY:
			return engine.Board.Enter.Y
		case COLOR:
			return int(ctx.Cell().Color)
		case FGCOLOR:
			return int(ctx.Cell().Color.FG())
		case BGCOLOR:
			return int(ctx.Cell().Color.BG())
		case CHAR:
			return int(ctx.Cell().Char)
		case PLAYERX:
			i := engine.Board.Player
			if i != -1 {
				return engine.Board.Objects[i].Position.X
			}
			return -1
		case PLAYERY:
			i := engine.Board.Player
			if i != -1 {
				return engine.Board.Objects[i].Position.Y
			}
			return -1
		default:
			if c := SaveFile.GetCounter(name); c != nil {
				return c.Value
			}
			return 0
		}
	}

	if scope == token.ScopeSelf {
		return ctx.EntityContext.GetInteger(scope, name)
	}

	if scope == token.ScopeBoard {
		if c := engine.Board.GetCounter(name); c != nil {
			return c.Value
		}
		return 0
	}

	return 0
}

func NewObjectContext(id int) *ScriptContext {
	object := &engine.Board.Objects[id]
	if object.Script == nil {
		return nil
	}
	return &ScriptContext{
		Script:     object.Script,
		CurrentPos: object.ScriptPosition,
		Event:      false,
		EntityContext: &ObjectContext{
			id:     id,
			object: object,
			event:  false,
		},
	}
}

func NewObjectEventContext(id int, event data.Name, skip int) *ScriptContext {
	if id == -1 {
		return nil
	}
	object := &engine.Board.Objects[id]
	if object.Script == nil {
		return nil
	}
	pos := object.LabelPosition(event, skip)
	if pos == -1 {
		return nil
	}
	return &ScriptContext{
		Script:     object.Script,
		CurrentPos: pos,
		Event:      true,
		EntityContext: &ObjectContext{
			id:     id,
			object: object,
			event:  true,
		},
	}
}

func (ctx *ScriptContext) GetObject() *data.Object {
	if i := ctx.ObjectID(); i != -1 {
		return &engine.Board.Objects[i]
	}
	return nil
}

func (ctx *ScriptContext) GetGroup() data.Name {
	if i := ctx.ObjectID(); i != -1 {
		return engine.Board.Objects[i].Group
	}
	return data.DefaultName
}

func (ctx *ScriptContext) RunEvent(event data.Name, skip int) bool {
	if obj := ctx.GetObject(); obj != nil {
		p := obj.ScriptPosition
		r := ctx.RunEvent(event, skip)
		if p != obj.ScriptPosition {
			ctx.CurrentPos = p
		}
		return r
	}
	return ctx.RunEvent(event, skip)
}

// ObjectContext

type ObjectContext struct {
	id     int
	object *data.Object
	event  bool
	repeat int
}

func (oc *ObjectContext) GetInteger(scope byte, name data.Name) int {
	if c := oc.object.GetCounter(name); c != nil {
		return c.Value
	}
	return 0
}

func (oc *ObjectContext) Cell() *data.Cell {
	return &oc.object.Cell
}

func (oc *ObjectContext) Position() data.Vector {
	return oc.object.Position
}

func (oc *ObjectContext) Move(value data.Vector) {
	oc.object.Position = value
	UpdateCamera(oc.id, engine.Board)
}

func (oc *ObjectContext) Collision() data.BitField {
	return oc.object.Collision
}

func (oc *ObjectContext) PushFlags() data.BitField {
	return oc.object.PushFlags
}

func (oc *ObjectContext) ObjectID() int {
	return oc.id
}

func (oc *ObjectContext) SetScriptPosition(value int) {
	oc.object.ScriptPosition = value
}

func (oc *ObjectContext) Die() {
	oc.object.Free = true
}

func (oc *ObjectContext) GetCounter(name data.Name) *data.Counter {
	return oc.object.GetCounter(name)
}

func (oc *ObjectContext) GetOrAddCounter(name data.Name) *data.Counter {
	return oc.object.GetOrAddCounter(name)
}

func (oc *ObjectContext) SetCounter(name data.Name, amount int) {
	oc.SetCounter(name, amount)
}

func (oc *ObjectContext) RunEvent(event data.Name, skip int) bool {
	result := ObjectEvent(oc.id, event, skip)
	return result
}

func (oc *ObjectContext) LabelPosition(label data.Name, skip int) int {
	return oc.object.LabelPosition(label, skip)
}

func (oc *ObjectContext) IsPlayer() bool {
	return oc.id == engine.Board.Player
}

func (oc *ObjectContext) SetRepeat(value int) {
	if oc.event {
		oc.repeat = value
	} else {
		oc.object.Repeat = value
	}
}

func (oc *ObjectContext) PopRepeat() bool {
	if oc.event {
		if oc.repeat == 0 {
			return false
		}
		oc.repeat--
		return true
	}
	if oc.object.Repeat == 0 {
		return false
	}
	oc.object.Repeat--
	return true
}

func (oc *ObjectContext) SetFlag(flag byte) {
	field, value := oc.flagData(flag)
	if field == nil {
		return
	}
	field.Set(value)
}

func (oc *ObjectContext) ClearFlag(flag byte) {
	field, value := oc.flagData(flag)
	if field == nil {
		return
	}
	field.Clear(value)
}

func (oc *ObjectContext) flagData(flag byte) (*data.BitField, data.BitField) {
	switch flag {
	case token.FlagCollision1:
		return &oc.object.Collision, data.Collision1
	case token.FlagCollision2:
		return &oc.object.Collision, data.Collision2
	case token.FlagCollision3:
		return &oc.object.Collision, data.Collision3
	case token.FlagCollision4:
		return &oc.object.Collision, data.Collision4
	case token.FlagCollision5:
		return &oc.object.Collision, data.Collision5
	case token.FlagCollision6:
		return &oc.object.Collision, data.Collision6
	case token.FlagCollision7:
		return &oc.object.Collision, data.Collision7
	case token.FlagCollision8:
		return &oc.object.Collision, data.Collision8

	case token.FlagCanPush:
		return &oc.object.PushFlags, data.CanPush
	case token.FlagPushableN:
		return &oc.object.PushFlags, data.PushableNorth
	case token.FlagPushableS:
		return &oc.object.PushFlags, data.PushableSouth
	case token.FlagPushableE:
		return &oc.object.PushFlags, data.PushableEast
	case token.FlagPushableW:
		return &oc.object.PushFlags, data.PushableWest
	case token.FlagPushableAll:
		return &oc.object.PushFlags, data.PushableAll
	case token.FlagPushableNS:
		return &oc.object.PushFlags, data.PushableNorthSouth
	case token.FlagPushableEW:
		return &oc.object.PushFlags, data.PushableEastWest

	case token.FlagPriority:
		return &oc.object.DrawFlags, data.DrawFlagPriority
	case token.FlagShowBG:
		return &oc.object.DrawFlags, data.DrawFlagShowBG

	default:
		return nil, 0
	}
}

func (oc *ObjectContext) GetTarget() data.Name {
	return oc.object.Target
}

func (oc *ObjectContext) SetTarget(name data.Name) {
	oc.object.Target = name
}


// Tile Context

func NewTileContext(id int, event data.Name, skip int) *ScriptContext {
	if id == -1 {
		return nil
	}
	tile := &engine.Board.Tiles[id]
	element := engine.ElementsFile.GetElement(tile.ID)
	if element == nil || element.Script == nil {
		return nil
	}
	pos := element.Script.LabelPosition(event, skip)
	if pos == -1 {
		return nil
	}
	return &ScriptContext{
		Script:     element.Script,
		CurrentPos: pos,
		Event:      true,
		EntityContext: &TileContext{
			id:      id,
			element: element,
			tile:    tile,
		},
	}
}

type TileContext struct {
	tile    *data.Tile
	element *data.Element
	id      int
	repeat  int
}

func (tc *TileContext) GetInteger(scope byte, name data.Name) int {
	return 0
}

func (tc *TileContext) Cell() *data.Cell {
	return &tc.tile.Cell
}

func (tc *TileContext) Position() data.Vector {
	return data.TilePos(tc.id)
}

func (tc *TileContext) Move(value data.Vector) {
	nId := data.TileID(value)
	engine.Board.Tiles[nId] = *tc.tile
	engine.Board.Tiles[tc.id] = data.Tile{}
	tc.id = nId
}

func (tc *TileContext) Collision() data.BitField {
	return tc.element.Collision
}

func (tc *TileContext) PushFlags() data.BitField {
	return tc.element.PushFlags
}

func (tc *TileContext) ObjectID() int {
	return -1
}

func (tc *TileContext) SetScriptPosition(value int) {
}

func (tc *TileContext) Die() {
	engine.Board.Tiles[tc.id] = data.Tile{}
}

func (tc *TileContext) GetCounter(name data.Name) *data.Counter {
	return nil
}

func (tc *TileContext) GetOrAddCounter(name data.Name) *data.Counter {
	return nil
}

func (tc *TileContext) SetCounter(name data.Name, amount int) {
}

func (tc *TileContext) RunEvent(event data.Name, skip int) bool {
	return TileEvent(tc.id, event, skip)
}

func (tc *TileContext) LabelPosition(label data.Name, skip int) int {
	return tc.element.Script.LabelPosition(label, skip)
}

func (tc *TileContext) IsPlayer() bool {
	return false
}

func (tc *TileContext) SetRepeat(value int) {
	tc.repeat = value
}

func (tc *TileContext) PopRepeat() bool {
	if tc.repeat == 0 {
		return false
	}
	tc.repeat--
	return true
}

func (tc *TileContext) SetFlag(flag byte) {}

func (tc *TileContext) ClearFlag(flag byte) {}

func (tc *TileContext) GetTarget() data.Name {
	return data.DefaultName
}

func (tc *TileContext) SetTarget(data.Name) {}
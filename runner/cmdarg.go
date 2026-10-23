package main

import (
	"encoding/binary"
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/token"
)

type CmdArg struct {
	Code byte
	Int  int16     // literal integer, color
	Byte byte      // scope, dir, command, length
	Name data.Name // label, counter, element, prefab, board, object, group
}

func (ctx *ScriptContext) NextArg() CmdArg {
	if ctx.ArgPos == ctx.CmdEndPos {
		return CmdArg{Code: token.ByteZero}
	}

	arg := CmdArg{}
	arg.Code = ctx.Script.Data[ctx.ArgPos]
	ctx.ArgPos++
	switch arg.Code {
	case token.ArgInteger:
		arg.Int = int16(binary.LittleEndian.Uint16(ctx.Script.Data[ctx.ArgPos : ctx.ArgPos+2]))
		ctx.ArgPos += 2
	case token.ArgScope, token.ArgCond, token.ArgMath, token.ArgDir, token.ArgPos, token.ArgObject, token.ArgFlag, token.ArgText:
		arg.Byte = ctx.Script.Data[ctx.ArgPos]
		ctx.ArgPos++
	case token.ArgElement, token.ArgPrefab, token.ArgName, token.ArgScopeName, token.ArgObjectName, token.ArgLabel, token.ArgCounter:
		id := int(binary.LittleEndian.Uint16(ctx.Script.Data[ctx.ArgPos : ctx.ArgPos+2]))
		ctx.ArgPos += 2
		if len(engine.NameTable.Names) <= id {
			ctx.Error(ErrBadToken)
			return arg
		}
		arg.Name = engine.NameTable.Names[id]
	case token.ArgTarget:
		arg.Name = ctx.GetTarget()
	}
	return arg
}

func (arg CmdArg) CanInt() bool {
	return arg.Code == token.ArgInteger || arg.Code == token.ArgScope || arg.Code == token.ArgCounter || arg.Code == token.ArgMath
}

func (ctx *ScriptContext) GetArgInt(arg CmdArg) (int, bool) {
	switch arg.Code {
	case token.ArgInteger:
		return int(arg.Int), true
	case token.ArgMath:
		return ctx.GetArgMath(arg)
	case token.ArgScope:
		arg2 := ctx.NextArg()
		return ctx.GetInteger(arg.Byte, arg2.Name), true
	case token.ArgCounter:
		return ctx.GetInteger(token.ScopeNone, arg.Name), true
	}
	return 0, false
}

func (arg CmdArg) CanName() bool {
	return arg.Code == token.ArgName || arg.Code == token.ArgTarget
}

func (arg CmdArg) CanObject() bool {
	return arg.Code == token.ArgObject || arg.Code == token.ArgObjectName
}

func (ctx *ScriptContext) GetArgObject(arg CmdArg) int {
	if arg.Code == token.ArgObjectName {
		this := ctx.ObjectID()
		pos := ctx.Position()
		dist := 0
		j := -1
		for i, obj := range engine.Board.Objects {
			if i != this && (obj.Name == arg.Name || obj.Group == arg.Name) {
				d := pos.Distance(obj.Position)
				if d == 0 || d == 1 {
					return i
				}
				if d != -1 && (j == -1 || d < dist) {
					dist = d
					j = i
				}
			}
		}
		return j
	}
	if arg.Code == token.ArgObject {
		switch arg.Byte {
		case token.ObjNone:
			return -1
		case token.ObjSelf:
			return ctx.ObjectID()
		case token.ObjPlayer:
			return engine.Board.Player
		}
	}
	return -1
}

func (arg CmdArg) GetElement() *data.Element {
	if arg.Code != token.ArgElement {
		return nil
	}
	return engine.ElementsFile.GetElementByName(arg.Name)
}

func (arg CmdArg) GetPrefab() *data.Prefab {
	if arg.Code != token.ArgPrefab {
		return nil
	}
	return engine.PrefabsFile.GetPrefabByName(arg.Name)
}

func (ctx *ScriptContext) NextName() data.Name {
	id := int(binary.LittleEndian.Uint16(ctx.Script.Data[ctx.ArgPos : ctx.ArgPos+2]))
	ctx.ArgPos += 2
	return engine.NameTable.Names[id]
}

func (ctx *ScriptContext) NextText() []byte {
	len := int(ctx.Script.Data[ctx.ArgPos])
	ctx.ArgPos++
	b := ctx.Script.Data[ctx.ArgPos : ctx.ArgPos+len]
	ctx.ArgPos += len
	return b
}

func (ctx *ScriptContext) GetArgCounter(arg CmdArg) (*data.Counter, bool) {
	if arg.Code == token.ArgCounter {
		return SaveFile.GetCounter(arg.Name), true
	}
	if arg.Code == token.ArgScope {
		argCounter := ctx.NextArg()
		if argCounter.Code != token.ArgCounter {
			return nil, false
		}
		switch arg.Byte {
		case token.ScopeSelf:
			return ctx.GetCounter(argCounter.Name), true
		case token.ScopeBoard:
			return engine.Board.GetCounter(argCounter.Name), true
		}
	}
	return nil, false
}

func (ctx *ScriptContext) GetArgText(arg CmdArg) ([]byte, bool) {
	if arg.Code != token.ArgText {
		return nil, false
	}
	len := int(arg.Byte)
	b := ctx.Script.Data[ctx.ArgPos : ctx.ArgPos+len]
	ctx.ArgPos += len
	return b, true
}

func (ctx *ScriptContext) GetArgElement(arg CmdArg) (*data.Element, bool) {
	if arg.Code != token.ArgElement || arg.Name.Empty() {
		return nil, false
	}
	element := engine.ElementsFile.GetElementByName(arg.Name)
	return element, element != nil
}

func (ctx *ScriptContext) NextColorElement() (*data.Color, *data.Element, bool) {
	var color *data.Color
	arg := ctx.NextArg()
	if arg.CanInt() {
		v, ok := ctx.GetArgInt(arg)
		if !ok {
			return nil, nil, false
		}
		c := data.Color(v)
		color = &c
		arg = ctx.NextArg()
	}
	element, ok := ctx.GetArgElement(arg)
	if !ok {
		return nil, nil, false
	}
	return color, element, true
}
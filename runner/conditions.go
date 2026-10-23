package main

import (
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/token"
)

func (arg CmdArg) CanCond() bool {
	return arg.Code == token.ArgCond
}

func (ctx *ScriptContext) GetArgCond(arg CmdArg) (bool, bool) {
	if arg.Code != token.ArgCond {
		return false, false
	}
	switch arg.Byte {
	case token.CondBlocked:
		arg2 := ctx.NextArg()
		p, ok := ctx.GetArgPos(arg2)
		if ok {
			return CollisionAt(p, ctx.Collision(), ctx.ObjectID()), true
		} else {
			return false, false
		}
	case token.CondContact:
		i := engine.Board.Player
		if i == -1 {
			return false, true
		}
		ppos := engine.Board.Objects[i].Position
		pos := ctx.Position()
		if pos == ppos || pos.Distance(ppos) == 1 {
			return true, true
		}
		return false, true
	case token.CondNot:
		arg2 := ctx.NextArg()
		if arg2.CanCond() {
			v, ok := ctx.GetArgCond(arg2)
			if !ok {
				return false, false
			}
			return !v, true
		}
	case token.CondAligned:
		i := engine.Board.Player
		if i == -1 {
			return false, true
		}
		ppos := engine.Board.Objects[i].Position
		pos := ctx.Position()
		if pos.X == ppos.X || pos.Y == ppos.Y {
			return true, true
		}
		return false, true
	case token.CondAny:
		color := -1
		arg2 := ctx.NextArg()
		if arg2.CanInt() {
			if v, ok := ctx.GetArgInt(arg2); ok {
				color = v
			}
			arg2 = ctx.NextArg()
		}
		if arg2.Code == token.ArgElement { // tiles
			element := engine.ElementsFile.GetElementByName(arg2.Name)
			if element == nil {
				return false, false
			}
			for i := range data.BoardLength {
				if engine.Board.Tiles[i].ID == element.ID {
					if color == -1 || engine.Board.Tiles[i].Color.FG() == data.Color(color) {
						return true, true
					}
				}
			}
		}
		if arg2.Code == token.ArgObjectName { // objects
			for _, obj := range engine.Board.Objects {
				if obj.Name == arg2.Name || obj.Group == arg2.Name {
					if color == -1 || obj.Cell.Color.FG() == data.Color(color) {
						return true, true
					}
				}
			}
		}
	case token.CondDetect: // pos, [color], element/object
		// get pos
		pos, ok := ctx.GetArgPos(ctx.NextArg())
		if !ok {
			ctx.Error(ErrBadParam)
			return false, false
		}

		if data.OutsideBoard(pos) {
			return false, true
		}
		// get [optional] color
		color := -1
		arg1 := ctx.NextArg()
		if arg1.CanInt() {
			if v, ok := ctx.GetArgInt(arg1); ok {
				color = v
			}
			arg1 = ctx.NextArg()
		}

		switch arg1.Code {
		case token.ArgElement:
			i := data.TileID(pos)
			tile := engine.Board.Tiles[i]
			element := engine.ElementsFile.GetElement(tile.ID)
			if element == nil {
				return false, false
			}
			if element.Name != arg1.Name {
				return false, true
			}
			if color == -1 {
				return true, true
			}
			if color >= 0 && color <= 15 && tile.Color.FG() == data.Color(color) {
				return true, true
			}
			return data.Color(color) == tile.Color, true
		case token.ArgObjectName:
			i := engine.Board.ObjectAt(pos)
			if i == -1 {
				return false, true
			}
			obj := engine.Board.Objects[i]
			if obj.Name != arg1.Name && obj.Group != arg1.Name {
				return false, true
			}
			if color == -1 {
				return true, true
			}
			if color >= 0 && color <= 15 && obj.Color.FG() == data.Color(color) {
				return true, true
			}
			return data.Color(color) == obj.Color, true
		case token.ArgObject:
			switch arg1.Byte {
			case token.ObjPlayer:
				if engine.Board.Player == -1 {
					return false, true
				}
				for i, obj := range engine.Board.Objects {
					if !obj.Free && obj.Position == pos && i == engine.Board.Player {
						return true, true
					}
				}
				return false, true
			case token.ObjGroup:
				group := ctx.GetGroup()
				if group.Empty() {
					return false, true
				}
				for _, obj := range engine.Board.Objects {
					if !obj.Free && obj.Position == pos && obj.Group == group {
						return true, true
					}
				}
				return false, true
			default:
				return false, false
			}
		default:
			return false, false
		}

	case token.CondIsTrue:
		arg1 := ctx.NextArg()
		value, ok := ctx.GetArgInt(arg1)
		if !ok {
			return false, false
		}
		return value > 0, true
	case token.CondEqual:
		arg1 := ctx.NextArg()
		value, ok := ctx.GetArgInt(arg1)
		if !ok {
			return false, false
		}
		arg2 := ctx.NextArg()
		if v, ok := ctx.GetArgInt(arg2); ok {
			return value == v, true
		}
	case token.CondNotEqual:
		arg1 := ctx.NextArg()
		value, ok := ctx.GetArgInt(arg1)
		if !ok {
			return false, false
		}
		arg2 := ctx.NextArg()
		if v, ok := ctx.GetArgInt(arg2); ok {
			return value != v, true
		}
	case token.CondGT:
		arg1 := ctx.NextArg()
		value, ok := ctx.GetArgInt(arg1)
		if !ok {
			return false, false
		}
		arg2 := ctx.NextArg()
		if v, ok := ctx.GetArgInt(arg2); ok {
			return value > v, true
		}
	case token.CondGTE:
		arg1 := ctx.NextArg()
		value, ok := ctx.GetArgInt(arg1)
		if !ok {
			return false, false
		}
		arg2 := ctx.NextArg()
		if v, ok := ctx.GetArgInt(arg2); ok {
			return value >= v, true
		}
	case token.CondLT:
		arg1 := ctx.NextArg()
		value, ok := ctx.GetArgInt(arg1)
		if !ok {
			return false, false
		}
		arg2 := ctx.NextArg()
		if v, ok := ctx.GetArgInt(arg2); ok {
			return value < v, true
		}
	case token.CondLTE:
		arg1 := ctx.NextArg()
		value, ok := ctx.GetArgInt(arg1)
		if !ok {
			return false, false
		}
		arg2 := ctx.NextArg()
		if v, ok := ctx.GetArgInt(arg2); ok {
			return value <= v, true
		}
	case token.CondBetween:
		arg1 := ctx.NextArg()
		value, ok := ctx.GetArgInt(arg1)
		if !ok {
			return false, false
		}
		arg2 := ctx.NextArg()
		if v, ok := ctx.GetArgInt(arg2); ok {
			arg3 := ctx.NextArg()
			if v2, ok := ctx.GetArgInt(arg3); ok {
				return value >= v && value <= v2, true
			}
		}
	case token.CondAccepted:
		return ctx.LastEvent, true
	default:
		return false, false
	}

	return false, false
}

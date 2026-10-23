package main

import (
	"silver-ink/data"
	"silver-ink/token"
	"silver-ink/engine"
)

func (arg CmdArg) CanPos() bool {
	return arg.Code == token.ArgPos || arg.Code == token.ArgDir
}

func (ctx *ScriptContext) GetArgPos(arg CmdArg) (data.Vector, bool) {
	if !arg.CanPos() {
		return data.VectorZero, false
	}
	pos, ok := ctx.getPos(arg)
	if !ok {
		return data.VectorZero, false
	}
	return pos, true
}

func (ctx *ScriptContext) getPos(arg CmdArg) (data.Vector, bool) {
	// convert Dir to Pos
	if arg.CanDir() {
		dir, ok := ctx.GetArgDir(arg)
		if !ok {
			return data.VectorZero, false
		}
		return ctx.Position().Add(dir), true
	}

	pos := data.Vector{}
	switch arg.Byte {
	case token.PosSelf:
		pos = ctx.Position()
	case token.PosPlayer:
		i := engine.Board.Player
		if i == -1 {
			return data.Vector{X:-1, Y:-1}, true
		}
		pos = engine.Board.Objects[i].Position
	case token.PosFollow:
		i := engine.Board.Player
		if i == -1 {
			return data.Vector{X:-1, Y:-1}, true
		}
		pos = engine.Board.Objects[i].Position
	case token.PosWhere:
		arg := ctx.NextArg()
		if arg.CanObject() {
			if id := ctx.GetArgObject(arg); id != -1 {
				pos = engine.Board.Objects[id].Position
			}
		}
	case token.PosAt: // AT <x> <y>
		argX := ctx.NextArg()
		if v, ok := ctx.GetArgInt(argX); ok {
			pos.X = v
		} else {
			return data.VectorZero, false
		}
		argY := ctx.NextArg()
		if v, ok := ctx.GetArgInt(argY); ok {
			pos.Y = v
		} else {
			return data.VectorZero, false
		}
	case token.PosBy: // BY <x> <y>
		pos = ctx.Position()
		argX := ctx.NextArg()
		if v, ok := ctx.GetArgInt(argX); ok {
			pos.X += v
		} else {
			return data.VectorZero, false
		}
		argY := ctx.NextArg()
		if !argY.CanInt() {
			return data.VectorZero, false
		}
		if v, ok := ctx.GetArgInt(argY); ok {
			pos.Y += v
		} else {
			return data.VectorZero, false
		}
	case token.PosTo:
		// starting pos
		arg2 := ctx.NextArg()
		if !arg2.CanPos() {
			return data.VectorZero, false
		}
		if p, ok := ctx.GetArgPos(arg2); ok {
			pos = p
		} else {
			return data.VectorZero, false
		}
		// add relative to first arg
		arg2 = ctx.NextArg()
		if arg2.CanDir() {
			if d, ok := ctx.GetArgDir(arg2); ok {
				pos.Add(d)
			} else {
				return data.VectorZero, false
			}
		}
		if arg2.CanPos() {
			if p, ok := ctx.GetArgPos(arg2); ok {
				pos.Add(p)
			} else {
				return data.VectorZero, false
			}
		}
	default:
		return data.Vector{}, false
	}

	return pos, true
}

func (arg CmdArg) CanDir() bool {
	return arg.Code == token.ArgDir 
}

func (ctx *ScriptContext) GetArgDir(arg CmdArg) (data.Vector, bool) {
	if !arg.CanDir() {
		return data.VectorZero, false
	}

	dir := data.Vector{}

	switch arg.Byte {
	case token.DirIdle:
		dir = data.VectorZero
	case token.DirNorth:
		dir = data.Vector{X: 0, Y: -1}
	case token.DirSouth:
		dir = data.Vector{X: 0, Y: 1}
	case token.DirEast:
		dir = data.Vector{X: 1, Y: 0}
	case token.DirWest:
		dir = data.Vector{X: -1, Y: 0}
	case token.DirCW:
		arg2 := ctx.NextArg()
		if d, ok := ctx.GetArgDir(arg2); ok {
			dir = data.Vector{X: -d.Y, Y: d.X}
		} else {
			return data.Vector{}, false
		}
	case token.DirCCW:
		arg2 := ctx.NextArg()
		if d, ok := ctx.GetArgDir(arg2); ok {
			dir = data.Vector{X: d.Y, Y: -d.X}
		} else {
			return data.Vector{}, false
		}
	case token.DirOPP:
		arg2 := ctx.NextArg()
		if d, ok := ctx.GetArgDir(arg2); ok {
			dir = data.Vector{X: -d.X, Y: -d.Y}
		} else {
			return data.Vector{}, false
		}
	case token.DirFlow:
		if obj := ctx.GetObject(); obj != nil {
			dir = obj.Step
		}
	case token.DirSeek:
		i := engine.Board.Player
		if i == -1 {
			return data.VectorZero, true
		}
		ppos := engine.Board.Objects[i].Position
		dir = stepTo(ctx.Position().Sub(ppos))
	case token.DirRnd:
		switch Rand.Int63n(4) {
		case 0:
			dir = data.Vector{X: 0, Y: -1}
		case 1:
			dir = data.Vector{X: 0, Y: 1}
		case 2:
			dir = data.Vector{X: 1, Y: 0}
		case 3:
			dir = data.Vector{X: -1, Y: 0}
		}
	case token.DirRndns:
		if RandBool() {
			dir = data.Vector{X: 0, Y: -1}
		}
		dir = data.Vector{X: 0, Y: 1}
	case token.DirRndne:
		if RandBool() {
			dir = data.Vector{X: 0, Y: -1}
		}
		dir = data.Vector{X: 1, Y: 0}
	case token.DirRnddp:
		arg2 := ctx.NextArg()
		if d, ok := ctx.GetArgDir(arg2); ok {
			if RandBool() {
				dir = data.Vector{X: -d.Y, Y: d.X}
			}
			dir = data.Vector{X: d.Y, Y: -d.X}
		} else {
			return data.Vector{}, false
		}
	default:
		return data.Vector{}, false
	}
	return dir, true
}

func stepTo(v data.Vector) data.Vector {
	ax := v.X
	if ax < 0 {
		ax = -ax
	}
	ay := v.Y
	if ay < 0 {
		ay = -ay
	}
	if ax >= ay {
		if v.X >= 0 {
			return data.Vector{X: -1, Y: 0}
		}
		return data.Vector{X: 1, Y: 0}
	}
	if v.Y >= 0 {
		return data.Vector{X: 0, Y: -1}
	}
	return data.Vector{X: 0, Y: 1}
}

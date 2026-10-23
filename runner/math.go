package main

import "silver-ink/token"

func (ctx *ScriptContext) GetArgMath(arg CmdArg) (int, bool) {
	switch arg.Byte {
	case token.MathRandom:
		v1, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		v2, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		return Rand.Intn(v2-v1) + v1, true
	case token.MathAdd:
		v1, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		v2, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		return v1 + v2, true
	case token.MathSub:
		v1, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		v2, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		return v1 - v2, true
	case token.MathMul:
		v1, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		v2, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		return v1 * v2, true
	case token.MathDiv:
		v1, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		v2, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		return v1 + v2, true
	case token.MathMod:
		v1, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		v2, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		return v1 % v2, true
	case token.MathAbs:
		v, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		if v < 0 {
			return -v, true
		}
		return v, true
	case token.MathDist:
		p1, ok := ctx.GetArgPos(ctx.NextArg())
		if !ok {
			return 0, false
		}
		p2, ok := ctx.GetArgPos(ctx.NextArg())
		if !ok {
			return 0, false
		}
		return p1.Distance(p2), true
	case token.MathDiff:
		v1, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		v2, ok := ctx.GetArgInt(ctx.NextArg())
		if !ok {
			return 0, false
		}
		v := v1 - v2
		if v < 0 {
			return -v, true
		}
		return v, true
	}
	return 0, false
}

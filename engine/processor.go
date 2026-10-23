package engine

import (
	"silver-ink/render"
	"github.com/jupiterrider/purego-sdl3/sdl"
)

var _stack []Processor
var _dirty bool
var _redraw bool

type Processor interface {
	Process()
	Active()
	Render(*render.Buffer, bool)
}

func PushStack(p Processor) {
	_stack = append(_stack, p)
	_dirty = true
	ClearFocus()
	ClearActionBuffer()
	ClearMouseState()
	activateTop()
}

func PopStack() {
	if len(_stack) > 0 {
		_stack[len(_stack)-1] = nil
		_stack = _stack[:len(_stack)-1]
		_dirty = true
		_redraw = true
		ClearFocus()
		ClearActionBuffer()
		ClearMouseState()
		activateTop()
	}
}

func Dirty() {
	_dirty = true
}

func activateTop() {
	if len(_stack) > 0 {
		_stack[len(_stack)-1].Active()
	}
}

func stackProcess() {
	if len(_stack) > 0 {
		_stack[len(_stack)-1].Process()
	}
}

func stackRender() {
	// redraw entire stack
	if _redraw {
		_redraw = false
		_dirty = true
		for i := 0; i < len(_stack)-1; i++ {
			_stack[i].Render(_buffer, false)
		}
	}

	// redraw top of the stack
	//if _dirty {
		_dirty = false
		if len(_stack) > 0 {
			_stack[len(_stack)-1].Render(_buffer, true)
		}
		renderBuffer()
	//}
}

func renderBuffer() {
	render.ClearBackground(_renderer)
	_buffer.Render()
	sdl.RenderPresent(_renderer)
}
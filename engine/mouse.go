package engine

import (
	"silver-ink/data"

	"github.com/Zyko0/go-sdl3/sdl"
)

type MouseState struct {
	Position           data.Vector
	Modifier           sdl.Keymod
	LeftButtonDown     bool
	RightButtonDown    bool
	LeftButtonClick    bool
	RightButtonClick   bool
	LeftButtonRelease  bool
	RightButtonRelease bool
}

var mouseStateQueue []MouseState = make([]MouseState, 0, 8)
var mouseLastPos data.Vector

func PullMouseState() []MouseState {
	msq := mouseStateQueue
	mouseStateQueue = mouseStateQueue[:0]
	return msq
}

func ClearMouseState() {
	mouseStateQueue = mouseStateQueue[:0]
}

func PeekMouseState() MouseState {
	button, _, _ := sdl.GetMouseState()
	return MouseState{
		Position:        mouseLastPos,
		LeftButtonDown:  button == sdl.BUTTON_LEFT,
		RightButtonDown: button == sdl.BUTTON_RIGHT,
	}
}

func pushMouseState(ms MouseState) {
	mouseStateQueue = append(mouseStateQueue, ms)
}

func mouseMouseMotion(motion *sdl.MouseMotionEvent) {
	p := data.Vector{X: int(motion.X) >> 3, Y: int(motion.Y) >> 3}
	if mouseLastPos != p {
		mouseLastPos = p
		button, _, _ := sdl.GetMouseState()
		ms := MouseState{
			Position:        p,
			Modifier:        sdl.GetModState(),
			LeftButtonDown:  button == sdl.BUTTON_LEFT,
			RightButtonDown: button == sdl.BUTTON_RIGHT,
		}
		pushMouseState(ms)
	}
}

func mouseMouseDown(button *sdl.MouseButtonEvent) {
	ms := MouseState{
		Position:         mouseLastPos,
		Modifier:         sdl.GetModState(),
		LeftButtonDown:   sdl.MouseButtonFlags(button.Button) == sdl.BUTTON_LEFT,
		RightButtonDown:  sdl.MouseButtonFlags(button.Button) == sdl.BUTTON_RIGHT,
		LeftButtonClick:  sdl.MouseButtonFlags(button.Button) == sdl.BUTTON_LEFT,
		RightButtonClick: sdl.MouseButtonFlags(button.Button) == sdl.BUTTON_RIGHT,
	}
	pushMouseState(ms)
}

func mouseMouseUp(button *sdl.MouseButtonEvent) {
	ms := MouseState{
		Position:           mouseLastPos,
		Modifier:           sdl.GetModState(),
		LeftButtonRelease:  sdl.MouseButtonFlags(button.Button) == sdl.BUTTON_LEFT,
		RightButtonRelease: sdl.MouseButtonFlags(button.Button) == sdl.BUTTON_RIGHT,
	}
	pushMouseState(ms)
}

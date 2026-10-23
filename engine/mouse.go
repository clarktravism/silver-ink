package engine

import (
	"silver-ink/data"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

type MouseState struct {
	Position data.Vector
	Modifier sdl.Keymod
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
	button := sdl.GetMouseState(nil, nil)
	return MouseState{
		Position:        mouseLastPos,
		LeftButtonDown:  button == sdl.ButtonLeft,
		RightButtonDown: button == sdl.ButtonRight,
	}
}

func pushMouseState(ms MouseState) {
	mouseStateQueue = append(mouseStateQueue, ms)
}

func mouseMouseMotion(motion sdl.MouseMotionEvent) {
	p := data.Vector{X: int(motion.X) >> 3, Y: int(motion.Y) >> 3}
	if mouseLastPos != p {
		mouseLastPos = p
		button := sdl.GetMouseState(nil, nil)
		ms := MouseState{
			Position:        p,
			Modifier:        sdl.GetModState(),
			LeftButtonDown:  button == sdl.ButtonLeft,
			RightButtonDown: button == sdl.ButtonRight,
		}
		pushMouseState(ms)
	}
}

func mouseMouseDown(button sdl.MouseButtonEvent) {
	ms := MouseState{
		Position:         mouseLastPos,
		Modifier:         sdl.GetModState(),
		LeftButtonDown:   sdl.MouseButtonFlags(button.Button) == sdl.ButtonLeft,
		RightButtonDown:  sdl.MouseButtonFlags(button.Button) == sdl.ButtonRight,
		LeftButtonClick:  sdl.MouseButtonFlags(button.Button) == sdl.ButtonLeft,
		RightButtonClick: sdl.MouseButtonFlags(button.Button) == sdl.ButtonRight,
	}
	pushMouseState(ms)
}

func mouseMouseUp(button sdl.MouseButtonEvent) {
	ms := MouseState{
		Position:           mouseLastPos,
		Modifier:           sdl.GetModState(),
		LeftButtonRelease:  sdl.MouseButtonFlags(button.Button) == sdl.ButtonLeft,
		RightButtonRelease: sdl.MouseButtonFlags(button.Button) == sdl.ButtonRight,
	}
	pushMouseState(ms)
}

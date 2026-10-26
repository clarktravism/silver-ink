package engine

import (
	"github.com/Zyko0/go-sdl3/sdl"
)

func poll() {
	var event sdl.Event
	for sdl.PollEvent(&event) {
		switch event.Type {
		case sdl.EVENT_QUIT:
			if quit_callback != nil {
				quit_callback()
			} else {
				Quit()
			}
		case sdl.EVENT_KEY_DOWN:
			key := event.KeyboardEvent()
			if FocusExists() {
				textKeyDown(key)
			} else {
				actionKeyCodeDown(key)
			}
		case sdl.EVENT_KEY_UP:
			key := event.KeyboardEvent()
			if !FocusExists() {
				actionKeyCodeUp(key)
			}
		case sdl.EVENT_TEXT_INPUT:
			text := event.TextInputEvent()
			if FocusExists() {
				textTextInput(text)
			}
		case sdl.EVENT_MOUSE_MOTION:
			_renderer.ConvertEventToRenderCoordinates(&event)
			motion := event.MouseMotionEvent()
			mouseMouseMotion(motion)
		case sdl.EVENT_MOUSE_BUTTON_DOWN:
			button := event.MouseButtonEvent()
			mouseMouseDown(button)
			actionMouseButton(button)
		case sdl.EVENT_MOUSE_BUTTON_UP:
			button := event.MouseButtonEvent()
			mouseMouseUp(button)
		case sdl.EVENT_MOUSE_WHEEL:
			wheel := event.MouseWheelEvent()
			actionMouseWheel(wheel)
		}
	}
}

var _quit = false
var quit_callback func()

func Quit() {
	_quit = true
}

func SetQuitCallback(callback func()) {
	quit_callback = callback
}

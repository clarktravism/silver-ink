package engine

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
)

func poll() {
	var event sdl.Event
	for sdl.PollEvent(&event) {
		switch event.Type() {
		case sdl.EventQuit:
			if quit_callback != nil {
				quit_callback()
			} else {
				Quit()
			}
		case sdl.EventKeyDown:
			key := event.Key()
			if FocusExists() {
				textKeyDown(key)
			} else {
				actionKeyCodeDown(key)
			}
		case sdl.EventKeyUp:
			if !FocusExists() {
				actionKeyCodeUp(event.Key())
			}
		case sdl.EventTextInput:
			if FocusExists() {
				textTextInput(event.Text())
			}
		case sdl.EventMouseMotion:
			sdl.ConvertEventToRenderCoordinates(_renderer, &event)
			mouseMouseMotion(event.Motion())
		case sdl.EventMouseButtonDown:
			button := event.Button()
			mouseMouseDown(button)
			actionMouseButton(button)
		case sdl.EventMouseButtonUp:
			mouseMouseUp(event.Button())
		case sdl.EventMouseWheel:
			actionMouseWheel(event.Wheel())
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
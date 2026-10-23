package engine

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
	"silver-ink/data"
)

var actionBuffer data.Name = data.DefaultName
var keycodeDown sdl.Keycode = 0
var actionRepeat int = 0

func GetAction() data.Name {
	if actionBuffer.Empty() {
		if actionRepeat == 0 {
			actionRepeat = 1
			return data.DefaultName
		}
		return currentKeyboardAction()
	}
	action := actionBuffer
	actionBuffer = data.DefaultName
	return action
}

func ClearActionBuffer() {
	actionBuffer = data.DefaultName
	keycodeDown = 0
	actionRepeat = 0
}

func actionKeyCodeDown(key sdl.KeyboardEvent) {
	if key.Key == sdl.KeycodeLShift {
		return
	}

	keycodeDown = key.Key
	actionRepeat = 0
	if actionBuffer.Empty() {
		actionBuffer = currentKeyboardAction()
	}
}

func actionKeyCodeUp(key sdl.KeyboardEvent) {
	if key.Key == sdl.KeycodeLShift {
		return
	}


	if keycodeDown == key.Key {
		keycodeDown = 0
		actionRepeat = 0
	}
}

func actionMouseWheel(wheel sdl.MouseWheelEvent) {
	if actionBuffer.Empty() {
		code := MouseCode(0)
		if wheel.IntegerY < 0 {
			code = MouseCodeWheelDown
		} else if wheel.IntegerY > 0 {
			code = MouseCodeWheelUp
		}
		if code != 0 {
			actionBuffer = getMouseAction(code)
		}
	}
}

func actionMouseButton(button sdl.MouseButtonEvent) {
	if button.Down && actionBuffer.Empty() {
		switch button.Button {
		case uint8(sdl.ButtonLeft):
			actionBuffer = getMouseAction(MouseCodeLeft)
		case uint8(sdl.ButtonRight):
			actionBuffer = getMouseAction(MouseCodeRight)
		case uint8(sdl.ButtonMiddle):
			actionBuffer = getMouseAction(MouseCodeMiddle)
		}
	}
}

const mod_filter = sdl.KeymodAlt | sdl.KeymodCtrl | sdl.KeymodShift

func currentKeyboardAction() data.Name {
	if keycodeDown == 0 {
		return data.DefaultName
	}
	modstate := sdl.GetModState() & mod_filter
	for _, kb := range kBinds {
		if kb.Keycode == keycodeDown && (kb.Modifier&modstate == modstate) {
			return kb.Action
		}
	}
	return data.DefaultName
}

func getMouseAction(code MouseCode) data.Name {
	modstate := sdl.GetModState()
	for _, mb := range mBinds {
		if mb.MouseCode == code && (mb.Modifier == 0 || mb.Modifier&modstate != 0) {
			return mb.Action
		}
	}
	return data.DefaultName
}

func IsShift(mod sdl.Keymod) bool {
	return mod&sdl.KeymodShift != 0
}

func IsCtrl(mod sdl.Keymod) bool {
	return mod&sdl.KeymodCtrl != 0
}

func IsAlt(mod sdl.Keymod) bool {
	return mod&sdl.KeymodAlt != 0
}
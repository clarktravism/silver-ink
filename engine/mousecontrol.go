package engine
/*
import (
	"silver-ink/data"
	"silver-ink/input"
)

// embed to control to allow them to respond to mouse events
type MouseControl struct {
	MouseRectangle data.Rect
	Click          Signal[input.MouseState]
	MouseMoved     Signal[input.MouseState]
	MouseEnter     Signal[input.MouseState]
	MouseLeave     Signal[input.MouseState]
	Hover          bool
	Processors     []MouseProcessor
}

type MouseProcessor interface {
	ProcessMouse(ms input.MouseState)
}

func (mc *MouseControl) AddMouseProcessor(processor MouseProcessor) {
	mc.Processors = append(mc.Processors, processor)
}

func (mc *MouseControl) ProcessMouseInput() {
	ms := input.PullMouseState()
	if ms.MouseMoved || ms.LeftButtonClick || ms.RightButtonClick {
		mc.ProcessMouse(ms)
	}
}

func (mc *MouseControl) ProcessMouse(ms input.MouseState) {
	inside := mc.MouseRectangle.IsInside(ms.Position)
	ms.Position = ms.Position.Sub(mc.MouseRectangle.Position)
	if inside {
		if !mc.Hover {
			mc.Hover = true
			mc.MouseEnter.Emit(ms)
		}
		if ms.MouseMoved {
			mc.MouseMoved.Emit(ms)
		}
		if ms.LeftButtonClick || ms.RightButtonClick {
			mc.Click.Emit(ms)
		}
		for _, mp := range mc.Processors {
			mp.ProcessMouse(ms)
		}
	} else if mc.Hover {
		mc.Hover = false
		mc.MouseLeave.Emit(ms)
	}
}
*/
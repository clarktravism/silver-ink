package engine

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
)

var _focused *TextInput
var _text_prop sdl.PropertiesID

func FocusExists() bool {
	return _focused != nil
}

func ClearFocus() {
	if _focused != nil {
		_focused.LoseFocus()
	}
}

func textTextInput(text sdl.TextInputEvent) {
	if _focused != nil {
		value := text.Text()
		_focused.appendText(value)
		if _focused.OnKeyPress != nil {
			_focused.OnKeyPress(value)
		}
	}
}

func textKeyDown(key sdl.KeyboardEvent) {
	if _focused != nil {
		_focused.keyCodePressed(key.Key)
	}
}

// embed to control to allow them to take focus
type TextInput struct {
	Focused bool
	Next    *TextInput
	Prev    *TextInput

	Text      string
	MaxLength int
	TextType  TextType

	OnKeyPress    func(string)
	OnBackspace   func()
	OnFocusLost   func()
	OnFocusGained func()
	OnTextUpdated func()
	OnCancel      func()
	OnAccept      func()
}

type TextType byte

const TextDefault TextType = 0
const TextCaps TextType = 1
const TextNumbers TextType = 2

func (ti *TextInput) TakeFocus() {
	if ti == nil {
		return
	}

	ClearFocus()
	_focused = ti
	ti.Focused = true

	ti.startText()
	if ti.OnFocusGained != nil {
		ti.OnFocusGained()
	} 
}

func (ti *TextInput) LoseFocus() {
	if ti == nil {
		return
	}
	if _focused == ti {
		_focused = nil
	}
	if ti.Focused {
		ti.Focused = false
		sdl.StopTextInput(_window)
		if ti.OnFocusLost != nil {
			ti.OnFocusLost()
		}
	}
}

func (fc *TextInput) FocusNext() {
	if fc == nil {
		return
	}
	fc.Next.TakeFocus()
}

func (fc *TextInput) FocusPrev() {
	if fc == nil {
		return
	}
	fc.Prev.TakeFocus()
}

func (tc *TextInput) appendText(text string) {
	if tc.MaxLength != 0 || len(tc.Text) < tc.MaxLength {
		tc.Text += text
		Dirty()
		if tc.OnTextUpdated != nil {
			tc.OnTextUpdated()
		}
	}
}

func (tc *TextInput) keyCodePressed(key sdl.Keycode) {
	switch key {
	case sdl.KeycodeBackspace:
		if len(tc.Text) > 0 {
			tc.Text = tc.Text[:len(tc.Text)-1]
		}
		Dirty()
		if tc.OnTextUpdated != nil {
			tc.OnTextUpdated()
		}
		if tc.OnBackspace != nil {
			tc.OnBackspace()
		}

	case sdl.KeycodeReturn:
		if tc.OnAccept != nil {
			tc.OnAccept()
		}
	case sdl.KeycodeTab:
		tc.FocusNext()
	case sdl.KeycodeEscape:
		if len(tc.Text) > 0 {
			tc.Text = ""
			if tc.OnTextUpdated != nil {
				tc.OnTextUpdated()
			} 
			Dirty()
		} else if tc.OnCancel != nil {
			tc.OnCancel()
		}
	}
}

func (tc *TextInput) startText() {
	if _text_prop == 0 {
		_text_prop = sdl.CreateProperties()
	}

	switch tc.TextType {
	case TextDefault:
		sdl.SetNumberProperty(_text_prop, sdl.PropTextInputTypeNumber, int64(sdl.TextInputTypeText))
		sdl.SetNumberProperty(_text_prop, sdl.PropTextInputCapitalizationNumber, int64(sdl.CapitalizeNone))
	case TextCaps:
		sdl.SetNumberProperty(_text_prop, sdl.PropTextInputTypeNumber, int64(sdl.TextInputTypeText))
		sdl.SetNumberProperty(_text_prop, sdl.PropTextInputCapitalizationNumber, int64(sdl.CapitalizeLetters))
	case TextNumbers:
		sdl.SetNumberProperty(_text_prop, sdl.PropTextInputTypeNumber, int64(sdl.TextInputTypeNumber))
		sdl.SetNumberProperty(_text_prop, sdl.PropTextInputCapitalizationNumber, int64(sdl.CapitalizeNone))
	}

	sdl.StartTextInputWithProperties(_window, _text_prop)
}

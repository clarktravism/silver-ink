package engine

import (
	"github.com/Zyko0/go-sdl3/sdl"
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

func textTextInput(text *sdl.TextInputEvent) {
	if _focused != nil {
		value := text.Text
		_focused.appendText(value)
		if _focused.OnKeyPress != nil {
			_focused.OnKeyPress(value)
		}
	}
}

func textKeyDown(key *sdl.KeyboardEvent) {
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
		_window.StopTextInput()
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
	case sdl.K_BACKSPACE:
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

	case sdl.K_RETURN:
		if tc.OnAccept != nil {
			tc.OnAccept()
		}
	case sdl.K_TAB:
		tc.FocusNext()
	case sdl.K_ESCAPE:
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
		_text_prop, _ = sdl.CreateProperties()
	}

	switch tc.TextType {
	case TextDefault:
		_text_prop.SetNumberProperty(PropTextInputTypeNumber, TextInputTypeText)
		//sdl.SetNumberProperty(_text_prop, PropTextInputTypeNumber, int64(sdl.TextInputTypeText))
		_text_prop.SetNumberProperty(PropTextInputCapitalizationNumber, CapitalizeNone)
		//sdl.SetNumberProperty(_text_prop, PropTextInputCapitalizationNumber, int64(sdl.CapitalizeNone))
	case TextCaps:
		_text_prop.SetNumberProperty(PropTextInputTypeNumber, TextInputTypeText)
		//sdl.SetNumberProperty(_text_prop, PropTextInputTypeNumber, int64(sdl.TextInputTypeText))
		_text_prop.SetNumberProperty(PropTextInputCapitalizationNumber, CapitalizeLetters)
		//sdl.SetNumberProperty(_text_prop, PropTextInputCapitalizationNumber, int64(sdl.CapitalizeLetters))
	case TextNumbers:
		_text_prop.SetNumberProperty(PropTextInputTypeNumber, TextInputTypeNumber)
		//sdl.SetNumberProperty(_text_prop, PropTextInputTypeNumber, int64(sdl.TextInputTypeNumber))
		_text_prop.SetNumberProperty(PropTextInputCapitalizationNumber, CapitalizeNone)
		//sdl.SetNumberProperty(_text_prop, PropTextInputCapitalizationNumber, int64(sdl.CapitalizeNone))
	}

	_window.StartTextInputWithProperties(_text_prop)
}

const (
	PropTextInputTypeNumber             = "SDL.textinput.type"
	PropTextInputCapitalizationNumber   = "SDL.textinput.capitalization"
	PropTextInputAutocorrectBoolean     = "SDL.textinput.autocorrect"
	PropTextInputMultilineBoolean       = "SDL.textinput.multiline"
	PropTextInputAndroidInputTypeNumber = "SDL.textinput.android.inputtype"
)

const (
	TextInputTypeText                  int64 = iota // The input is text.
	TextInputTypeTextName                           // The input is a person's name.
	TextInputTypeTextEmail                          // The input is an e-mail address.
	TextInputTypeTextUsername                       // The input is a username.
	TextInputTypeTextPasswordHidden                 // The input is a secure password that is hidden.
	TextInputTypeTextPasswordVisible                // The input is a secure password that is visible.
	TextInputTypeNumber                             // The input is a number.
	TextInputTypeNumberPasswordHidden               // The input is a secure PIN that is hidden.
	TextInputTypeNumberPasswordVisible              // The input is a secure PIN that is visible.
)

const (
	CapitalizeNone      int64 = iota // No auto-capitalization will be done.
	CapitalizeSentences              // The first letter of sentences will be capitalized.
	CapitalizeWords                  // The first letter of words will be capitalized.
	CapitalizeLetters                // All letters will be capitalized.
)

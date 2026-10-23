package engine

import (
	"fmt"
	"os"
	"strings"
	"encoding/json"
	"silver-ink/data"
	"github.com/jupiterrider/purego-sdl3/sdl"
)

var kBinds []KeyBind
var mBinds []MouseBind

type KeyBind struct {
	Keycode sdl.Keycode
	Modifier sdl.Keymod
	Action data.Name
}

type MouseBind struct {
	MouseCode MouseCode
	Modifier sdl.Keymod
	Action data.Name
}

type MouseCode byte
const (
	MouseCodeLeft = 1
	MouseCodeRight = 2
	MouseCodeMiddle = 3
	MouseCodeWheelUp = 4
	MouseCodeWheelDown = 5
)

type binding struct {
	Key string `json:"Key"`
	Action data.Name `json:"Action"`
}

func LoadKeyBinds(path string) error {
	kBinds = kBinds[:0]

	// read and decode
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	bindings := struct{
		Bindings []binding `json:"Bindings"` 
	}{}
	if err := json.Unmarshal(b, &bindings); err != nil {
		return fmt.Errorf("error while reading keybinds: %v", err)
	}

	for _, b := range bindings.Bindings {
		if err := parseBinding(b.Key, b.Action); err != nil {
			return fmt.Errorf("error while reading keybinds: %v", err)
		}
	}
	return nil
}

func parseBinding(key string, action data.Name) error {
	key = strings.ToLower(key)

	key, mod, err := parseKeyMod(key)
	if err != nil {
		return err
	}

	if strings.HasPrefix(key, "mouse") {
		if err := parseMouseBind(key, mod, action); err != nil {
			return err
		}
	} else {
		if err := parseKeyBind(key, mod, action); err != nil {
			return err
		}
	}
	return nil
}

func parseKeyMod(key string) (string, sdl.Keymod, error) {
	keymod := sdl.Keymod(0)
	modStr, keyStr, ok := strings.Cut(key, "+")
	if ok {
		key = keyStr
		switch modStr {
		case "shift":
			keymod = sdl.KeymodShift
		case "ctrl":
			keymod = sdl.KeymodCtrl
		case "alt":
			keymod = sdl.KeymodAlt
		default:
			return "", 0, fmt.Errorf("invalid modifer: %v", modStr)
		}
		return keyStr, keymod, nil
	}
	return key, keymod, nil
}

func parseMouseBind(key string, mod sdl.Keymod, action data.Name) error {
	mouseCode := MouseCode(0)
	switch key {
	case "mouse_left":
		mouseCode = MouseCodeLeft
	case "mouse_right":
		mouseCode = MouseCodeRight
	case "mouse_middle":
		mouseCode = MouseCodeMiddle
	case "mouse_wheel_up":
		mouseCode = MouseCodeWheelUp
	case "mouse_wheel_down":
		mouseCode = MouseCodeWheelDown
	default:
		return fmt.Errorf("invalid key: %v", key)
	}

	mBinds = append(mBinds, MouseBind{
		MouseCode: mouseCode,
		Modifier: mod,
		Action: action,
	})
	return nil
}

func parseKeyBind(key string, mod sdl.Keymod, action data.Name) error {
	keyCode, err := parseKeyCode(key)
	if err != nil {
		return err
	}
	kBinds = append(kBinds, KeyBind{
		Keycode: keyCode,
		Modifier: mod,
		Action: action,
	})
	return nil
}

func parseKeyCode(key string) (sdl.Keycode, error) {
	// parse scanscode
	switch key {
	case "space":
		return sdl.KeycodeSpace, nil
	case "enter":
		return sdl.KeycodeReturn, nil
	case "return":
		return sdl.KeycodeReturn, nil
	case "escape":
		return sdl.KeycodeEscape, nil
	case "up":
		return sdl.KeycodeUp, nil
	case "down":
		return sdl.KeycodeDown, nil
	case "left":
		return sdl.KeycodeLeft, nil
	case "right":
		return sdl.KeycodeRight, nil
	case "f1":
		return sdl.KeycodeF1, nil
	case "f2":
		return sdl.KeycodeF2, nil
	case "f3":
		return sdl.KeycodeF3, nil
	case "f4":
		return sdl.KeycodeF4, nil
	case "f5":
		return sdl.KeycodeF5, nil
	case "f6":
		return sdl.KeycodeF6, nil
	case "f7":
		return sdl.KeycodeF7, nil
	case "f8":
		return sdl.KeycodeF8, nil
	case "f9":
		return sdl.KeycodeF9, nil
	case "f10":
		return sdl.KeycodeF10, nil
	case "f11":
		return sdl.KeycodeF11, nil
	case "f12":
		return sdl.KeycodeF12, nil
	}

	if len(key) != 1 {
		return 0, fmt.Errorf("invalid key: %v", key)
	}

	k := sdl.Keycode(key[0])

	// most ascii characters
	if (k >= sdl.KeycodeExclaim && k <= sdl.KeycodeAt) || (k >= sdl.KeycodeLeftBracket && k <= sdl.KeycodeTilde) {
		return k, nil
	}

	return 0, fmt.Errorf("invalid key: %v", key)
}
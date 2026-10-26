package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"silver-ink/data"
	"strings"

	"github.com/Zyko0/go-sdl3/sdl"
)

var kBinds []KeyBind
var mBinds []MouseBind

type KeyBind struct {
	Keycode  sdl.Keycode
	Modifier sdl.Keymod
	Action   data.Name
}

type MouseBind struct {
	MouseCode MouseCode
	Modifier  sdl.Keymod
	Action    data.Name
}

type MouseCode byte

const (
	MouseCodeLeft      = 1
	MouseCodeRight     = 2
	MouseCodeMiddle    = 3
	MouseCodeWheelUp   = 4
	MouseCodeWheelDown = 5
)

type binding struct {
	Key    string    `json:"Key"`
	Action data.Name `json:"Action"`
}

func LoadKeyBinds(path string) error {
	kBinds = kBinds[:0]

	// read and decode
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	bindings := struct {
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
			keymod = sdl.KMOD_SHIFT
		case "ctrl":
			keymod = sdl.KMOD_CTRL
		case "alt":
			keymod = sdl.KMOD_ALT
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
		Modifier:  mod,
		Action:    action,
	})
	return nil
}

func parseKeyBind(key string, mod sdl.Keymod, action data.Name) error {
	keyCode, err := parseKeyCode(key)
	if err != nil {
		return err
	}
	kBinds = append(kBinds, KeyBind{
		Keycode:  keyCode,
		Modifier: mod,
		Action:   action,
	})
	return nil
}

func parseKeyCode(key string) (sdl.Keycode, error) {
	// parse scanscode
	switch key {
	case "space":
		return sdl.K_SPACE, nil
	case "enter":
		return sdl.K_RETURN, nil
	case "return":
		return sdl.K_RETURN, nil
	case "escape":
		return sdl.K_ESCAPE, nil
	case "up":
		return sdl.K_UP, nil
	case "down":
		return sdl.K_DOWN, nil
	case "left":
		return sdl.K_LEFT, nil
	case "right":
		return sdl.K_RIGHT, nil
	case "f1":
		return sdl.K_F1, nil
	case "f2":
		return sdl.K_F2, nil
	case "f3":
		return sdl.K_F3, nil
	case "f4":
		return sdl.K_F4, nil
	case "f5":
		return sdl.K_F5, nil
	case "f6":
		return sdl.K_F6, nil
	case "f7":
		return sdl.K_F7, nil
	case "f8":
		return sdl.K_F8, nil
	case "f9":
		return sdl.K_F9, nil
	case "f10":
		return sdl.K_F10, nil
	case "f11":
		return sdl.K_F11, nil
	case "f12":
		return sdl.K_F12, nil
	}

	if len(key) != 1 {
		return 0, fmt.Errorf("invalid key: %v", key)
	}

	k := sdl.Keycode(key[0])

	// most ascii characters
	if (k >= sdl.K_EXCLAIM && k <= sdl.K_AT) || (k >= sdl.K_LEFTBRACKET && k <= sdl.K_TILDE) {
		return k, nil
	}

	return 0, fmt.Errorf("invalid key: %v", key)
}

package engine

import (
	"silver-ink/render"
	"time"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
)

var _library library
var _window *sdl.Window
var _renderer *sdl.Renderer
var _buffer *render.Buffer
var _tpf time.Duration
var tick uint8

var BufferWidth int
var BufferHeight int
var BufferLength int

func Setup(title string, width int, height int, fps int) error {
	_library = binsdl.Load()

	BufferWidth = width
	BufferHeight = height
	BufferLength = width * height
	base_resolution_width := width * 8
	base_resolution_height := height * 8

	// sdl setup
	if err := sdl.SetHint(sdl.HINT_RENDER_VSYNC, "1"); err != nil {
		return err
	}
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	var err error
	_window, _renderer, err = sdl.CreateWindowAndRenderer(title, base_resolution_width*2, base_resolution_height*2, sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}

	_renderer.SetLogicalPresentation(int32(base_resolution_width), int32(base_resolution_height), sdl.LOGICAL_PRESENTATION_INTEGER_SCALE)

	// buffer setup
	_buffer = render.NewBuffer(width, height, _renderer)
	_tpf = time.Duration(time.Second / time.Duration(fps))

	return nil
}

func Close() {
	CloseAudio()
	render.DestroyCharsets()
	_window.Destroy()
	_renderer.Destroy()
	sdl.Quit()
	if _library != nil { 
		_library.Unload()
	}
}

func Run() {
	sdl.RunLoop(func() error {
		frameStart := time.Now()
		poll()
		if _quit {
			return sdl.EndLoop
		}

		stackProcess()
		stackRender()

		// cap fps
		frameDuration := time.Since(frameStart)
		if frameDuration < _tpf {
			sdl.DelayNS(uint64(_tpf - frameDuration))
		}
		tick++
		return nil
	})
}

func Title(text string) {
	_window.SetTitle(text)
}

func FPS(fps int) {
	if fps <= 0 {
		fps = 10
	}
	_tpf = time.Duration(time.Second / time.Duration(fps))
}

func CycleTick(cycle uint8) bool {
	return cycle != 0 && tick%cycle == 0
}


type library interface {
	Unload()
}
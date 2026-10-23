package engine

import (
	"log"
	"silver-ink/render"
	"time"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

var _window *sdl.Window
var _renderer *sdl.Renderer
var _buffer *render.Buffer
var _tpf time.Duration
var tick uint8

var BufferWidth int
var BufferHeight int
var BufferLength int

func Setup(title string, width int, height int, fps int) {
	BufferWidth = width
	BufferHeight = height
	BufferLength = width*height
	base_resolution_width := int32(width * 8)
	base_resolution_height := int32(height * 8)

	// sdl setup
	if !sdl.SetHint(sdl.HintRenderVSync, "1") {
		log.Println(sdl.GetError())
		panic(sdl.GetError())
	}

	if !sdl.Init(sdl.InitVideo) {
		log.Println(sdl.GetError())
		panic(sdl.GetError())
	}
	if !sdl.CreateWindowAndRenderer(title, base_resolution_width*2, base_resolution_height*2, sdl.WindowResizable, &_window, &_renderer) {
		log.Println(sdl.GetError())
		panic(sdl.GetError())
	}

	sdl.SetRenderLogicalPresentation(_renderer, base_resolution_width, base_resolution_height, sdl.LogicalPresentationIntegerScale)

	// buffer setup
	_buffer = render.NewBuffer(width, height, _renderer)

	_tpf = time.Duration(time.Second / time.Duration(fps))
}

func Close() {
	sdl.Quit()
	sdl.DestroyWindow(_window)
	sdl.DestroyRenderer(_renderer)
	render.DestroyCharsets()
}

func Run() {
	for {
		frameStart := time.Now()

		poll()
		if _quit {
			return
		}

		stackProcess()
		stackRender()

		// cap fps
		frameDuration := time.Since(frameStart)
		if frameDuration < _tpf {
			sdl.DelayNS(uint64(_tpf - frameDuration))
		}
		tick++
	}
}

func Title(text string) {
	sdl.SetWindowTitle(_window, text)
}

func FPS(fps int) {
	if fps <= 0 { fps = 10 }
	_tpf = time.Duration(time.Second / time.Duration(fps))
}

func CycleTick(cycle uint8) bool {
	return cycle != 0 && tick % cycle == 0
}
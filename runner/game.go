package main

import (
	"silver-ink/engine"
	"silver-ink/render"
)

type Game struct {}

func PushGame(boardName string) {
	g := Game{}
	StartNewGame(boardName)
	engine.PushStack(&g)
}

func (g *Game) Active() {
	engine.FPS(SaveFile.FPS)
}

func (g *Game) Process() {
	g.processBoard()
	Text.Process()
	engine.Dirty()
}

func (g *Game) Render(buffer *render.Buffer, top bool) {
	CameraRender(buffer)
	ViewRender(buffer)
	Text.Render(buffer)
}
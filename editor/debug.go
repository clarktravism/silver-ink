package main

import (
	"fmt"
	"log"
	"os/exec"
	"silver-ink/engine"
	"silver-ink/render"
)

func PushDebug() {
	engine.PushStack(DebugBoard{})
}

type DebugBoard struct {}

func (db DebugBoard) Render(buffer *render.Buffer, top bool) {}

func (db DebugBoard) Active() {
	go func() {
		path, err := exec.LookPath("./runner")
		if err != nil {
			log.Println(err)
			return
		}
		cmd := exec.Command(path, fmt.Sprintf("-w=%v", engine.WorldPath.Path), fmt.Sprintf("-b=%v", engine.Board.Name.String()))
		if err := cmd.Run(); err != nil {
			log.Println(err)
		}
		engine.PopStack()
	}()
}

func (db DebugBoard) Process() {}

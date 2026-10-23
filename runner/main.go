package main

import (
	"flag"
	"log"
	"path/filepath"
	"silver-ink/engine"
)

const buffer_width = 40
const buffer_height = 25
const buffer_length = buffer_width*buffer_height

func main() {
	// flag setup
	boardname := ""
	flag.StringVar(&boardname, "b", "", "board name")
	worldpath := ""
	flag.StringVar(&worldpath, "w", "", "world path")
	flag.Parse()
	if worldpath == "" {
		worldpath = settings.DefaultWorld
	}

	// log setup
	engine.SetLog(filepath.Join(worldpath, "runner.log"))
	defer engine.CloseLog()

	// catch and log any panics
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	engine.Setup("Silver-Ink", buffer_width, buffer_height, 10)
	defer engine.Close()
	engine.SetQuitCallback(onQuit)
	if err := StartWorld(worldpath); err != nil {
		log.Println(err)
		return
	}

	PushGame(boardname)
	engine.Run()
}

func onQuit() {
	engine.Quit()
}
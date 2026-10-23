package main

import (
	"log"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"silver-ink/data"
	"silver-ink/engine"
)

const world_file = "world.json"

type SelectWorld struct {
	Force bool
	Callback func(SelectWorldCallback)
	
	worlds []MetaWorld
}

type SelectWorldCallback struct {
	Accepted bool
	WorldName string
	WorldPath string
}

func (sw *SelectWorld) Push() {
	// grab worlds
	worlds, err := getWorldsList()
	if err != nil {
		log.Printf("error while reading worlds: %v\r\n", err)
	}
	sw.worlds = worlds

	// no worlds, open new world
	if len(worlds) == 0 {
		sw.pushNewWorld()
		return
	}

	sw.pushSelectWorld()
}

func (sw *SelectWorld) pushNewWorld() {
	tp := engine.TextPrompt{
		Title: "New World Name",
		Length: 20,
		AllowCancel: !sw.Force,
		Callback: sw.callbackNewWorld,
		Style: text_prompt_style,
	}
	tp.Push()
}

func (sw *SelectWorld) callbackNewWorld(callback engine.TextPromptCallback) {
	if !callback.Accepted {
		// restart to beginning
		sw.Push()
		return
	}

	path := filepath.Join(settings.WorldsPath, callback.Text)
	_, err := os.Stat(path)
	if err == nil {
		// path already exists, retry
		sw.pushNewWorld()
		return
	}
	if !errors.Is(err, os.ErrNotExist) {
		log.Printf("error while checking world file path: %v\r\n", err)
		return
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Printf("error while checking world file path: %v\r\n", err)
		return
	}

	if err := saveNewWorld(callback.Text, path); err != nil {
		log.Printf("error while creating world: %v\r\n", err)
		return
	}

	sw.Callback(SelectWorldCallback{
		Accepted: true,
		WorldName: callback.Text,
		WorldPath: path,
	})
}

const new_world_template = "templates/new_world"

func saveNewWorld(name string, worldpath string) error {

	// copy template files to new world path
	templateFs := os.DirFS(new_world_template)

	if err := os.CopyFS(worldpath, templateFs); err != nil {
		return err
	}
	engine.SetWorldPath(worldpath)

	// update world file with name
	engine.WorldFile.Name = name
	engine.WorldFile.TitleBoard = data.NewName("TITLE")
	engine.WorldFile.StartBoard = data.NewName("BOARD1")

	// title board
	titleBoard := data.NewBoard()
	titleBoard.Flags = 0
	titleBoard.ViewID = 2
	titleBoard.Name = engine.WorldFile.TitleBoard
	if err := engine.WorldPath.SaveBoard(&titleBoard); err != nil {
		return err
	}

	// start board
	startBoard := data.NewBoard()
	startBoard.Name = engine.WorldFile.StartBoard
	if err := engine.WorldPath.SaveBoard(&startBoard); err != nil {
		return err
	}

	if err := engine.WorldFile.Save(); err != nil {
		return err
	}

	return nil
}

func getWorldsList() ([]MetaWorld, error) {
	paths, err := os.ReadDir(settings.WorldsPath)
	if err != nil {
		return nil, err
	}
	worlds := make([]MetaWorld, 0)
	for _, path := range paths {
		if path.IsDir() {
			metaWorld, err := getMetaWorld(filepath.Join(settings.WorldsPath, path.Name()))
			if err == nil {
				worlds = append(worlds, metaWorld)
			}
		}
	}
	return worlds, nil
}

func getMetaWorld(worldpath string) (MetaWorld, error) {
	world := MetaWorld{
		Path: worldpath,
	}

	worldfilepath := filepath.Join(worldpath, world_file)

	b, err := os.ReadFile(worldfilepath)
	if err != nil {
		return MetaWorld{}, err
	}
	if err := json.Unmarshal(b, &world); err != nil {
		return MetaWorld{}, err
	}

	return world, nil
}


var NEW_WORLD = data.NewName("NEW_WORLD")
var OPEN_WORLD = data.NewName("OPEN_WORLD")

func (sw *SelectWorld) pushSelectWorld() {
	size := data.Vector{ X:38, Y: 25 }

	sb := engine.ScrollBox{
		Rectangle:  data.Rect{
			Position: engine.CenterPos(size),
			Size: size,
		},
		Title: "Select World",
		Style: scroll_box_style,
		Callback: sw.callbackSelectWorld,
		ForceSelect: sw.Force,
	}

	for _, world := range sw.worlds {
		sb.AppendLink(world.Name, OPEN_WORLD, world)
	}

	sb.AppendText("", false)
	sb.AppendLink("Create New World", NEW_WORLD, nil)

	sb.Push()
}

func (sw *SelectWorld) callbackSelectWorld(callback engine.ScrollBoxCallback) {
	
	if !callback.Accepted {
		sw.Callback(SelectWorldCallback{
			Accepted: false,
		})
	}

	if callback.Link == NEW_WORLD {
		sw.pushNewWorld()
		return
	}

	if callback.Link == OPEN_WORLD {
		world, _ := callback.Ext.(MetaWorld)
		
		sw.Callback(SelectWorldCallback{
			Accepted: true,
			WorldName: world.Name,
			WorldPath: world.Path,
		})
	}
}

// MetaWorld

type MetaWorld struct {
	Name string `json:"Name"`
	Path string
}
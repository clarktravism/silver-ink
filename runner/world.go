package main

import (
	"os"
	"silver-ink/engine"
)

const world_file = "world.json"

func StartWorld(path string) error {
	return loadWorld(path)
}

func inWorldFolder() bool {
	_, err := os.Stat(world_file)
	if err == nil {
		return true
	}
	return false
}

func loadWorld(path string) error {
	if err := engine.SetWorldPath(path); err != nil {
		return err
	}
	if err := engine.LoadKeyBinds(engine.WorldPath.Join("keybinds.json")); err != nil {
		return err
	}
	if err := engine.LoadScripts(); err != nil {
		return err
	}
	if err := LoadTextboxes(engine.WorldPath); err != nil {
		return err
	}
	return nil
}
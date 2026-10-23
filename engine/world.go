package engine

import (
	"os"
	"path/filepath"
	"encoding/json"
)

const world_file = "world.json"

func GetWorldsList(worldsDir string) ([]WorldInfo, error) {
	paths, err := os.ReadDir(worldsDir)
	if err != nil {
		return nil, err
	}
	worlds := make([]WorldInfo, 0)
	for _, path := range paths {
		if path.IsDir() {
			metaWorld, err := getMetaWorld(filepath.Join(worldsDir, path.Name()))
			if err == nil {
				worlds = append(worlds, metaWorld)
			}
		}
	}
	return worlds, nil
}

func getMetaWorld(worldpath string) (WorldInfo, error) {
	world := WorldInfo{
		Path: worldpath,
	}

	worldfilepath := filepath.Join(worldpath, world_file)

	b, err := os.ReadFile(worldfilepath)
	if err != nil {
		return WorldInfo{}, err
	}
	if err := json.Unmarshal(b, &world); err != nil {
		return WorldInfo{}, err
	}

	return world, nil
}

type WorldInfo struct {
	Name string `json:"Name"`
	Path string `json:"-"`
}
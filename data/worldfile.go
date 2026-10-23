package data

import (
	"encoding/json"
	"os"
)

const world_file = "world.json"

type WorldFile struct {
	Name       string `json:"Name"`
	TitleBoard Name   `json:"TitleBoard"`
	StartBoard Name   `json:"StartBoard"`
	FPS        int    `json:"FPS"`
	InputCycle uint8  `json:"InputCycle"`
	Counters
	path string
}

// world.json
func (wp WorldPath) LoadWorldFile() (*WorldFile, error) {
	world := WorldFile{path: wp.Join(world_file)}

	// read and decode
	b, err := os.ReadFile(world.path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &world); err != nil {
		return nil, err
	}

	return &world, nil
}

func (wf *WorldFile) Save() error {
	b, err := json.MarshalIndent(wf, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(wf.path, b, 0644); err != nil {
		return err
	}
	return nil
}

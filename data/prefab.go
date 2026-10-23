package data

import (
	"encoding/json"
	"os"
	"time"
)

type Prefab struct {
	ID    int   `json:"ID"`
	Name  Name  `json:"Name"`
	Group Name  `json:"Group"`
	Cycle uint8 `json:"Cycle"`
	Cell
	Target    Name     `json:"Target"`
	DrawFlags BitField `json:"DrawFlags"`
	Collision BitField `json:"Collision"`
	PushFlags BitField `json:"PushFlags"`
	Follow    bool     `json:"Follow"`
	Player    bool     `json:"Player"`
	Counters
	ScriptName Name `json:"Script"`
	Script     *Script
}

func (prefab *Prefab) Object() Object {
	return Object{
		Name:       prefab.Name,
		Group:      prefab.Group,
		Cycle:      prefab.Cycle,
		Cell:       prefab.Cell,
		Target:     prefab.Target,
		DrawFlags:  prefab.DrawFlags,
		Collision:  prefab.Collision,
		PushFlags:  prefab.PushFlags,
		Counters:   prefab.Counters.Copy(), // need deep copy
		ScriptName: prefab.ScriptName,
		Script:     prefab.Script,
	}
}

// prefabs file

const prefabs_file = "prefabs.json"

type PrefabsFile struct {
	Prefabs []Prefab
	names   map[Name]*Prefab
	modtime time.Time
	path    string
}

func (pf *PrefabsFile) GetPrefabByName(name Name) *Prefab {
	if pf == nil {
		return nil
	}
	return pf.names[name]
}

func (wp WorldPath) LoadPrefabFile() (*PrefabsFile, error) {
	pf := PrefabsFile{
		path: wp.Join(prefabs_file),
	}
	if err := pf.load(); err != nil {
		return nil, err
	}
	return &pf, nil
}

func (pf *PrefabsFile) HasChanged() bool {
	fi, err := os.Stat(pf.path)
	if err != nil {
		return false // unknown
	}
	return !pf.modtime.Equal(fi.ModTime())
}

func (pf *PrefabsFile) Reload() error {
	return pf.load()
}

func (pf *PrefabsFile) load() error {
	// get file info
	fi, err := os.Stat(pf.path)
	if err != nil {
		return err
	}
	pf.modtime = fi.ModTime()

	// read and decode
	b, err := os.ReadFile(pf.path)
	if err != nil {
		return err
	}

	v := struct{
		Prefabs []Prefab `json:"Prefabs"`
	}{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	pf.Prefabs = v.Prefabs

	// map prefabs
	pf.names = make(map[Name]*Prefab, len(pf.Prefabs))
	for i := range pf.Prefabs {
		pf.names[pf.Prefabs[i].Name] = &pf.Prefabs[i]
	}

	return nil
}

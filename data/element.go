package data

import (
	"encoding/json"
	"os"
	"time"
)

const elements_file = "elements.json"

type Element struct {
	ID        byte     `json:"ID"`
	Name      Name     `json:"Name"`
	Collision BitField `json:"Collision"`
	PushFlags BitField `json:"PushFlags"`
	Cell
	ScriptName Name `json:"Script"`
	Script     *Script
}

type ElementsFile struct {
	Elements []Element `json:"Elements"`
	names    map[Name]*Element
	modtime  time.Time
	path     string
}

func (ef *ElementsFile) GetElementByName(name Name) *Element {
	if ef == nil {
		return nil
	}
	return ef.names[name]
}

func (ef *ElementsFile) GetElement(id byte) *Element {
	if ef == nil {
		return nil
	}
	if int(id) >= len(ef.Elements) {
		return nil
	}
	return &ef.Elements[id]
}

func (wp WorldPath) LoadElementsFile() (*ElementsFile, error) {
	ef := ElementsFile{
		path: wp.Join(elements_file),
	}
	if err := ef.load(); err != nil {
		return nil, err
	}
	return &ef, nil
}

func (ef *ElementsFile) HasChanged() bool {
	fi, err := os.Stat(ef.path)
	if err != nil {
		return false // unknown
	}
	return !ef.modtime.Equal(fi.ModTime())
}

func (ef *ElementsFile) Reload() error {
	return ef.load()
}

func (ef *ElementsFile) load() error {

	// get file info
	fi, err := os.Stat(ef.path)
	if err != nil {
		return err
	}
	ef.modtime = fi.ModTime()

	// read and decode
	b, err := os.ReadFile(ef.path)
	if err != nil {
		return err
	}
	
	v := struct{
		Elements []Element `json:"Elements"`
	}{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	// get max id, load scripts
	max := byte(0)
	for i := range v.Elements {
		if v.Elements[i].ID > max {
			max = v.Elements[i].ID
		}
	}

	// set and map elements
	ef.names = make(map[Name]*Element, len(v.Elements))
	ef.Elements = make([]Element, max+1)
	for i := range v.Elements {
		// set to slice
		id := v.Elements[i].ID
		ef.Elements[id] = v.Elements[i]
		// set to map
		ef.names[ef.Elements[id].Name] = &ef.Elements[id]
	}

	return nil
}

package engine

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"silver-ink/data"
	"silver-ink/token"
	"strings"
)

const scripts_folder = "scripts"
const script_ext = ".obj"

var scripts map[data.Name]*data.Script

var ErrWorldPathNotLoaded = errors.New("world path not loaded")

func LoadScripts() error {
	if WorldPath == nil {
		return ErrWorldPathNotLoaded
	}

	scripts = make(map[data.Name]*data.Script, 32)
	paths, err := os.ReadDir(WorldPath.Join(scripts_folder))
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(make([]byte, 0, 1024))

	for _, path := range paths {
		if fi, _ := path.Info(); !fi.IsDir() && strings.EqualFold(filepath.Ext(fi.Name()), script_ext) {
			name, _, _ := strings.Cut(fi.Name(), ".")
			filename := WorldPath.Join(scripts_folder, fi.Name())
			if err := loadScript(data.NewName(name), filename, buffer); err != nil {
				return fmt.Errorf("%v %v", fi.Name(), err.Error())
			}
		}
	}

	if err := loadScriptsToElements(); err != nil {
		return err
	}
	if err := loadScriptsToPrefabs(); err != nil {
		return err
	}

	return nil
}

func GetScript(name data.Name) *data.Script {
	if script, ok := scripts[name]; ok {
		return script
	}
	return nil
}

func loadScript(name data.Name, filename string, buffer *bytes.Buffer) error {
	buffer.Reset()

	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	// read and scan file
	labeler := labelStore{}
	scanner := token.NewByteCodeScanner(f, &labeler, &NameTable)
	for scanner.Scan() {
		if _, err := buffer.Write(scanner.Bytes()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	code := make([]byte, buffer.Len())
	copy(code, buffer.Bytes())

	scripts[name] = &data.Script{
		Name:   name,
		Data:   code,
		Labels: labeler.labels,
	}
	return nil
}

type labelStore struct {
	labels map[data.Name][]int
}

func (ls *labelStore) AddLabel(name data.Name, pos int) {
	if ls.labels == nil {
		ls.labels = make(map[data.Name][]int)
	}
	ls.labels[name] = append(ls.labels[name], pos)
}

type nameStore struct {
	names []data.Name
}

func (ns *nameStore) AddName(name data.Name) int {
	for i := range ns.names {
		if ns.names[i] == name {
			return i
		}
	}
	ns.names = append(ns.names, name)
	return len(ns.names) - 1
}

func LoadScriptsToBoard(board *data.Board) error {
	for i, obj := range board.Objects {
		if !obj.ScriptName.Empty() {
			script := GetScript(obj.ScriptName)
			if script == nil {
				return fmt.Errorf("script not found: %v", obj.ScriptName.String())
			}
			board.Objects[i].Script = script
		}
	}
	return nil
}

func loadScriptsToElements() error {
	for i, ele := range ElementsFile.Elements {
		if !ele.ScriptName.Empty() {
			script := GetScript(ele.ScriptName)
			if script == nil {
				return fmt.Errorf("script not found: %v", ele.ScriptName.String())
			}
			ElementsFile.Elements[i].Script = script
		}

	}
	return nil
}

func loadScriptsToPrefabs() error {
	for i, pf := range PrefabsFile.Prefabs {
		if !pf.ScriptName.Empty() {
			script := GetScript(pf.ScriptName)
			if script == nil {
				return fmt.Errorf("script not found: %v", pf.ScriptName.String())
			}
			PrefabsFile.Prefabs[i].Script = script
		}
	}
	return nil
}

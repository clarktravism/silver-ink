package data
/*
import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AssetFile struct {
	Assets  []AssetItem `json:"Assets"`
	ids     map[int16]int
	names   map[Name]int
	modtime time.Time
	path    string // path the json file
	folder  string // folder of where the assets are located
	ext     string // ext of assets
}

type AssetItem struct {
	ID   int16  `json:"ID"`
	Name Name   `json:"Name"`
	Path string `json:"Path"`
}

func (wp WorldPath) LoadAssetFile(path string, folder string, ext string) (*AssetFile, error) {
	af := AssetFile{
		path:   wp.Join(path),
		folder: wp.Join(folder),
		ext:    ext,
	}
	if err := af.load(); err != nil {
		return nil, err
	}
	af.Scan()
	return &af, nil
}

func (af *AssetFile) HasChanged() bool {
	fi, err := os.Stat(af.path)
	if err != nil {
		return false // unknown
	}
	return !af.modtime.Equal(fi.ModTime())
}

func (af *AssetFile) Reload() error {
	return af.load()
}

func (af *AssetFile) load() error {
	// get file info
	fi, err := os.Stat(af.path)
	if err != nil {
		return err
	}
	af.modtime = fi.ModTime()

	// read and decode
	b, err := os.ReadFile(af.path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, af); err != nil {
		return err
	}

	// set maps
	af.names = make(map[Name]int, len(af.Assets))
	af.ids = make(map[int16]int, len(af.Assets))
	for i, asset := range af.Assets {
		if asset.Name.Empty() || asset.ID <= 0 || asset.Path == "" {
			continue
		}
		af.names[asset.Name] = i
		af.ids[asset.ID] = i
	}
	return nil
}

func (af *AssetFile) newID() int16 {
	id := int16(1)
	for i := range af.Assets {
		if af.Assets[i].ID >= id {
			id = af.Assets[i].ID + 1
		}
	}
	return id
}

func (af *AssetFile) pathName(name Name) string {
	if name.Empty() {
		return ""
	}
	return strings.ToLower(name.String() + af.ext)
}

// scans folder to new, renamed, deleted assets
func (af *AssetFile) Scan() {

	found := make([]bool, len(af.Assets))
	entries, _ := os.ReadDir(af.folder)
	for _, entry := range entries {
		// filter out files without matching ext
		if entry.IsDir() {
			continue
		}
		info, _ := entry.Info()
		path := strings.ToLower(info.Name())
		if filepath.Ext(path) != af.ext {
			continue
		}
		// find matching path
		if i, ok := af.findPath(path); ok {
			found[i] = true
		} else {
			af.AddPath(path)
		}
	}

	// remove missing files
	for i := range found {
		if !found[i] {
			af.remove(i)
		}
	}

	af.Save()
}

func (af *AssetFile) FindID(id int16) (AssetItem, bool) {
	if i, ok := af.ids[id]; ok {
		return af.Assets[i], true
	}
	return AssetItem{}, false
}

func (af *AssetFile) FindName(name Name) (AssetItem, bool) {
	if af == nil { return AssetItem{}, false }
	if i, ok := af.names[name]; ok {
		return af.Assets[i], true
	}
	return AssetItem{}, false
}

func (af *AssetFile) FindPath(path string) (AssetItem, bool) {
	if i, ok := af.findPath(path); ok {
		return af.Assets[i], true
	}
	return AssetItem{}, false
}

func (af *AssetFile) findPath(path string) (int, bool) {
	for i := range af.Assets {
		if strings.EqualFold(af.Assets[i].Path, path) {
			return i, true
		}
	}
	return 0, false
}

func (af *AssetFile) AddPath(path string) AssetItem {
	name, _ := strings.CutSuffix(path, af.ext)
	asset := AssetItem{
		Path: path,
		Name: NewName(name),
	}
	asset.ID = af.newID()
	af.Assets = append(af.Assets, asset)
	i := len(af.Assets) - 1
	af.names[asset.Name] = i
	af.ids[asset.ID] = i
	return asset
}

func (af *AssetFile) remove(i int) {
	asset := af.Assets[i]

	// remove from map
	delete(af.ids, asset.ID)
	delete(af.names, asset.Name)

	// cut out of assets slice
	af.Assets[i] = af.Assets[len(af.Assets)-1]
	af.Assets = af.Assets[:len(af.Assets)-1]

	// update maps
	if len(af.Assets) > i {
		af.ids[af.Assets[i].ID] = i
		af.names[af.Assets[i].Name] = i
	}
}

func (af *AssetFile) Rename(old Name, new Name) error {
	i, ok := af.names[old]
	if !ok {
		return nil
	}
	delete(af.names, old)
	af.names[new] = i

	asset := af.Assets[i]
	oldPath := filepath.Join(af.folder, asset.Path)
	asset.Name = new
	asset.Path = af.pathName(new)
	af.Assets[i] = asset
	newPath := filepath.Join(af.folder, asset.Path)

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	return af.Save()
}

func (af *AssetFile) PathToID(id int16) string {
	if asset, ok := af.FindID(id); ok {
		return filepath.Join(af.folder, asset.Path)
	}
	return ""
}

func (af *AssetFile) PathToName(name Name) string {
	if asset, ok := af.FindName(name); ok {
		return filepath.Join(af.folder, asset.Path)
	}
	return ""
}

func (af *AssetFile) FullPath(asset AssetItem) string {
	return filepath.Join(af.folder, asset.Path)
}

func (af *AssetFile) Save() error {
	b, err := json.MarshalIndent(af, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(af.path, b, 0644); err != nil {
		return err
	}
	return nil
}
*/
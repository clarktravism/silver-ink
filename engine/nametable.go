package engine

import (
	"silver-ink/data"
)

var NameTable nameTable

// implements token.NameRef
type nameTable struct {
	Names []data.Name
}

func (nt *nameTable) IsElement(name data.Name) bool {
	if e := ElementsFile.GetElementByName(name); e != nil {
		return true
	}
	return false
}

func (nt *nameTable) IsPrefab(name data.Name) bool {
	if p := PrefabsFile.GetPrefabByName(name); p != nil {
		return true
	}
	return false
}

func (nt *nameTable) GetName(name data.Name) int {
	for i := range nt.Names {
		if nt.Names[i] == name {
			return i
		}
	}
	nt.Names = append(nt.Names, name)
	return len(nt.Names)-1
}

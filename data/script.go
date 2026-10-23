package data

const scripts_folder = "scripts"
const script_ext = ".obj"

type Script struct {
	Name   Name
	Data   []byte
	Labels map[Name][]int
	refs int
}

func (scr *Script) LabelPosition(label Name, skip int) int {
	if scr == nil || scr.Labels == nil || skip < 0 {
		return -1
	}
	v, ok := scr.Labels[label]
	if !ok {
		return -1
	}
	if skip >= len(v) {
		return -1
	}
	return v[skip]
}
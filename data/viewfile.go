package data

import (
	"encoding/json"
	"os"
)

const viewfile_path = "views.json"

type ViewFile struct {
	Views map[int]View
	path string
}

type View struct {
	ID       int            `json:"ID"`
	Camera   Camera         `json:"Camera"`
	Panels   []Panel        `json:"Panels"`
	Icons    []PanelIcon    `json:"Icons"`
	Labels   []PanelLabel   `json:"Labels"`
	Flags    []PanelFlag    `json:"Flags"`
	Counters []PanelCounter `json:"Counters"`
}

type Camera struct {
	Position Vector `json:"Position"`
	Size     Vector `json:"Size"`
	Offset   Vector `json:"Offset"`
}

type Panel struct {
	Cell
	Position Vector `json:"Position"`
	Size     Vector `json:"Size"`
}

type PanelLabel struct {
	Text     string `json:"Text"`
	Color    Color  `json:"Color"`
	Position Vector `json:"Position"`
}

type PanelFlag struct {
	Counter Name `json:"Counter"`
	Cell
	Position Vector `json:"Position"`
}

type PanelIcon struct {
	Cell
	Position Vector `json:"Position"`
}

type PanelCounter struct {
	Counter  Name   `json:"Counter"`
	Color    Color  `json:"Color"`
	Length   int    `json:"Length"`
	Position Vector `json:"Position"`
}

func (wp *WorldPath) LoadViewFile() (*ViewFile, error) {
	vf := ViewFile{
		path: wp.Join(viewfile_path),
	}
	if err := vf.load(); err != nil {
		return nil, err
	}
	return &vf, nil
}

func (vf *ViewFile) load() error {
	// read and decode
	b, err := os.ReadFile(vf.path)
	if err != nil {
		return err
	}
	views := struct {
		Views []View `json:"Views"`
	}{}
	if err := json.Unmarshal(b, &views); err != nil {
		return err
	}
	vf.Views = make(map[int]View, len(views.Views))
	for _, view := range views.Views {
		vf.Views[view.ID] = view
	}
	return nil
}

func (vf *ViewFile) GetView(id int) View {
	if view, ok := vf.Views[id]; ok {
		return view
	}
	return default_view
}

var default_view = View{
	ID: 0,
	Camera: Camera{
		Position: VectorZero,
		Size:     Vector{X: 40, Y: 25},
		Offset:   Vector{X: 20, Y: 12},
	},
}

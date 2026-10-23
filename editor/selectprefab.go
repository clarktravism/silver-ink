package main

import (
	"log"
	"silver-ink/data"
	"silver-ink/engine"
)

const prefab_file = "prefabs.json"

type SelectPrefab struct {
	Callback func(SelectPrefabCallback)
}

type SelectPrefabCallback struct {
	Accepted bool
	Prefab data.Prefab
}

func (sp *SelectPrefab) Push() {

	if engine.PrefabsFile.HasChanged() {
		if err := engine.PrefabsFile.Reload(); err != nil {
			log.Printf("error while reading prefabs: %v\r\n", err)
		}
	}

	// create scroll box
	size := data.Vector{ X:38, Y: 41 }
	sb := engine.ScrollBox{
		Rectangle:  data.Rect{
			Position: engine.CenterPos(size),
			Size: size,
		},
		Title: "Select Prefab",
		Style: scroll_box_style,
		Callback: sp.callbackScrollBox,
		ForceSelect: false,
	}

	for _, prefab := range engine.PrefabsFile.Prefabs {
		sb.AppendLine(engine.ScrollBoxLine{
			Text: prefab.Name.String(),
			Link: SELECT_PREFAB,
			Ext: prefab,
			Centered: false,
			Icon: prefab.Cell,
		})
	}

	sb.AppendText("", false)
	sb.AppendLink("Cancel", CANCEL, nil)
	sb.Push()
}

var SELECT_PREFAB = data.NewName("SELECT_ELEMENT")

func (sp *SelectPrefab) callbackScrollBox(callback engine.ScrollBoxCallback) {
	if !callback.Accepted || callback.Link == CANCEL {
		sp.Callback(SelectPrefabCallback{
			Accepted: false,
		})
		return
	}

	if callback.Link == SELECT_PREFAB {
		prefab, _ := callback.Ext.(data.Prefab)
		
		sp.Callback(SelectPrefabCallback{
			Accepted: true,
			Prefab: prefab,
		})
	}
}
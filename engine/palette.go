package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

var palette_name data.Name
var palette_folder = "palettes"

func LoadPalette(name data.Name) error {
	if palette_name == name {
		return nil
	}
	filename := name.String() + png_ext
	if WorldPath == nil {
		return render.LoadPalette(filename)
	}
	return render.LoadPalette(WorldPath.Join(palette_folder, filename))
}
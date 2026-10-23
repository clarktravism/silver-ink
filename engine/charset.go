package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

var charset_name data.Name
var png_ext = ".png"
var charset_folder = "charsets"

func LoadCharset(name data.Name) error {
	if charset_name == name {
		return nil
	}
	filename := name.String() + png_ext
	if WorldPath == nil {
		return render.LoadCharset(_renderer, filename)
	}
	return render.LoadCharset(_renderer, WorldPath.Join(charset_folder, filename))
}
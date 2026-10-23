package data

import (
	"path/filepath"
)

type WorldPath struct {
	Path string
}

func (wp *WorldPath) Join(elem ...string) string {
	elem = append([]string{wp.Path}, elem...)
	return filepath.Join(elem...)
}
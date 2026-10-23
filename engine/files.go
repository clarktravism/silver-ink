package engine

import (
	"log"
	"silver-ink/data"
	"strings"
)

var WorldPath *data.WorldPath
var WorldFile *data.WorldFile
var ElementsFile *data.ElementsFile
var PrefabsFile *data.PrefabsFile
//var CharsetsFile *data.AssetFile
//var PalettesFile *data.AssetFile
var ViewFile *data.ViewFile
var Board *data.Board

func SetWorldPath(path string) error {
	WorldPath = nil
	WorldFile = nil
	ElementsFile = nil
	PrefabsFile = nil
	Board = nil

	WorldPath = &data.WorldPath{Path: path}

	var err error
	WorldFile, err = WorldPath.LoadWorldFile()
	if err != nil {
		log.Println(err)
	}
	ElementsFile, err = WorldPath.LoadElementsFile()
	if err != nil {
		log.Println(err)
	}
	PrefabsFile, err = WorldPath.LoadPrefabFile()
	if err != nil {
		log.Println(err)
	}
	ViewFile, err = WorldPath.LoadViewFile()
	if err != nil {
		log.Println(err)
	}
	
	return nil
}

func SetBoard(board *data.Board) {
	Board = board
	if err := LoadCharset(board.Charset); err != nil {
		log.Println(err)
	}
	if err := LoadPalette(board.Palette); err != nil {
		log.Println(err)
	}
}

func FormatPath(name string, ext string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "_") + ext
}

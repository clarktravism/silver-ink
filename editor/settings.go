package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

const settings_file = "./settings.json"

type Settings struct {
	WorldsPath string
}

func init() {
	if err := loadSettings(); err != nil {
		log.Printf("error while loading settings: %v\r\n", err)
	}
}

var settings Settings

func loadSettings() error {
	// read and decode
	b, err := os.ReadFile(settings_file)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		settings = default_settings
		return saveSettings()
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &settings); err != nil {
		return err
	}
	return nil
}

func saveSettings() error {
	b, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	if err := os.WriteFile(settings_file, b, 0644); err != nil {
		return err
	}
	return nil
}

var default_settings = Settings{
	WorldsPath: "./worlds",
}
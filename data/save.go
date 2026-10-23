package data
/*
import (
	"path/filepath"
	"encoding/json"
	"os"
)

const work_dir = "work"
const saves_dir = "saves"
const save_file = "save.json"

type Save struct {
	CurrentBoard uint16 `json:"CurrentBoard"`
	Counters
}

// save.json
func (save *Save) LoadSaveFile(pathname string) error {
	savePath := PathTo(filepath.Join(saves_dir, pathname))
	workPath := PathTo(work_dir)
	
	//copy save folder to work files
	if err := os.RemoveAll(workPath); err != nil {
		return err
	}
	saveFS := os.DirFS(savePath)
	if err := os.CopyFS(workPath, saveFS); err != nil {
		return err
	}
	
	// read and decode save.json
	saveFile := filepath.Join(savePath, save_file)
	b, err := os.ReadFile(saveFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &save); err != nil {
		return err
	}
	return nil
}

func (save *Save) SaveGame(pathname string) error {
	savePath := PathTo(filepath.Join(saves_dir, pathname))
	workPath := PathTo(work_dir)

	//copy work files to save folder
	if err := os.RemoveAll(savePath); err != nil {
		return err
	}
	workFS := os.DirFS(workPath)
	if err := os.CopyFS(savePath, workFS); err != nil {
		return err
	}

	//save file json
	b, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return err
	}
	saveFilePath := filepath.Join(savePath, save_file) 
	if err := os.WriteFile(saveFilePath, b, 0644); err != nil {
		return err
	}

	return nil
}

func (world *World) NewWorkSave() Save {
	return Save{
		Counters: world.Counters,
		CurrentBoard: world.TitleBoard,
	}
}
*/
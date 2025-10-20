package jsonwriter

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveJSON(path string, data interface{}) error {
	if data == nil {
		return fmt.Errorf("no data to save")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

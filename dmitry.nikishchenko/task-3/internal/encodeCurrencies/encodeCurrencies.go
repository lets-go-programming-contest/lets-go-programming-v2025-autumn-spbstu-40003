package encodeCurrencies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveCurrencies(path string, data interface{}) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create directories: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("cannot close file: %w", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("cannot encode file properly: %w", err)
	}

	return nil
}

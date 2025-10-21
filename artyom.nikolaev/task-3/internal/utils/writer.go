package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArtttNik/task-3/internal/parser"
)

func WriteJSONToFile(path string, data []parser.Currency) error {
	const permissions = 0o755

	err := os.MkdirAll(filepath.Dir(path), permissions)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close file: %w", err))
		}
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	err = enc.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode JSON data to file: %w", err)
	}

	return nil
}

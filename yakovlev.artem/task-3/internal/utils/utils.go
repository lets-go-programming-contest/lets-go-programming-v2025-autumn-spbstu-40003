package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const dirPerm = 0o755

// WriteJSONToFile creates directories (if needed) and writes pretty JSON.
func WriteJSONToFile(path string, data any) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(fmt.Errorf("close file: %w", cerr))
		}
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Egor1726/task-3/internal/parser"
)

const dirPermissions = 0o755

func WriteJSONToFile(path string, data []parser.Currency) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPermissions); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(fmt.Errorf("failed to close file: %w", cerr))
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

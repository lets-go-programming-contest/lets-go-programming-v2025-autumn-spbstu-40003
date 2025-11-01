package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Artem-Hack/task-3/internal/parser"
)

func ExportJSON(path string, data []parser.Currency) error {
	const filePerm = 0o755

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, filePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Println("close error:", cerr)
		}
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("JSON encoding error: %w", err)
	}

	return nil
}

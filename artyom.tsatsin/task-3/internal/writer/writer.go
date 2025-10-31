package writer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Artem-Hack/task-3/internal/parser"
)

func ExportJSON(path string, data []parser.Currency) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("Failed to create directory: " + err.Error())
	}

	file, err := os.Create(path)
	if err != nil {
		panic("Failed to create JSON: " + err.Error())
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(data); err != nil {
		panic("JSON encoding error: " + err.Error())
	}

	return nil
}

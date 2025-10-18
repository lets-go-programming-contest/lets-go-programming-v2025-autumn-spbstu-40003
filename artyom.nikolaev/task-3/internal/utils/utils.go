package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ArtttNik/task-3/internal/parser"
)

func WriteJSONToFile(path string, data []parser.Currency) error {
	const permissions = 0o755

	lastSlash := strings.LastIndex(path, "/")
	if lastSlash != -1 {
		dir := path[:lastSlash]

		err := os.MkdirAll(dir, permissions)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
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

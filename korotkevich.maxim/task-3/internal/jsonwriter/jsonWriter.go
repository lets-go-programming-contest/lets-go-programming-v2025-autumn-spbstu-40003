package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteJSON(outputPath string, data interface{}) error {
	const permissions = 0o755

	directory := filepath.Dir(outputPath)
	if err := os.MkdirAll(directory, permissions); err != nil {
		return fmt.Errorf("error: failed to create output directory %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error: fail with creating file: %w", err)
	}

	defer func() {
		err := outputFile.Close()
		if err != nil {
			panic(fmt.Errorf("error: problem with closing file: %w", err))
		}
	}()

	enc := json.NewEncoder(outputFile)
	enc.SetIndent("", "  ")

	err = enc.Encode(data)
	if err != nil {
		return fmt.Errorf("error: failed to encode JSON in file: %w", err)
	}

	return nil
}

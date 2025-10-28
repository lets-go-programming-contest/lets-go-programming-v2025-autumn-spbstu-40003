package datawriter

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveAsJSON(path string, data any) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	return nil
}

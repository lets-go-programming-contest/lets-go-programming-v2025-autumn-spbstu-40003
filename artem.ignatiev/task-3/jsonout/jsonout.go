package jsonout

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteJSON(file string, data any) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("cannot encode json: %w", err)
	}
	return nil
}

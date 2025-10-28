package jsonwriter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrNoData = errors.New("no data to save")

const dirPerm = 0755

func Save(path string, data interface{}) error {
	if data == nil {
		return ErrNoData
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("cannot create directory: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	defer func(f *os.File) {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", cerr)
		}
	}(file)

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("cannot encode JSON: %w", err)
	}

	return nil
}

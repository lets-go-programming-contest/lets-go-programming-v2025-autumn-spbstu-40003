package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Save(valutes []JSONValute, path string) error {
	valutesJSON, err := json.MarshalIndent(valutes, "", " ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	const (
		filePermissionCode      = 0o644
		directoryPermissionCode = 0o755
	)

	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, directoryPermissionCode); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	err = os.WriteFile(path, valutesJSON, filePermissionCode)
	if err != nil {
		return fmt.Errorf("json write: %w", err)
	}

	return nil
}

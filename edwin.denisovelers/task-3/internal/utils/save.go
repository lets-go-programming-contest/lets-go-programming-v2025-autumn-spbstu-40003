package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func Save(valutes []JSONValute, path string) error {
	valutesJSON, err := json.MarshalIndent(valutes, "", " ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	const permissionCode = 0o0644

	_, err = os.Create(path)
	if err != nil {
		return fmt.Errorf("json file create: %w", err)
	}

	err = os.WriteFile(path, valutesJSON, permissionCode)
	if err != nil {
		return fmt.Errorf("json write: %w", err)
	}

	return nil
}

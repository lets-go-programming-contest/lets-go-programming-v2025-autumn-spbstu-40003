package valutemanager

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"

	"github.com/KiRy6A/task-3/internal/data"
)

const (
	filePermission      = 0o644
	directoryPermission = 0o755
)

func Read(path string) (data.Valutes, error) {
	var valutesData data.Valutes

	file, err := os.Open(path)
	if err != nil {
		return valutesData, fmt.Errorf("failed opening file: %w", err)
	}

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&valutesData)
	if err != nil {
		return valutesData, fmt.Errorf("failed decoding data: %w", err)
	}

	err = file.Close()
	if err != nil {
		return valutesData, fmt.Errorf("failed closing file: %w", err)
	}

	return valutesData, nil
}

func Write(path string, valutes data.Valutes) error {
	data, err := json.Marshal(valutes.AllValutes)
	if err != nil {
		return fmt.Errorf("failed serialization data: %w", err)
	}

	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, directoryPermission); err != nil {
		return fmt.Errorf("failed creating directory: %w", err)
	}

	err = os.WriteFile(path, data, filePermission)
	if err != nil {
		return fmt.Errorf("failed writing file: %w", err)
	}

	return nil
}

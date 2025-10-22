package valutemanager

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
)

const (
	filePermission      = 0o644
	directoryPermission = 0o755
)

type Valute struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

type Valutes struct {
	AllValutes []Valute `xml:"Valute"`
}

func Read(path string) (Valutes, error) {
	var valutesData Valutes

	file, err := os.Open(path)
	if err != nil {
		return valutesData, fmt.Errorf("failed opening file: %w", err)
	}

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&valutesData)
	if err != nil {
		file.Close()
		return valutesData, fmt.Errorf("failed decoding data: %w", err)
	}

	file.Close()

	return valutesData, nil
}

func Write(path string, valutes Valutes) error {
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

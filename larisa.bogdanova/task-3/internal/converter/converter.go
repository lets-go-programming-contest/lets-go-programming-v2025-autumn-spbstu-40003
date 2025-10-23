package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/badligyg/task-3/internal/models"
	"golang.org/x/net/html/charset"
)

const dirPerm = 0o755

func ReadCurrencies(filename string) ([]models.Currency, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open XML file: %w", err)
	}
	defer file.Close()

	var data models.ValCurs
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("parse XML: %w", err)
	}

	sort.Slice(data.Items, func(i, j int) bool {
		return data.Items[i].Value > data.Items[j].Value
	})

	return data.Items, nil
}

func WriteCurrencies(currencies []models.Currency, filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(currencies); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	return nil
}

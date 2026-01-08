package processor

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"task-3/internal/currency"

	"golang.org/x/net/html/charset"
)

func LoadExchangeData(path string) (*currency.ExchangeData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Println("close error:", cerr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var exchangeData currency.ExchangeData
	if err := decoder.Decode(&exchangeData); err != nil {
		return nil, fmt.Errorf("XML decoding error: %w", err)
	}

	return &exchangeData, nil
}

func GetSortedCurrencies(data *currency.ExchangeData) []currency.Item {
	items := make([]currency.Item, len(data.Items))
	copy(items, data.Items)

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].RateValue > items[j].RateValue
	})

	return items
}

func ExportToJSON(path string, data []currency.Item) error {
	const filePerm = 0o755

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, filePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Println("close error:", cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("JSON encoding error: %w", err)
	}

	return nil
}

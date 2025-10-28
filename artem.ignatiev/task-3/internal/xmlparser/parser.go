package xmlparser

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Amount float64

func (a *Amount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string
	if err := d.DecodeElement(&raw, &start); err != nil {
		return fmt.Errorf("failed to decode value: %w", err)
	}

	raw = strings.ReplaceAll(raw, ",", ".")
	raw = strings.TrimSpace(raw)

	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fmt.Errorf("failed to parse float: %w", err)
	}

	*a = Amount(val)
	return nil
}

type Currency struct {
	CodeNum  int    `xml:"NumCode"  json:"num_code"`
	CodeChar string `xml:"CharCode" json:"char_code"`
	Value    Amount `xml:"Value"    json:"value"`
}

type CurrencyList struct {
	Currencies []Currency `xml:"Valute"`
}

func ParseXML(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open XML file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data CurrencyList
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("parse XML: %w", err)
	}

	sort.Slice(data.Currencies, func(i, j int) bool {
		return data.Currencies[i].Value > data.Currencies[j].Value
	})

	return data.Currencies, nil
}

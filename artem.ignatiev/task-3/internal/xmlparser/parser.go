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
	var valueStr string
	if err := d.DecodeElement(&valueStr, &start); err != nil {
		return fmt.Errorf("failed to decode value: %w", err)
	}

	valueStr = strings.ReplaceAll(valueStr, ",", ".")
	valueStr = strings.TrimSpace(valueStr)

	val, err := strconv.ParseFloat(valueStr, 64)
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

type CurrencyFile struct {
	Currencies []Currency `xml:"Valute"`
}

func ReadCurrencies(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open XML file: %w", err)
	}
	defer func(f *os.File) {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", cerr)
		}
	}(file)

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data CurrencyFile
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	sort.Slice(data.Currencies, func(i, j int) bool {
		return data.Currencies[i].Value > data.Currencies[j].Value
	})

	return data.Currencies, nil
}

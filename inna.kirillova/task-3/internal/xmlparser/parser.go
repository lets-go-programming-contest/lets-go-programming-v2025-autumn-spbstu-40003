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

type ValueFloat float64

func (v *ValueFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valStr string

	if err := d.DecodeElement(&valStr, &start); err != nil {
		return fmt.Errorf("failed to decode value: %w", err)
	}

	valStr = strings.ReplaceAll(valStr, ",", ".")
	valStr = strings.TrimSpace(valStr)

	f, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse float: %w", err)
	}

	*v = ValueFloat(f)

	return nil
}

type ExchangeTrade struct {
	NumCode  int       `json:"num_code"  xml:"NumCode"`
	CharCode string    `json:"char_code" xml:"CharCode"`
	Value    ValueFloat `json:"value"     xml:"Value"`
}

type ExchangeData struct {
	Trades []ExchangeTrade `xml:"Valute"`
}

func ParseXML(path string) ([]ExchangeTrade, error) {
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

	var data ExchangeData
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("parse xml: %w", err)
	}

	sort.Slice(data.Trades, func(i, j int) bool {
		return data.Trades[i].Value > data.Trades[j].Value
	})

	return data.Trades, nil
}

package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ValueFloat float64

func (v *ValueFloat) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var raw string

	if err := dec.DecodeElement(&raw, &start); err != nil {
		return fmt.Errorf("failed to decode Value: %w", err)
	}

	raw = strings.TrimSpace(strings.ReplaceAll(raw, ",", "."))

	if val, err := strconv.ParseFloat(raw, 64); err != nil {
		return fmt.Errorf("invalid Value: %w", err)
	} else {
		*v = ValueFloat(val)
	}

	return nil
}

type Currency struct {
	Code  int        `json:"num_code"  xml:"NumCode"`
	Char  string     `json:"char_code" xml:"CharCode"`
	Value ValueFloat `json:"value"     xml:"Value"`
}

func LoadXML(path string) ([]Currency, error) {
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

	var list []Currency

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("XML reading error: %w", err)
		}

		start, ok := token.(xml.StartElement)

		if !ok || start.Name.Local != "Valute" {
			continue
		}

		var curr Currency
		if err := decoder.DecodeElement(&curr, &start); err != nil {
			return nil, fmt.Errorf("currency decode error: %w", err)
		}

		list = append(list, curr)
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Value > list[j].Value
	})

	return list, nil
}

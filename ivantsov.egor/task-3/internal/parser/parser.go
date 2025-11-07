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

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("token read error: %w", err)
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "NumCode":
				var val string
				if err := decoder.DecodeElement(&val, &elem); err != nil {
					return fmt.Errorf("decode NumCode error: %w", err)
				}
				val = strings.TrimSpace(val)
				if val != "" {
					c.NumCode, err = strconv.Atoi(val)
					if err != nil {
						return fmt.Errorf("invalid NumCode: %w", err)
					}
				}
			case "CharCode":
				if err := decoder.DecodeElement(&c.CharCode, &elem); err != nil {
					return fmt.Errorf("decode CharCode error: %w", err)
				}
				c.CharCode = strings.TrimSpace(c.CharCode)
			case "Value":
				var val string
				if err := decoder.DecodeElement(&val, &elem); err != nil {
					return fmt.Errorf("decode Value error: %w", err)
				}
				val = strings.ReplaceAll(strings.TrimSpace(val), ",", ".")
				if val != "" {
					c.Value, err = strconv.ParseFloat(val, 64)
					if err != nil {
						return fmt.Errorf("invalid Value: %w", err)
					}
				}
			}
		case xml.EndElement:
			if elem.Name.Local == start.Name.Local {
				return nil
			}
		}
	}
	return nil
}

func ParseFile(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(fmt.Errorf("failed to close XML file: %w", cerr))
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var currencies []Currency
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("XML read error: %w", err)
		}

		startElem, ok := token.(xml.StartElement)
		if !ok || startElem.Name.Local != "Valute" {
			continue
		}

		var currency Currency
		if err := decoder.DecodeElement(&currency, &startElem); err != nil {
			return nil, fmt.Errorf("decode element error: %w", err)
		}

		currencies = append(currencies, currency)
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies, nil
}

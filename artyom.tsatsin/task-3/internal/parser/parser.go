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
	Code  int     `json:"num_code"  xml:"NumCode"`
	Char  string  `json:"char_code" xml:"CharCode"`
	Value float64 `json:"value"     xml:"Value"`
}

func (c *Currency) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("XML token error: %w", err)
		}

		switch elem := token.(type) {
		case xml.StartElement:
			if err := c.handleStartElement(dec, elem); err != nil {
				return err
			}
		case xml.EndElement:
			if elem.Name.Local == start.Name.Local {
				return nil
			}
		}
	}

	return nil
}

func (c *Currency) handleStartElement(dec *xml.Decoder, elem xml.StartElement) error {
	switch elem.Name.Local {
	case "NumCode":
		var numCodeStr string

		if err := dec.DecodeElement(&numCodeStr, &elem); err != nil {
			return fmt.Errorf("failed to decode NumCode: %w", err)
		}

		c.Code, _ = strconv.Atoi(strings.TrimSpace(numCodeStr))

	case "CharCode":
		var charCodeStr string

		if err := dec.DecodeElement(&charCodeStr, &elem); err != nil {
			return fmt.Errorf("failed to decode CharCode: %w", err)
		}

		c.Char = strings.TrimSpace(charCodeStr)

	case "Value":
		var valueStr string
		if err := dec.DecodeElement(&valueStr, &elem); err != nil {
			return fmt.Errorf("failed to decode Value: %w", err)
		}

		valueStr = strings.ReplaceAll(strings.TrimSpace(valueStr), ",", ".")
		val, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return fmt.Errorf("invalid Value: %w", err)
		}

		c.Value = val
	}

	return nil
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

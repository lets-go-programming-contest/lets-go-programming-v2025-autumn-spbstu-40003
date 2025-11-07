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

//nolint:gocognit // function intentionally complex due to XML decoding structure
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
			if err := c.processElement(decoder, elem); err != nil {
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

func (c *Currency) processElement(decoder *xml.Decoder, elem xml.StartElement) error {
	switch elem.Name.Local {
	case "NumCode":
		return c.decodeNumCode(decoder, elem)
	case "CharCode":
		return c.decodeCharCode(decoder, elem)
	case "Value":
		return c.decodeValue(decoder, elem)
	}

	return nil
}

func (c *Currency) decodeNumCode(decoder *xml.Decoder, elem xml.StartElement) error {
	var val string

	if err := decoder.DecodeElement(&val, &elem); err != nil {
		return fmt.Errorf("decode NumCode error: %w", err)
	}

	val = strings.TrimSpace(val)

	if val == "" {
		return nil
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("parse NumCode error: %w", err)
	}

	c.NumCode = num

	return nil
}

func (c *Currency) decodeCharCode(decoder *xml.Decoder, elem xml.StartElement) error {
	var val string

	if err := decoder.DecodeElement(&val, &elem); err != nil {
		return fmt.Errorf("decode CharCode error: %w", err)
	}

	val = strings.TrimSpace(val)

	if val != "" {
		c.CharCode = val
	}

	return nil
}

func (c *Currency) decodeValue(decoder *xml.Decoder, elem xml.StartElement) error {
	var val string

	if err := decoder.DecodeElement(&val, &elem); err != nil {
		return fmt.Errorf("decode Value error: %w", err)
	}

	val = strings.ReplaceAll(strings.TrimSpace(val), ",", ".")

	if val == "" {
		return nil
	}

	value, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fmt.Errorf("parse Value error: %w", err)
	}

	c.Value = value

	return nil
}

func ParseFile(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file error: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(fmt.Errorf("failed to close file: %w", cerr))
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

		var cur Currency

		if err := decoder.DecodeElement(&cur, &startElem); err != nil {
			return nil, fmt.Errorf("currency decode error: %w", err)
		}

		currencies = append(currencies, cur)
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies, nil
}

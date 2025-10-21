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

func (currency *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return fmt.Errorf("decode token error: %w", err)
		}

		switch elem := token.(type) {
		case xml.StartElement:
			err := currency.processElement(decoder, elem)
			if err != nil {
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

func (currency *Currency) processElement(decoder *xml.Decoder, elem xml.StartElement) error {
	switch elem.Name.Local {
	case "NumCode":
		return currency.decodeNumCode(decoder, elem)
	case "CharCode":
		return currency.decodeCharCode(decoder, elem)
	case "Value":
		return currency.decodeValue(decoder, elem)
	}

	return nil
}

func (currency *Currency) decodeNumCode(decoder *xml.Decoder, elem xml.StartElement) error {
	var codeValue string

	err := decoder.DecodeElement(&codeValue, &elem)
	if err != nil {
		return fmt.Errorf("failed to decode element: %w", err)
	}

	codeValue = strings.TrimSpace(codeValue)
	if codeValue == "" {
		return nil
	}

	parsedCode, err := strconv.Atoi(codeValue)
	if err != nil {
		return fmt.Errorf("invalid NumCode: %w", err)
	}

	currency.NumCode = parsedCode

	return nil
}

func (currency *Currency) decodeCharCode(decoder *xml.Decoder, element xml.StartElement) error {
	var charCodeValue string

	err := decoder.DecodeElement(&charCodeValue, &element)
	if err != nil {
		return fmt.Errorf("failed to decode element: %w", err)
	}

	currency.CharCode = strings.TrimSpace(charCodeValue)

	return nil
}

func (currency *Currency) decodeValue(decoder *xml.Decoder, element xml.StartElement) error {
	var valueString string

	err := decoder.DecodeElement(&valueString, &element)
	if err != nil {
		return fmt.Errorf("failed to decode element: %w", err)
	}

	valueString = strings.ReplaceAll(strings.TrimSpace(valueString), ",", ".")
	if valueString == "" {
		return nil
	}

	parsedValue, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		return fmt.Errorf("invalid Value: %w", err)
	}

	currency.Value = parsedValue

	return nil
}

func ParseFile(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no such file or directory: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("failed to close file: %w", err))
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

			return nil, fmt.Errorf("failed to read token: %w", err)
		}

		startElem, ok := token.(xml.StartElement)
		if !ok || startElem.Name.Local != "Valute" {
			continue
		}

		var currency Currency

		err = decoder.DecodeElement(&currency, &startElem)
		if err != nil {
			return nil, fmt.Errorf("failed to decode currency: %w", err)
		}

		currencies = append(currencies, currency)
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies, nil
}

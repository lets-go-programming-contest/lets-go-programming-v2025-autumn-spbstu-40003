package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type valCurs struct {
	Valutes []valute `xml:"Valute"`
}

type valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Value    string `xml:"Value"`
}

func ParseCBR(path string) ([]Currency, error) {
	currencies, err := ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return currencies, nil
}

func ParseFile(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no such file or directory: %w", err)
		}

		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("failed to close file: %w", err))
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var curs valCurs

	err = decoder.Decode(&curs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	currencies := convertToModel(curs.Valutes)
	sortCurrencies(currencies)

	return currencies, nil
}

func convertToModel(valutes []valute) []Currency {
	result := make([]Currency, len(valutes))
	for i, val := range valutes {
		result[i] = convertSingle(val)
	}

	return result
}

func convertSingle(valute valute) Currency {
	numCode, err := parseInt(strings.TrimSpace(valute.NumCode))
	if err != nil {
		numCode = 0
	}

	charCode := strings.TrimSpace(valute.CharCode)

	value, err := parseValue(valute.Value)
	if err != nil {
		value = 0
	}

	return Currency{
		NumCode:  numCode,
		CharCode: charCode,
		Value:    value,
	}
}

func parseInt(s string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, fmt.Errorf("failed to convert to int: %w", err)
	}

	return value, nil
}

func parseValue(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}

	return value, nil
}

func sortCurrencies(currencies []Currency) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
}

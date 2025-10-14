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
	Values []valute `xml:"Valute"`
}

type valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Value    string `xml:"Value"`
}

func ParseCBR(path string) ([]Currency, error) {
	list, err := parseFile(path)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	return list, nil
}

func parseFile(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no such file: %w", err)
		}

		return nil, fmt.Errorf("open xml: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(fmt.Errorf("close xml: %w", cerr))
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var curs valCurs

	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	result := convertToModel(curs.Values)

	sortCurrencies(result)

	return result, nil
}

func convertToModel(vals []valute) []Currency {
	result := make([]Currency, len(vals))

	for i, valuteItem := range vals {
		result[i] = convertSingle(valuteItem)
	}

	return result
}

func convertSingle(valuteItem valute) Currency {
	numCode, err := parseInt(valuteItem.NumCode)
	if err != nil {
		numCode = 0
	}

	value, err := parseFloat(valuteItem.Value)
	if err != nil {
		value = 0
	}

	return Currency{
		NumCode:  numCode,
		CharCode: strings.TrimSpace(valuteItem.CharCode),
		Value:    value,
	}
}

func parseInt(s string) (int, error) {
	s = strings.TrimSpace(s)

	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("atoi %q: %w", s, err)
	}

	return n, nil
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", "."))

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parseFloat %q: %w", s, err)
	}

	return f, nil
}

func sortCurrencies(cc []Currency) {
	sort.Slice(cc, func(i, j int) bool {
		return cc[i].Value > cc[j].Value
	})
}

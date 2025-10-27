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
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Nominal  int     `json:"-"         xml:"Nominal"`
	Value    float64 `json:"value"     xml:"Value"`
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type raw struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Nominal  string `xml:"Nominal"`
		Value    string `xml:"Value"`
	}

	var rawVal raw
	if err := decoder.DecodeElement(&rawVal, &start); err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	if num, err := parseIntFromString(rawVal.NumCode); err == nil {
		c.NumCode = num
	}
	c.CharCode = strings.TrimSpace(rawVal.CharCode)
	if nom, err := parseIntFromString(rawVal.Nominal); err == nil {
		c.Nominal = nom
	}
	if val, err := parseFloatFromString(rawVal.Value); err == nil {
		c.Value = val
	}

	return nil
}

func ParseCBR(path string) ([]Currency, error) {
	return parseFile(path)
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

	var curs struct {
		Values []Currency `xml:"Valute"`
	}
	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	sortCurrencies(curs.Values)
	return curs.Values, nil
}

func parseIntFromString(str string) (int, error) {
	s := strings.TrimSpace(str)
	number, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("atoi %q: %w", s, err)
	}
	return number, nil
}

func parseFloatFromString(str string) (float64, error) {
	s := strings.TrimSpace(strings.ReplaceAll(str, ",", "."))
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parseFloat %q: %w", s, err)
	}
	return value, nil
}

func sortCurrencies(cc []Currency) {
	sort.Slice(cc, func(i, j int) bool {
		return cc[i].Value > cc[j].Value
	})
}

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

	numCode, err := parseInt(rawVal.NumCode)
	if err != nil {
		c.NumCode = 0
	} else {
		c.NumCode = numCode
	}

	nominal, err := parseInt(rawVal.Nominal)
	if err != nil {
		c.Nominal = 0
	} else {
		c.Nominal = nominal
	}

	c.CharCode = strings.TrimSpace(rawVal.CharCode)

	val, err := parseFloat(rawVal.Value)
	if err != nil {
		c.Value = 0
	} else {
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

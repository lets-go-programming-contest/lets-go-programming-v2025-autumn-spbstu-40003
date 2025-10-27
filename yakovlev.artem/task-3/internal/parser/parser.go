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
	Value    float64 `json:"value"     xml:"Value"`
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type raw struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}
	var r raw
	if err := decoder.DecodeElement(&r, &start); err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	num, _ := parseIntFromString(r.NumCode)
	val, _ := parseFloatFromString(r.Value)

	c.NumCode = num
	c.CharCode = strings.TrimSpace(r.CharCode)
	c.Value = val

	return nil
}

func ParseCBR(path string) ([]Currency, error) {
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

	sort.Slice(curs.Values, func(i, j int) bool {
		return curs.Values[i].Value > curs.Values[j].Value
	})

	return curs.Values, nil
}

func parseIntFromString(s string) (int, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("atoi %q: %w", s, err)
	}

	return n, nil
}

func parseFloatFromString(s string) (float64, error) {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", "."))
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parseFloat %q: %w", s, err)
	}

	return f, nil
}

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
	NumCode  int     `xml:"NumCode"  json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Nominal  int     `xml:"Nominal"  json:"-"`
	Value    float64 `xml:"Value"    json:"value"`
}

func (c *Currency) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type raw struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Nominal  string `xml:"Nominal"`
		Value    string `xml:"Value"`
	}
	var r raw
	if err := d.DecodeElement(&r, &start); err != nil {
		return err
	}

	var err error
	if c.NumCode, err = parseInt(r.NumCode); err != nil {
		c.NumCode = 0
	}
	c.CharCode = strings.TrimSpace(r.CharCode)
	if c.Nominal, err = parseInt(r.Nominal); err != nil {
		c.Nominal = 0
	}
	if c.Value, err = parseFloat(r.Value); err != nil {
		c.Value = 0
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

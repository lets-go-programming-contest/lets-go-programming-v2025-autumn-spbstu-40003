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
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no such file: %w", err)
		}
		return nil, fmt.Errorf("open xml: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(fmt.Errorf("close xml: %w", err))
		}
	}()

	dec := xml.NewDecoder(f)
	dec.CharsetReader = charset.NewReaderLabel

	var curs valCurs
	if err := dec.Decode(&curs); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	out := convertToModel(curs.Values)
	sortCurrencies(out)
	return out, nil
}

func convertToModel(vs []valute) []Currency {
	res := make([]Currency, len(vs))
	for i, v := range vs {
		res[i] = convertSingle(v)
	}
	return res
}

func convertSingle(v valute) Currency {
	num, err := parseInt(v.NumCode)
	if err != nil {
		num = 0
	}
	val, err := parseFloat(v.Value)
	if err != nil {
		val = 0
	}
	return Currency{
		NumCode:  num,
		CharCode: strings.TrimSpace(v.CharCode),
		Value:    val,
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
	sort.Slice(cc, func(i, j int) bool { return cc[i].Value > cc[j].Value })
}

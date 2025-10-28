package xmlhandler

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type CurrencyList struct {
	Date     string     `xml:"Date,attr"`
	Name     string     `xml:"name,attr"`
	Currency []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int      `xml:"NumCode" json:"num_code"`
	CharCode string   `xml:"CharCode" json:"char_code"`
	Value    FloatNum `xml:"Value" json:"value"`
}

type FloatNum float64

func (f *FloatNum) UnmarshalText(text []byte) error {
	s := strings.ReplaceAll(string(text), ",", ".")
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid float format: %q", s)
	}
	*f = FloatNum(val)
	return nil
}

func LoadCurrencies(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open XML file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var list CurrencyList
	if err := decoder.Decode(&list); err != nil && err != io.EOF {
		return nil, fmt.Errorf("decode XML: %w", err)
	}

	return list.Currency, nil
}

func SortDescending(currencies []Currency) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
}

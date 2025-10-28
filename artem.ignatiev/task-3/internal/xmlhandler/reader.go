package xmlhandler

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

var ErrInvalidFloat = errors.New("invalid float format")

type CurrencyList struct {
	Date     string     `xml:"Date,attr"`
	Name     string     `xml:"name,attr"`
	Currency []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int      `json:"num_code"  xml:"NumCode"`
	CharCode string   `json:"char_code" xml:"CharCode"`
	Value    FloatNum `json:"value"     xml:"Value"`
}

type FloatNum float64

func (f *FloatNum) UnmarshalText(text []byte) error {
	s := strings.ReplaceAll(string(text), ",", ".")

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("%w: %q", ErrInvalidFloat, s)
	}

	*f = FloatNum(val)

	return nil
}

func LoadCurrencies(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open XML file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(cerr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var list CurrencyList

	if err := decoder.Decode(&list); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("decode XML: %w", err)
	}

	return list.Currency, nil
}

func SortDescending(currencies []Currency) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
}

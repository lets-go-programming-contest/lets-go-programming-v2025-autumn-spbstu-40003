package xmlreader

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Amount float64

func (a *Amount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	s = strings.ReplaceAll(s, ",", ".")
	val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return fmt.Errorf("convert to float: %w", err)
	}

	*a = Amount(val)
	return nil
}

type Currency struct {
	CodeNum  int    `xml:"NumCode" json:"num_code"`
	CodeChar string `xml:"CharCode" json:"char_code"`
	Value    Amount `xml:"Value" json:"value"`
}

type CurrencyFile struct {
	List []Currency `xml:"Valute"`
}

func ReadCurrencies(path string) ([]Currency, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read XML file: %w", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(content)))
	decoder.CharsetReader = charset.NewReaderLabel

	var data CurrencyFile
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("cannot decode XML: %w", err)
	}

	sort.Slice(data.List, func(i, j int) bool {
		return data.List[i].Value > data.List[j].Value
	})

	return data.List, nil
}

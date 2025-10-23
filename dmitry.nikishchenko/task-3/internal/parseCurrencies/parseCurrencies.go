package parseCurrencies

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type FloatComma float64

func (f *FloatComma) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.ReplaceAll(s, ",", ".")

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*f = FloatComma(v)
	return nil
}

type ValCurs struct {
	XMLrootName xml.Name `xml:"Valuta"`
	Valutes     []Valute `xml:"Item"`
}

type Valute struct {
	NumCode  int        `xml:"ISO_Num_Code" json:"num_code"`
	CharCode string     `xml:"ISO_Char_Code" json:"char_code"`
	Value    FloatComma `xml:"Nominal" json:"nominal"`
}

func LoadCurrencies(path string) ([]Valute, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	var valCurs ValCurs
	if err := xml.Unmarshal(data, &valCurs); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	sort.Slice(valCurs.Valutes, func(i, j int) bool {
		return valCurs.Valutes[i].Value > valCurs.Valutes[j].Value
	})

	return valCurs.Valutes, nil
}

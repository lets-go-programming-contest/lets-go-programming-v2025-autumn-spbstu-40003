package parseCurrencies

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
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
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID        string     `xml:"ID,attr" json:"id"`
	NumCode   int        `xml:"NumCode" json:"num_code"`
	CharCode  string     `xml:"CharCode" json:"char_code"`
	Nominal   int        `xml:"Nominal" json:"nominal"`
	Name      string     `xml:"Name" json:"name"`
	Value     FloatComma `xml:"Value" json:"value"`
	VunitRate FloatComma `xml:"VunitRate" json:"vunit_rate"`
}

func LoadCurrencies(path string) ([]Valute, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	sort.Slice(valCurs.Valutes, func(i, j int) bool {
		return valCurs.Valutes[i].Value > valCurs.Valutes[j].Value
	})

	return valCurs.Valutes, nil
}

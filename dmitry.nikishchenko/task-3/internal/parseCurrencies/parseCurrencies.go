package parsecurrencies

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

func (floatField *FloatComma) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var inputField string
	if err := d.DecodeElement(&inputField, &start); err != nil {
		return fmt.Errorf("failed to decode element: %w", err)
	}

	if inputField == "" {
		*floatField = 0.0

		return nil
	}

	inputField = strings.ReplaceAll(inputField, ",", ".")

	v, err := strconv.ParseFloat(inputField, 64)
	if err != nil {
		return fmt.Errorf("failed to parse float: %w", err)
	}

	*floatField = FloatComma(v)

	return nil
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID        string     `json:"id"         xml:"ID,attr"`
	NumCode   int        `json:"num_code"   xml:"NumCode"`
	CharCode  string     `json:"char_code"  xml:"CharCode"`
	Nominal   int        `json:"nominal"    xml:"Nominal"`
	Name      string     `json:"name"       xml:"Name"`
	Value     FloatComma `json:"value"      xml:"Value"`
	VunitRate FloatComma `json:"vunit_rate" xml:"VunitRate"`
}

func LoadCurrencies(path string) ([]Valute, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("cannot close file: %v", err)
		}
	}()

	decoder := xml.NewDecoder(file)
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

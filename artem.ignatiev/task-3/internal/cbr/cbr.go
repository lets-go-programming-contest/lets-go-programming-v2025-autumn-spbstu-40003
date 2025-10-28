package cbr

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

func ParseCBR(path string) ([]Valute, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("cannot parse xml: %w", err)
	}

	for i := range valCurs.Valutes {
		v := valCurs.Valutes[i]
		if v.NumCode == "" {
			return nil, fmt.Errorf("invalid NumCode for element %d", i)
		}

		if _, err := strconv.ParseFloat(replaceComma(v.Value), 64); err != nil {
			return nil, fmt.Errorf("invalid Value for %s: %w", v.NumCode, err)
		}
	}

	return valCurs.Valutes, nil
}

func replaceComma(s string) string {
	result := ""
	for _, c := range s {
		if c == ',' {
			result += "."
		} else {
			result += string(c)
		}
	}
	return result
}

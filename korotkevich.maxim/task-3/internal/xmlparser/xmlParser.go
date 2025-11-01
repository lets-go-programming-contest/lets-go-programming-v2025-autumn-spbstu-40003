package xmlparser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

var (
	ErrNumCode      = errors.New("error: invalid num code")
	ErrEmptyValue   = errors.New("error: empty Value")
	ErrInvalidValue = errors.New("error: invalid value")
)

type CustomFloat float64

type ExchangeTrade struct {
	NumCode  int         `json:"num_code"  xml:"NumCode"`
	CharCode string      `json:"char_code" xml:"CharCode"`
	Value    CustomFloat `json:"value"     xml:"Value"`
}

type ExchangeData struct {
	Valutes []ExchangeTrade `xml:"Valute"`
}

func (f *CustomFloat) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var val string

	if err := dec.DecodeElement(&val, &start); err != nil {
		return fmt.Errorf("error with decoding element: %w", err)
	}

	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, ",", ".")

	if val == "" {
		return ErrEmptyValue
	}

	valFLoat, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return ErrInvalidValue
	}

	*f = CustomFloat(valFLoat)

	return nil
}

func XMLParse(path string) ([]ExchangeTrade, error) {
	var exData ExchangeData

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open XML file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "error with closing XML file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&exData); err != nil {
		return nil, fmt.Errorf("parse XML: %w", err)
	}

	sort.Slice(exData.Valutes, func(i, j int) bool {
		return exData.Valutes[i].Value > exData.Valutes[j].Value
	})

	return exData.Valutes, nil
}

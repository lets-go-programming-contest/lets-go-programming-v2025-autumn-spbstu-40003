package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (v *CurrencyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valueStr string

	if err := d.DecodeElement(&valueStr, &start); err != nil {
		return fmt.Errorf("decode value: %w", err)
	}

	cleanValue := strings.Replace(valueStr, ",", ".", 1)

	parsedValue, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		return fmt.Errorf("parse value %q: %w", valueStr, err)
	}

	*v = CurrencyValue(parsedValue)

	return nil
}

type Currency struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}

type ValCurs struct {
	XMLName xml.Name   `xml:"ValCurs"`
	Items   []Currency `xml:"Valute"`
}

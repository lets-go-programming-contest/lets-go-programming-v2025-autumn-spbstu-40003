package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type ValCurs struct {
	XMLName xml.Name   `xml:"ValCurs"`
	Items   []Currency `xml:"Valute"`
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type temp struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var tempData temp

	if err := decoder.DecodeElement(&tempData, &start); err != nil {
		return fmt.Errorf("decode XML element: %w", err)
	}

	c.NumCode, _ = strconv.Atoi(strings.TrimSpace(tempData.NumCode))
	c.CharCode = strings.TrimSpace(tempData.CharCode)

	cleanValue := strings.Replace(tempData.Value, ",", ".", 1)
	c.Value, _ = strconv.ParseFloat(cleanValue, 64)

	return nil
}

package data

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type DotFloat float64

type Valute struct {
	NumCode  int      `json:"num_code" xml:"NumCode"`
	CharCode string   `json:"char_code" xml:"CharCode"`
	Value    DotFloat `json:"value" xml:"Value"`
}

type Valutes struct {
	AllValutes []Valute `xml:"Valute"`
}

func (dotFloat *DotFloat) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var str string

	err := decoder.DecodeElement(&str, &start)
	if err != nil {
		return fmt.Errorf("failed decoding element: %w", err)
	}

	str = strings.ReplaceAll(str, ",", ".")

	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("failed parsing float: %w", err)
	}

	*dotFloat = DotFloat(num)

	return nil
}

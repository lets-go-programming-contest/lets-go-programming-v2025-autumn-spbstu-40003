package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ExchangeData struct {
	XMLName    xml.Name `xml:"ValCurs"`
	DateValue  string   `xml:"Date,attr"`
	SourceName string   `xml:"name,attr"`
	Items      []Item   `xml:"Valute"`
}

type moneyValue float64

type Item struct {
	NumCode   int        `json:"num_code"  xml:"NumCode"`
	CharCode  string     `json:"char_code" xml:"CharCode"`
	RateValue moneyValue `json:"value"     xml:"Value"`
}

func (mv *moneyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var text string
	if err := d.DecodeElement(&text, &start); err != nil {
		return fmt.Errorf("value decoding error: %w", err)
	}

	if text == "" {
		*mv = moneyValue(0)

		return nil
	}

	text = strings.ReplaceAll(text, ",", ".")

	num, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf("value parsing error '%s': %w", text, err)
	}

	*mv = moneyValue(num)

	return nil
}

package cbr

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type XMLFloat float64

func (f *XMLFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = XMLFloat(v)
	return nil
}

type ValuteXML struct {
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  int      `xml:"Nominal"`
	ValueRaw XMLFloat `xml:"Value"`
}

type ValCursXML struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Valutes []ValuteXML `xml:"Valute"`
}

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ParseXML(r io.Reader) ([]Currency, error) {
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReaderLabel

	var vc ValCursXML
	if err := dec.Decode(&vc); err != nil {
		return nil, err
	}

	out := make([]Currency, 0, len(vc.Valutes))
	for _, v := range vc.Valutes {
		num := 0
		if strings.TrimSpace(v.NumCode) != "" {
			var err error
			num, err = strconv.Atoi(v.NumCode)
			if err != nil {
				return nil, fmt.Errorf("invalid NumCode %q: %w", v.NumCode, err)
			}
		}

		valuePerUnit := float64(v.ValueRaw)
		if v.Nominal > 0 {
			valuePerUnit /= float64(v.Nominal)
		}

		out = append(out, Currency{
			NumCode:  num,
			CharCode: v.CharCode,
			Value:    valuePerUnit,
		})
	}

	return out, nil
}

package cbr

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type valCurs struct {
	XMLName xml.Name  `xml:"ValCurs"`
	Valutes []valuteX `xml:"Valute"`
}

type valuteX struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	ValueRaw string `xml:"Value"`
}

func ParseXML(r io.Reader) ([]Currency, error) {
	var vc valCurs
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&vc); err != nil {
		return nil, err
	}

	out := make([]Currency, 0, len(vc.Valutes))
	for _, v := range vc.Valutes {
		num, err := strconv.Atoi(strings.TrimSpace(v.NumCode))
		if err != nil {
			return nil, fmt.Errorf("invalid NumCode %q: %w", v.NumCode, err)
		}
		char := strings.TrimSpace(v.CharCode)

		nominal, err := strconv.Atoi(strings.TrimSpace(v.Nominal))
		if err != nil {
			return nil, fmt.Errorf("invalid Nominal %q for %s: %w", v.Nominal, char, err)
		}
		valStr := strings.ReplaceAll(strings.TrimSpace(v.ValueRaw), ",", ".")
		valF, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Value %q for %s: %w", v.ValueRaw, char, err)
		}
		valuePerUnit := valF / float64(nominal)

		out = append(out, Currency{
			NumCode:  num,
			CharCode: char,
			Value:    valuePerUnit,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})

	return out, nil
}

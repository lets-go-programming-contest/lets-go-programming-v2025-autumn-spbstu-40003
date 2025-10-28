package cbr

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
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
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReaderLabel

	var vc valCurs
	if err := dec.Decode(&vc); err != nil {
		return nil, err
	}

	out := make([]Currency, 0, len(vc.Valutes))
	for _, v := range vc.Valutes {
		if strings.TrimSpace(v.NumCode) == "" || strings.TrimSpace(v.CharCode) == "" || strings.TrimSpace(v.Nominal) == "" || strings.TrimSpace(v.ValueRaw) == "" {
			continue
		}

		num, err := strconv.Atoi(strings.TrimSpace(v.NumCode))
		if err != nil {
			continue
		}
		char := strings.TrimSpace(v.CharCode)

		nominal, err := strconv.Atoi(strings.TrimSpace(v.Nominal))
		if err != nil {
			continue
		}
		valStr := strings.ReplaceAll(strings.TrimSpace(v.ValueRaw), ",", ".")
		valF, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}
		valuePerUnit := valF / float64(nominal)

		out = append(out, Currency{
			NumCode:  num,
			CharCode: char,
			Value:    valuePerUnit,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid currency entries found after parsing")
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})

	return out, nil
}

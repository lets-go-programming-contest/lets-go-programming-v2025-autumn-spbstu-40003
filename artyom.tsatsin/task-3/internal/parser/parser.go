package parser

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Currency struct {
	Code  int     `json:"num_code" xml:"NumCode"`
	Char  string  `json:"char_code" xml:"CharCode"`
	Value float64 `json:"value" xml:"Value"`
}

func (c *Currency) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic("XML token error: " + err.Error())
		}

		switch el := token.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "NumCode":
				var s string
				if err := dec.DecodeElement(&s, &el); err != nil {
					panic("Failed to decode NumCode: " + err.Error())
				}
				c.Code, _ = strconv.Atoi(strings.TrimSpace(s))
			case "CharCode":
				var s string
				if err := dec.DecodeElement(&s, &el); err != nil {
					panic("Failed to decode CharCode: " + err.Error())
				}
				c.Char = strings.TrimSpace(s)
			case "Value":
				var s string
				if err := dec.DecodeElement(&s, &el); err != nil {
					panic("Failed to decode Value: " + err.Error())
				}
				s = strings.ReplaceAll(strings.TrimSpace(s), ",", ".")
				val, err := strconv.ParseFloat(s, 64)
				if err != nil {
					panic("Invalid Value: " + err.Error())
				}
				c.Value = val
			}
		case xml.EndElement:
			if el.Name.Local == start.Name.Local {
				return nil
			}
		}
	}
	return nil
}

func LoadXML(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		panic("Failed to open XML file: " + err.Error())
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var list []Currency
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic("XML reading error: " + err.Error())
		}

		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "Valute" {
			continue
		}

		var curr Currency
		if err := decoder.DecodeElement(&curr, &start); err != nil {
			panic("Currency decode error: " + err.Error())
		}

		list = append(list, curr)
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Value > list[j].Value
	})

	return list, nil
}

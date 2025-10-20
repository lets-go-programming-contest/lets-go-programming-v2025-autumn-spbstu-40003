package xmlparser

import (
	"encoding/xml"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ExchangeTrade struct {
	NumCode  int     `json:"num_code" xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value" xml:"Value"`
}

func (e *ExchangeTrade) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type exchangeTradeXML struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var xmlData exchangeTradeXML
	if err := d.DecodeElement(&xmlData, &start); err != nil {
		return err
	}

	if strings.TrimSpace(xmlData.NumCode) != "" {
		e.NumCode, _ = strconv.Atoi(strings.TrimSpace(xmlData.NumCode))
	}

	if strings.TrimSpace(xmlData.Value) != "" {
		valueStr := strings.ReplaceAll(xmlData.Value, ",", ".")
		e.Value, _ = strconv.ParseFloat(valueStr, 64)
	}

	e.CharCode = strings.TrimSpace(xmlData.CharCode)

	return nil
}

type ExchangeData struct {
	Trades []ExchangeTrade `xml:"Valute"`
}

func ParseXML(path string) ([]ExchangeTrade, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data ExchangeData
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	sort.Slice(data.Trades, func(i, j int) bool {
		return data.Trades[i].Value > data.Trades[j].Value
	})

	return data.Trades, nil
}

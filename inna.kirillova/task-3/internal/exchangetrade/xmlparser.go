package exchangetrade

import (
	"encoding/xml"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ExchangeTrade struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type xmlData struct {
	Items []struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

func parseXML(path string) []ExchangeTrade {
	file, err := os.Open(path)
	if err != nil {
		panic("XML error: " + err.Error())
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data xmlData
	if err := decoder.Decode(&data); err != nil {
		panic("bad XML: " + err.Error())
	}

	var trades []ExchangeTrade
	for _, item := range data.Items {
		numCode, _ := strconv.Atoi(strings.TrimSpace(item.NumCode))
		value, err := strconv.ParseFloat(strings.Replace(item.Value, ",", ".", -1), 64)
		if err != nil {
			continue
		}

		trades = append(trades, ExchangeTrade{
			NumCode:  numCode,
			CharCode: strings.TrimSpace(item.CharCode),
			Value:    value,
		})
	}

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].Value > trades[j].Value
	})

	return trades
}

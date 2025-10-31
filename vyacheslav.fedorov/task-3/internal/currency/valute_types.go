package currency

type CurrencyList struct {
	Date    string         `xml:"Date"`
	Name    string         `xml:"Name"`
	Valutes []CurrencyItem `xml:"Valute"`
}

type (
	Value64      float64
	CurrencyItem struct {
		NumCode  int     `json:"num_code"  xml:"NumCode"`
		CharCode string  `json:"char_code" xml:"CharCode"`
		Value    Value64 `json:"value"     xml:"Value"`
	}
)

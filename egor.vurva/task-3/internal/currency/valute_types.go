package currency

type ValCurs struct {
	Date    string   `xml:"Date"`
	Name    string   `xml:"Name"`
	Valutes []Valute `xml:"Valute"`
}

type (
	Value64 float64
	Valute  struct {
		NumCode  int     `json:"num_code"  xml:"NumCode"`
		CharCode string  `json:"char_code" xml:"CharCode"`
		Value    Value64 `json:"value"     xml:"Value"`
	}
)

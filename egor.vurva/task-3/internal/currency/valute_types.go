package currency

type Value64 float64
type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	Date    string   `xml:"Date"`
	Name    string   `xml:"Name"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID        string  `xml:"ID"        json:"-"`
	NumCode   int     `xml:"NumCode"   json:"NumCode"`
	CharCode  string  `xml:"CharCode"  json:"CharCode"`
	Nominal   int     `xml:"Nominal"   json:"-"`
	Name      string  `xml:"Name"      json:"-"`
	Value     Value64 `xml:"Value"     json:"Value"`
	VunitRate string  `xml:"VunitRate" json:"-"`
}

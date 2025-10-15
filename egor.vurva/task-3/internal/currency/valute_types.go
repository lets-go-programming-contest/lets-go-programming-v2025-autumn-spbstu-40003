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
	NumCode   int     `xml:"NumCode"   json:"num_code"`
	CharCode  string  `xml:"CharCode"  json:"char_code"`
	Nominal   int     `xml:"Nominal"   json:"-"`
	Name      string  `xml:"Name"      json:"-"`
	Value     Value64 `xml:"Value"     json:"value"`
	VunitRate string  `xml:"VunitRate" json:"-"`
}

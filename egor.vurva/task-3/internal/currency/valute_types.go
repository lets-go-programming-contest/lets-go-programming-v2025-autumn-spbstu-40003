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
	ID        string  `json:"-"         xml:"ID"`
	NumCode   int     `json:"num_code"  xml:"NumCode"`
	CharCode  string  `json:"char_code" xml:"CharCode"`
	Nominal   int     `json:"-"         xml:"Nominal"`
	Name      string  `json:"-"         xml:"Name"`
	Value     Value64 `json:"value"     xml:"Value"`
	VunitRate string  `json:"-"         xml:"VunitRate"`
}

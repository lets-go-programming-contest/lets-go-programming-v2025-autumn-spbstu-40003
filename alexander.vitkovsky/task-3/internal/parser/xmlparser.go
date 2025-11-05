package parser

// Файл с парсингом данных из .xml

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

type Valute struct {
	ID        string `xml:"ID,attr"`
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   string `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

func ParseXML(path string) ValCurs {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	valCurs := ValCurs{}
	err = decoder.Decode(&valCurs)
	if err != nil {
		panic(err)
	}

	return valCurs
}

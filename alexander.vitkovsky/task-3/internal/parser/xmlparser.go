package parser

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
		panic("failed to open input XML file" + err.Error())
	}

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)

	if err := file.Close(); err != nil {
		panic("failed to close XML file" + err.Error())
	}

	if err != nil {
		panic("failed to decode XML file" + err.Error())
	}

	if len(valCurs.Valutes) == 0 {
		panic("empty XML or invalid structure")
	}

	return valCurs
}

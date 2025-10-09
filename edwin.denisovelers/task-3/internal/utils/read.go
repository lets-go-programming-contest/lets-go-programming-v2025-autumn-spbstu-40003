package utils

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type Exchange struct {
	Valutes []Valute `xml:"Valute"`
}

func Read(path string) (*Exchange, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var exchange Exchange
	if err := decoder.Decode(&exchange); err != nil {
		return nil, fmt.Errorf("decode input file: %w", err)
	}

	return &exchange, nil
}

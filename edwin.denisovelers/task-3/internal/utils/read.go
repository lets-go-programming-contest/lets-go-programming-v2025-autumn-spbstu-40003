package utils

import (
	"encoding/xml"
	"errors"
	"os"

	"golang.org/x/net/html/charset"
)

var (
	errInputFileRead  = errors.New("cannot read input file")
	errInputFileParse = errors.New("cannot parse input file")
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
		return nil, errInputFileRead
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	exchange := Exchange{nil}
	if err := decoder.Decode(&exchange); err != nil {
		return nil, errInputFileParse
	}

	return &exchange, nil
}

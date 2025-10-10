package parsevalutes

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
)

type Valutes struct {
	ValuteElements []ValuteElement `xml:"Valute"`
}

type ValuteElement struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func ParseInput(filePath string) (Valutes, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("error opening input file: %w", err))
	}

	defer func() {
		err := inputFile.Close()
		if err != nil {
			panic(fmt.Errorf("error closing config file: %w", err))
		}
	}()

	xmlDecoder := xml.NewDecoder(inputFile)
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	var curValute Valutes
	if err := xmlDecoder.Decode(&curValute); err != nil && errors.Is(err, io.EOF) {
		panic(fmt.Errorf("invalid signature of input file: %w", err))
	}

	return curValute, nil
}

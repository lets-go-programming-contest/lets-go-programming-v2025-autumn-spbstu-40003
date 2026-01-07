package convert

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func LoadXMLData[T any](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening XML file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("error closing file")
		}
	}()

	xmlDecoder := xml.NewDecoder(file)
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	data := new(T)
	if err := xmlDecoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding XML: %w", err)
	}

	return data, nil
}

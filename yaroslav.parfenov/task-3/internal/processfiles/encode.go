	package processfiles

	import (
		"encoding/json"
		"encoding/xml"
		"errors"
		"fmt"
		"io"
		"os"
		"path"
		"strconv"
		"strings"

		"golang.org/x/net/html/charset"
	)

	type Float64WithComma float64

	type Valutes struct {
		Valutes []Valute `xml:"Valute"`
	}

	type Valute struct {
		NumCode  int              `xml:"NumCode" json:"num_code"`
		CharCode string           `xml:"CharCode" json:"char_code"`
		Value    Float64WithComma `xml:"Value" json:"value"`
	}

	func (value *Float64WithComma) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
		var readVal string

		if err := d.DecodeElement(&readVal, &start); err != nil {
			return fmt.Errorf("error unmarshalling XML: %w", err)
		}

		readVal = strings.ReplaceAll(readVal, ",", ".")

		parsedVal, err := strconv.ParseFloat(readVal, 64)
		if err != nil {
			return fmt.Errorf("error parsing float: %w", err)
		}

		*value = Float64WithComma(parsedVal)

		return nil
	}

	func ParseInput(filePath string) (Valutes, error) {
		var curValute Valutes

		inputFile, err := os.Open(filePath)
		if err != nil {
			return curValute, fmt.Errorf("error opening input file: %w", err)
		}

		defer func() {
			err := inputFile.Close()
			if err != nil {
				panic(fmt.Errorf("error closing config file: %w", err))
			}
		}()

		xmlDecoder := xml.NewDecoder(inputFile)
		xmlDecoder.CharsetReader = charset.NewReaderLabel

		if err := xmlDecoder.Decode(&curValute); err != nil && !errors.Is(err, io.EOF) {
			return curValute, fmt.Errorf("invalid signature of input file: %w", err)
		}

		return curValute, nil
	}

	func prepareOutputFile(outputPath string) (*os.File, error) {
		const permissions = 0o755

		var outputFile *os.File

		dirPath := path.Dir(outputPath)

		if _, err := os.Stat(dirPath); err != nil {
			err := os.MkdirAll(dirPath, permissions)
			if err != nil {
				return nil, fmt.Errorf("error creating output directory: %w", err)
			}
		}

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("error creating output file: %w", err)
		}

		return outputFile, nil
	}

	func OutputEncodedValutes(outputPath string, encodedValutes []Valute) error {
		outputFile, err := prepareOutputFile(outputPath)
		if err != nil {
			return fmt.Errorf("error preparing output file: %w", err)
		}

		defer func() {
			err := outputFile.Close()
			if err != nil {
				panic(fmt.Errorf("error closing config file: %w", err))
			}
		}()

		jsonEncoder := json.NewEncoder(outputFile)
		jsonEncoder.SetIndent("", "  ")

		err = jsonEncoder.Encode(encodedValutes)
		if err != nil {
			return fmt.Errorf("error encoding valutes: %w", err)
		}

		return nil
	}

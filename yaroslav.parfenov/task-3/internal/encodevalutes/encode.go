package encodevalutes

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/gituser549/task-3/internal/parsevalutes"
)

type EncodedValutes struct {
	Valutes []EncodedValute `json:"valutes"`
}

type EncodedValute struct {
	NumCode  int64   `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func PrepareValutesForEncode(valutes *parsevalutes.Valutes) ([]EncodedValute, error) {
	encoded := make([]EncodedValute, 0, len(valutes.ValuteElements))

	for _, elem := range valutes.ValuteElements {
		numCode, err := strconv.ParseInt(elem.NumCode, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("error preparing NumCode of valute: %w", err)
		}

		value, err := strconv.ParseFloat(strings.ReplaceAll(elem.Value, ",", "."), 64)
		if err != nil {
			return nil, fmt.Errorf("error preparing Value of valute: %w", err)
		}

		encoded = append(encoded, EncodedValute{numCode, elem.CharCode, value})
	}

	sort.Slice(encoded, func(i, j int) bool { return encoded[i].Value > encoded[j].Value })

	return encoded, nil
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

func OutputEncodedValutes(outputPath string, encodedValutes []EncodedValute) error {
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

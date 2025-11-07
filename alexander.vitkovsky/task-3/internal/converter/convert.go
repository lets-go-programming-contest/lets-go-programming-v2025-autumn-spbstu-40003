package converter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alexpi3/task-3/internal/parser"
)

type ValuteResult struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ToResult(valCurs parser.ValCurs) []ValuteResult {
	valuteResults := make([]ValuteResult, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		var valuteResult ValuteResult

		valuteResult.NumCode, _ = strconv.Atoi(valute.NumCode)
		valuteResult.CharCode = valute.CharCode

		valueStr := strings.ReplaceAll(valute.Value, ",", ".")
		value, _ := strconv.ParseFloat(valueStr, 64)
		valuteResult.Value = value

		valuteResults = append(valuteResults, valuteResult)
	}

	return valuteResults
}

func SortByValueDesc(items []ValuteResult) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})
}

func SaveToJSON(path string, data []ValuteResult) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		panic("failed to create directories for output" + err.Error())
	}

	file, err := os.Create(path)
	if err != nil {
		panic("failed to create output file" + err.Error())
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		panic("failed to convert in JSON" + err.Error())
	}

	if err := file.Close(); err != nil {
		panic("failed to close output file" + err.Error())
	}
}

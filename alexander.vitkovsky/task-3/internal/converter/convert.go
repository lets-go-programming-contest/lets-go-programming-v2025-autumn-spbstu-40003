package converter

// Файл с обработкой, сортировкой и конвертацией данных в json

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/alexpi3/task-3/internal/parser"
)

type ValuteResult struct {
	/* "Дораспарсенные" данные для сортировки и вывода в json
	Только нужные поля в нужном формате */
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ToResult(valCurs parser.ValCurs) []ValuteResult {
	// "дораспарсивание"
	var (
		err           error = nil
		valuteResults []ValuteResult
	)

	for _, valute := range valCurs.Valutes {
		valuteResult := ValuteResult{}

		valuteResult.NumCode, err = strconv.Atoi(valute.NumCode)
		if err != nil {
			fmt.Println(err)
		}

		valuteResult.CharCode = valute.CharCode

		valuteResult.Value, err = strconv.ParseFloat(strings.ReplaceAll(valute.VunitRate, ",", "."), 64)
		if err != nil {
			fmt.Println(err)
		}

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
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		panic(err)
	}
}

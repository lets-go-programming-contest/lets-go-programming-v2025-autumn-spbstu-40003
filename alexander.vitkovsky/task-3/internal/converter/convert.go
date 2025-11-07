package converter

// Файл с обработкой, сортировкой и конвертацией данных в json

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
	/* "Дораспарсенные" данные для сортировки и вывода в json
	Только нужные поля в нужном формате */
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ToResult(valCurs parser.ValCurs) []ValuteResult {
	// "дораспарсивание"
	var valuteResults []ValuteResult

	// По условию входные данные всегда валидны, ошибки конвертации не отливливаю

	for _, valute := range valCurs.Valutes {
		var valuteResult ValuteResult

		valuteResult.NumCode, _ = strconv.Atoi(valute.NumCode)
		valuteResult.CharCode = valute.CharCode

		valueStr := strings.ReplaceAll(valute.Value, ",", ".")
		value, _ := strconv.ParseFloat(valueStr, 64)

		nominal, _ := strconv.Atoi(valute.Nominal)
		valuteResult.Value = value / float64(nominal)

		valuteResults = append(valuteResults, valuteResult)
	}

	return valuteResults
}

func SortByValueAsc(items []ValuteResult) {
	// компаратор для валют
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})
}

func SaveToJSON(path string, data []ValuteResult) {
	// Нужно создать директории, если их нет. Иначе ругаются тесты
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		panic("failed to create directores for output" + err.Error())
	}

	file, err := os.Create(path)
	if err != nil {
		panic("failed to ceate output file" + err.Error())
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		panic("failed to convert in JSON" + err.Error())
	}
}

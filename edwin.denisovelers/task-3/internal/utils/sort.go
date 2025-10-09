package utils

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

var errNormalize = errors.New("normalization error")

type JSONValute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func Sort(valutes []Valute) ([]JSONValute, error) {
	normalized, err := normalize(valutes)
	if err != nil {
		return nil, err
	}

	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].Value > normalized[j].Value
	})

	return normalized, nil
}

func normalize(valutes []Valute) ([]JSONValute, error) {
	normalized := make([]JSONValute, len(valutes))

	for valuteIndex := range valutes {
		value := strings.ReplaceAll(valutes[valuteIndex].Value, ",", ".")

		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, errNormalize
		}

		normalized[valuteIndex] = JSONValute{valutes[valuteIndex].NumCode, valutes[valuteIndex].CharCode, floatValue}
	}

	return normalized, nil
}

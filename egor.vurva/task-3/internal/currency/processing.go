package currency

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func SortValute(valutes []Valute) {
	slices.SortFunc(valutes, func(first, second Valute) int {
		return cmp.Compare(second.Value, first.Value)
	})
}

func (value *Value64) UnmarshalText(text []byte) error {
	strValue := strings.Replace(string(text), ",", ".", 1)

	number, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return fmt.Errorf("cannot convert %s to float64", string(text))
	}

	*value = Value64(number)

	return nil
}

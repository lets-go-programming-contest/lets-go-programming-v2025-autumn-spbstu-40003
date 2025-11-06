// task-3/internal/currency/processing.go
package currency

import (
	"cmp"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

func LoadCurrencyData(path string) (CurrencyList, error) {
	file, err := os.Open(path)
	if err != nil {
		return CurrencyList{}, fmt.Errorf("open %q: %w", path, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(cerr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var curs CurrencyList

	derr := decoder.Decode(&curs)
	if derr != nil && !errors.Is(derr, io.EOF) {
		return CurrencyList{}, fmt.Errorf("decode xml %q: %w", path, derr)
	}

	return curs, nil
}

func SortByValue(valutes []CurrencyItem) {
	slices.SortFunc(valutes, func(first, second CurrencyItem) int {
		return cmp.Compare(second.Value, first.Value)
	})
}

var ErrValue64Parse = errors.New("cannot convert to float64")

func (value *Value64) UnmarshalText(text []byte) error {
	strValue := strings.Replace(string(text), ",", ".", 1)

	number, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return fmt.Errorf("%w: %q", ErrValue64Parse, string(text))
	}

	*value = Value64(number)

	return nil
}

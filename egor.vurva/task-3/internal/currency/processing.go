package currency

import (
	"cmp"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

func ReadValCurs(path string) (curs ValCurs, err error) {
	file, err := os.Open(path)
	if err != nil {
		return ValCurs{}, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close %q: %w", path, cerr)
		}
	}()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	if derr := dec.Decode(&curs); derr != nil && derr != io.EOF {
		return ValCurs{}, fmt.Errorf("decode xml %q: %w", path, derr)
	}
	return curs, nil
}

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

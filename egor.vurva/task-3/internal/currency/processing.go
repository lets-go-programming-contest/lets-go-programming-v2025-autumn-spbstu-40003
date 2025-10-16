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

func ReadValCurs(path string) (ValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return ValCurs{}, fmt.Errorf("open %q: %w", path, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(cerr)
		}
	}()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	curs := ValCurs{}

	derr := dec.Decode(&curs)
	if derr != nil && !errors.Is(derr, io.EOF) {
		return ValCurs{}, fmt.Errorf("decode xml %q: %w", path, derr)
	}

	return curs, nil
}

func SortValute(valutes []Valute) {
	slices.SortFunc(valutes, func(first, second Valute) int {
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

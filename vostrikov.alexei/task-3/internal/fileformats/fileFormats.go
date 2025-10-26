package fileformats

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

var ErrNoValCursRoot = errors.New("no ValCurs root")

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

type Val struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Vunit    float64 `json:"value"`
}

type valCursXML struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Valutes []valuteXML `xml:"Valute"`
}

type valuteXML struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func ReadXML(xmlPATH string) ([]Val, error) {
	xmlFile, err := os.ReadFile(xmlPATH)
	if err != nil {
		panic("no such file or directory")
	}

	decoder := xml.NewDecoder(strings.NewReader(string(xmlFile)))
	decoder.CharsetReader = charset.NewReaderLabel

	var root valCursXML
	if err := decoder.Decode(&root); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	if len(root.Valutes) == 0 {
		return nil, ErrNoValCursRoot
	}

	vals := make([]Val, 0, len(root.Valutes))

	for _, v := range root.Valutes {
		vals = append(vals, Val{
			NumCode:  v.NumCode,
			CharCode: strings.TrimSpace(v.CharCode),
			Vunit:    atofComma(v.Value),
		})
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Vunit > vals[j].Vunit
	})

	return vals, nil
}

func atofComma(s string) float64 {
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", ".")
	f, _ := strconv.ParseFloat(s, 64)

	return f
}

func BuildJSON(val []Val, jsonPATH string) error {
	dir := filepath.Dir(jsonPATH)

	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("cannot create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return fmt.Errorf("error while marshaling JSON: %w", err)
	}

	if err := os.WriteFile(jsonPATH, data, filePerm); err != nil {
		return fmt.Errorf("error while writing JSON file: %w", err)
	}

	return nil
}

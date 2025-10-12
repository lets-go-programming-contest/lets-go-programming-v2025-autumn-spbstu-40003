package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

var ErrNoValCursRoot = errors.New("no ValCurs root")

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

type config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type Val struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Vunit    float64 `json:"value"`
}

func ReadYaml(configPATH string) (config, error) {
	str, err := os.ReadFile(configPATH)
	if err != nil {
		panic("no such file or directory")
	}

	var cfg config
	if err := yaml.Unmarshal(str, &cfg); err != nil {
		panic("did not find expected key")
	}

	return cfg, nil
}

func ReadXML(xmlPATH string) ([]Val, error) {
	xmlFile, err := os.Open(xmlPATH)
	if err != nil {
		panic("no such file or directory")
	}

	defer func() {
		if cerr := xmlFile.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close XML file: %v\n", cerr)
		}
	}()

	reader, err := charset.NewReader(xmlFile, "windows-1251")
	if err != nil {
		return nil, fmt.Errorf("error in charset %w", err)
	}

	doc, err := xmlquery.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("error while parsing xml %w", err)
	}

	root := xmlquery.FindOne(doc, "/ValCurs")
	if root == nil {
		return nil, ErrNoValCursRoot
	}

	valutes := xmlquery.Find(doc, "//Valute")

	vals := make([]Val, 0, len(valutes))

	for _, v := range valutes {
		numCode := atoi(text(v, "NumCode"))
		charCode := text(v, "CharCode")
		vunit := atofComma(text(v, "Value"))

		vals = append(vals, Val{numCode, charCode, vunit})
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Vunit > vals[j].Vunit
	})

	return vals, nil
}

func text(node *xmlquery.Node, relPath string) string {
	if node == nil {
		return ""
	}

	child := xmlquery.FindOne(node, relPath)
	if child == nil {
		return ""
	}

	return strings.TrimSpace(child.InnerText())
}

func atofComma(s string) float64 {
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", ".")
	f, _ := strconv.ParseFloat(s, 64)

	return f
}

func atoi(s string) int {
	s = strings.TrimSpace(s)
	f, _ := strconv.Atoi(s)

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

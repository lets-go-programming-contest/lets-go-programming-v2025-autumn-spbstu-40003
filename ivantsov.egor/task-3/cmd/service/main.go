package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type Valute struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type valCurs struct {
	Items []Valute `xml:"Valute"`
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (v *Valute) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type shadow struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}
	var sh shadow
	if err := d.DecodeElement(&sh, &start); err != nil {
		return fmt.Errorf("value parse error: %w", err)
	}
	num, err := strconv.Atoi(strings.TrimSpace(sh.NumCode))
	if err != nil {
		return fmt.Errorf("num code parse error: %w", err)
	}
	valStr := strings.ReplaceAll(strings.TrimSpace(sh.Value), ",", ".")
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fmt.Errorf("float parse error: %w", err)
	}
	v.NumCode = num
	v.CharCode = strings.TrimSpace(sh.CharCode)
	v.Value = val
	return nil
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to YAML config file")
	flag.Parse()
	if configPath == "" {
		panic("use -config <path>")
	}

	data, err := os.ReadFile(configPath)
	must(err)

	var cfg Config
	must(yaml.Unmarshal(data, &cfg))

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		panic("input or output path missing")
	}

	xmlBytes, err := os.ReadFile(cfg.InputFile)
	must(err)

	var root valCurs
	decoder := xml.NewDecoder(bytes.NewReader(xmlBytes))
	decoder.CharsetReader = charset.NewReaderLabel
	must(decoder.Decode(&root))

	sort.Slice(root.Items, func(i, j int) bool {
		return root.Items[i].Value > root.Items[j].Value
	})

	dir := filepath.Dir(cfg.OutputFile)
	must(os.MkdirAll(dir, 0o755))

	out, err := os.Create(cfg.OutputFile)
	must(err)
	defer func() { _ = out.Close() }()

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	must(enc.Encode(root.Items))
}

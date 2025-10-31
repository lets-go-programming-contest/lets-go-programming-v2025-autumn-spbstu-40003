package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type Valute struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

func (v *Valute) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type shadow struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}
	var sh shadow
	if err := d.DecodeElement(&sh, &start); err != nil {

		return err
	}
	num, err := strconv.Atoi(strings.TrimSpace(sh.NumCode))
	if err != nil {

		return fmt.Errorf("NumCode parse error: %w", err)
	}
	valStr := strings.ReplaceAll(strings.TrimSpace(sh.Value), ",", ".")
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fmt.Errorf("Value parse error: %w", err)
	}
	v.NumCode = num
	v.CharCode = strings.TrimSpace(sh.CharCode)
	v.Value = val

	return nil
}

type valCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Items   []Valute `xml:"Valute"`
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func openExisting(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	return f
}

func createAll(path string) *os.File {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		must(os.MkdirAll(dir, 0o755))
	}
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	return f
}

func decodeYAMLConfig(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		panic(err)
	}
	if cfg.InputFile == "" || cfg.OutputFile == "" {
		panic("invalid config")
	}

	return cfg
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to YAML config file")
	flag.Parse()
	if configPath == "" {
		panic("use -config <path>")
	}
	cfg := decodeYAMLConfig(configPath)
	in := openExisting(cfg.InputFile)
	defer in.Close()
	xmlBytes, err := io.ReadAll(in)
	must(err)
	var root valCurs
	must(xml.Unmarshal(xmlBytes, &root))
	sort.Slice(root.Items, func(i, j int) bool {

		return root.Items[i].Value > root.Items[j].Value
	})
	out := createAll(cfg.OutputFile)
	defer out.Close()
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	must(enc.Encode(root.Items))
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kryjkaqq/task-3/internal/cbr"
	"github.com/kryjkaqq/task-3/internal/config"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		fmt.Printf("cannot load config: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Open(cfg.InputFile)
	if err != nil {
		fmt.Printf("cannot open XML file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	currencies, err := cbr.ParseXML(f)
	if err != nil {
		fmt.Printf("cannot parse XML: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		fmt.Printf("cannot create JSON file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(currencies); err != nil {
		fmt.Printf("cannot encode JSON: %v\n", err)
		os.Exit(1)
	}
}

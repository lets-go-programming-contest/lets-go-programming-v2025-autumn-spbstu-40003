package main

import (
	"flag"
	"fmt"

	"github.com/Egor1726/task-3/internal/config"
	"github.com/Egor1726/task-3/internal/parser"
	"github.com/Egor1726/task-3/internal/writer"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	currencies, err := parser.ParseFile(cfg.InputFile)
	if err != nil {
		panic(fmt.Errorf("failed to parse XML: %w", err))
	}

	err = writer.WriteJSONToFile(cfg.OutputFile, currencies)
	if err != nil {
		panic(fmt.Errorf("failed to write JSON: %w", err))
	}

	fmt.Printf("Result saved to %s\n", cfg.OutputFile)
}

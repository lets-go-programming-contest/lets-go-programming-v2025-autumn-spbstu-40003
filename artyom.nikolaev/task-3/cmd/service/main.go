package main

import (
	"flag"
	"fmt"

	"github.com/ArtttNik/task-3/internal/config"
	"github.com/ArtttNik/task-3/internal/parser"
	"github.com/ArtttNik/task-3/internal/writer"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	currencies, err := parser.ParseFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	err = writer.WriteJSONToFile(cfg.OutputFile, currencies)
	if err != nil {
		panic(fmt.Errorf("unable to write output JSON: %w", err))
	}

	fmt.Printf("Result saved to %s\n", cfg.OutputFile)
}

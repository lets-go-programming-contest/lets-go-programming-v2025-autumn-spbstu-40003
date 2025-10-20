package main

import (
	"flag"
	"fmt"

	"github.com/kirinnah/task-3/internal/config"
	"github.com/kirinnah/task-3/internal/jsonwriter"
	"github.com/kirinnah/task-3/internal/xmlparser"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	trades, err := xmlparser.ParseXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	err = jsonwriter.SaveJSON(cfg.OutputFile, trades)
	if err != nil {
		panic(fmt.Errorf("can't write output JSON: %w", err))
	}
}

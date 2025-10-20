package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kirinnah/task-3/internal/config"
	"github.com/kirinnah/task-3/internal/jsonwriter"
	"github.com/kirinnah/task-3/internal/xmlparser"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	trades, err := xmlparser.ParseXML(cfg.InputFile)
	if err != nil {
		fmt.Printf("XML parsing error: %v\n", err)
		os.Exit(1)
	}

	if err := jsonwriter.SaveJSON(cfg.OutputFile, trades); err != nil {
		fmt.Printf("JSON save error: %v\n", err)
		os.Exit(1)
	}
}

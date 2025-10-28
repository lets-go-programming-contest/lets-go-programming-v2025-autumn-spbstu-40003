package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kryjkaqq/task-3/internal/datawriter"
	"github.com/kryjkaqq/task-3/internal/setup"
	"github.com/kryjkaqq/task-3/internal/xmlhandler"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to YAML configuration file")
	flag.Parse()

	if configPath == "" {
		panic("missing required flag: -config")
	}

	cfg, err := setup.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Errorf("load config error: %w", err))
	}

	valutes, err := xmlhandler.LoadCurrencies(cfg.InputFile)
	if err != nil {
		panic(fmt.Errorf("read XML error: %w", err))
	}

	xmlhandler.SortDescending(valutes)

	if err := os.MkdirAll(filepath.Dir(cfg.OutputFile), os.ModePerm); err != nil {
		panic(fmt.Errorf("cannot create output directory: %w", err))
	}

	if err := datawriter.SaveAsJSON(cfg.OutputFile, valutes); err != nil {
		panic(fmt.Errorf("cannot write JSON: %w", err))
	}

}

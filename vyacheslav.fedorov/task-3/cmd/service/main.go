package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/S1avVv/task-3/internal/config"
	"github.com/S1avVv/task-3/internal/currency"
)

const (
	directoryPermission = 0o755
	filePermission      = 0o644
)

func main() {
	var configFilePath string

	flag.StringVar(&configFilePath, "config", "", "Provide config fileJSON path")

	flag.Parse()

	if configFilePath == "" {
		panic("missing required -cfg path")
	}

	cfg, err := config.ParseConfiguration(configFilePath)
	if err != nil {
		panic(err)
	}

	currencyData, err := currency.LoadCurrencyData(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	currency.SortByValue(currencyData.Valutes)

	jsonOutput, err := json.MarshalIndent(currencyData.Valutes, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Dir(cfg.OutputFile), directoryPermission)
	if err != nil {
		panic(fmt.Errorf("make parent dirs %q: %w", cfg.OutputFile, err))
	}

	outputFile, err := os.OpenFile(cfg.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePermission)
	if err != nil {
		panic(fmt.Errorf("open output file %q: %w", cfg.OutputFile, err))
	}

	_, err = outputFile.Write(jsonOutput)
	if err != nil {
		panic(err)
	}
}

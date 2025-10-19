package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Vurvaa/task-3/internal/config"
	"github.com/Vurvaa/task-3/internal/currency"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "", "Provide config fileJSON path")

	flag.Parse()

	if configFile == "" {
		panic("missing required -cfg path")
	}

	cfg, err := config.ParseConfig(configFile)
	if err != nil {
		panic(err)
	}

	valCurs, err := currency.ReadValCurs(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	currency.SortValute(valCurs.Valutes)

	valutsJSON, err := json.MarshalIndent(valCurs.Valutes, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Dir(cfg.OutputFile), dirPerm)
	if err != nil {
		panic(fmt.Errorf("make parent dirs %q: %w", cfg.OutputFile, err))
	}

	fileJSON, err := os.OpenFile(cfg.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		panic(fmt.Errorf("open output file %q: %w", cfg.OutputFile, err))
	}

	_, err = fileJSON.Write(valutsJSON)
	if err != nil {
		panic(err)
	}
}

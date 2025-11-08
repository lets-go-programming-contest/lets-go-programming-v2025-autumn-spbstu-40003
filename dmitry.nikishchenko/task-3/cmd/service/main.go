package main

import (
	encodeCurrencies "github.com/d1mene/task-3/internal/encodeCurrencies"
	parseCurrencies "github.com/d1mene/task-3/internal/parseCurrencies"
	processConfig "github.com/d1mene/task-3/internal/processConfig"
)

func main() {
	cfg, err := processConfig.LoadConfig()
	if err != nil {
		panic(err)
	}

	valutes, err := parseCurrencies.LoadCurrencies(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	saveErr := encodeCurrencies.SaveCurrencies(cfg.OutputFile, valutes)
	if saveErr != nil {
		panic(saveErr)
	}
}

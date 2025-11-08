package main

import (
	"flag"

	"github.com/badligyg/task-3/internal/config"
	"github.com/badligyg/task-3/internal/converter"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	settings, err := config.LoadSettings(*configPath)
	if err != nil {
		panic(err)
	}

	currencies, err := converter.ReadCurrencies(settings.InputFile)
	if err != nil {
		panic(err)
	}

	err = converter.WriteCurrencies(currencies, settings.OutputFile)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"fmt"

	"task-3/internal/config"
	"task-3/internal/processor"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config.yaml")
	flag.Parse()

	cfg, err := config.Read(*cfgPath)
	if err != nil {
		panic("Error reading configuration: " + err.Error())
	}

	exchangeData, err := processor.LoadExchangeData(cfg.InputFile)
	if err != nil {
		panic("XML parsing error: " + err.Error())
	}

	sortedCurrencies := processor.GetSortedCurrencies(exchangeData)

	if err := processor.ExportToJSON(cfg.OutputFile, sortedCurrencies); err != nil {
		panic("JSON write error: " + err.Error())
	}

	fmt.Printf("Data successfully saved to JSON\n")
}

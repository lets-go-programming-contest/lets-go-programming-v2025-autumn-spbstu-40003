package main

import (
	"flag"

	"TASK-3/internal/config"
	"TASK-3/internal/convert"
	"TASK-3/internal/currency"
)

const (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
	DirPerms          = 0755
)

func main() {
	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	settings, err := config.LoadSettings(*configPath)
	if err != nil {
		panic(err)
	}

	exchangeData, err := convert.LoadXMLData[currency.ExchangeData](settings.InputFile)
	if err != nil {
		panic(err)
	}

	convert.SortItemsByRate(&exchangeData.Items)

	if err := convert.SaveItemsAsJSON(exchangeData.Items, settings.OutputFile, DirPerms); err != nil {
		panic(err)
	}
}

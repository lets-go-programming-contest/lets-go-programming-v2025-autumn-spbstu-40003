package main

import (
	"flag"
	"fmt"

	"github.com/kryjkaqq/task-3/internal/config"
	"github.com/kryjkaqq/task-3/internal/jsonwriter"
	"github.com/kryjkaqq/task-3/internal/xmlparser"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}

	currencies, err := xmlparser.ParseXML(cfg.InputPath)
	if err != nil {
		panic(err)
	}

	if err := jsonwriter.Save(cfg.OutputPath, currencies); err != nil {
		panic(fmt.Errorf("cannot save JSON: %w", err))
	}
}

package main

import (
	"flag"

	"github.com/gituser549/task-3/internal/processconfig"
	"github.com/gituser549/task-3/internal/processfiles"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "config file")
	flag.Parse()

	cfg, err := processconfig.GetConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	valutes, err := processfiles.ParseInput(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	err = processfiles.OutputEncodedValutes(cfg.OutputFile, valutes.Valutes)
	if err != nil {
		panic(err)
	}
}

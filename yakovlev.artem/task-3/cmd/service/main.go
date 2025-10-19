package main

import (
	"flag"

	"github.com/nxgmwv/task-3/internal/config"
	"github.com/nxgmwv/task-3/internal/parser"
	"github.com/nxgmwv/task-3/internal/utils"
)

func main() {
	cfgPath := flag.String("config", "", "path to YAML config")
	flag.Parse()

	if *cfgPath == "" {
		panic("config flag must be set: -config=/path/to/config.yaml")
	}

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	data, err := parser.ParseCBR(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	if err := utils.WriteJSONToFile(cfg.OutputFile, data); err != nil {
		panic(err)
	}
}

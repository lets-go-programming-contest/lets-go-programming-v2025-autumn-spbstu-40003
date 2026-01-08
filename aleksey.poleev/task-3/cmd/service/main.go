package main

import (
	"flag"
	"fmt"

	"github.com/out1ow/task-3/internal/config"
	"github.com/out1ow/task-3/internal/parser"
	"github.com/out1ow/task-3/internal/writer"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	currency, err := parser.ParseFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	err = writer.WriteJSONToFile(cfg.OutputFile, currency)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", cfg.OutputFile)
}

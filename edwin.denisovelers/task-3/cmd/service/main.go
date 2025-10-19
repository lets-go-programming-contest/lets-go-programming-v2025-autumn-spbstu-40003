package main

import (
	"flag"

	"github.com/wedwincode/task-3/internal/config"
	"github.com/wedwincode/task-3/internal/reader"
	"github.com/wedwincode/task-3/internal/sorter"
	"github.com/wedwincode/task-3/internal/writer"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.ParseConfig(*configPath)
	if err != nil {
		panic(err)
	}

	exchange, err := reader.Read(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sorted := sorter.Sort(exchange.Valutes)

	if err := writer.Save(sorted, cfg.OutputFile); err != nil {
		panic(err)
	}
}

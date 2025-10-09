package main

import (
	"errors"
	"flag"

	"github.com/wedwincode/task-3/internal/utils"
)

var errNoConfig = errors.New("no config path is provided")

func main() {
	config := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	if *config == "" {
		panic(errNoConfig)
	}

	cfg, err := utils.ParseConfig(*config)
	if err != nil {
		panic(err)
	}

	exchange, err := utils.Read(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sorted, err := utils.Sort(exchange.Valutes)
	if err != nil {
		panic(err)
	}

	if err := utils.Save(sorted, cfg.OutputFile); err != nil {
		panic(err)
	}
}

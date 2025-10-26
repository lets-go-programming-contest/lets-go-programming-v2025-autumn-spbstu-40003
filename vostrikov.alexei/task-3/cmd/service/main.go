package main

import (
	"flag"
	"fmt"

	"github.com/lexachanskii/task-3/internal/config"
	"github.com/lexachanskii/task-3/internal/currency"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML configuration file")

	flag.Parse()

	cfg, err := config.GetConfig(*configPath)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	val, err := currency.GetValues(cfg.InputFile)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	err = currency.WriteValues(val, cfg.OutputFile)
	if err != nil {
		fmt.Println(err.Error())

		return
	}
}

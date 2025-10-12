package main

import (
	"flag"
	"fmt"

	"github.com/lexachanskii/task-3/internal/currency"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML configuration file")

	flag.Parse()

	val, err := currency.GetValues(*configPath)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	err = currency.WriteValues(val, *configPath)
	if err != nil {
		fmt.Println(err.Error())

		return
	}
}

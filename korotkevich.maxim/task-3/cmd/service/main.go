package main

import (
	"flag"
	"fmt"

	"github.com/KrrMaxim/task-3/internal/config"
	"github.com/KrrMaxim/task-3/internal/jsonwriter"
	"github.com/KrrMaxim/task-3/internal/xmlparser"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to YAML config")
	flag.Parse()

	config, err := config.ConfigLoad(*cfgPath)
	if err != nil {
		panic(err)
	}
	
	valutes, err := xmlparser.XMLParse(config.InputFile)
	if err != nil {
		panic(err)
	}

	if err := jsonwriter.WriteJSON(config.OutputFile, valutes); err != nil {
		panic(fmt.Errorf("error: problem with writing output JSON: %w", err))
	}

	fmt.Printf("Data Saved in JSON file")
}

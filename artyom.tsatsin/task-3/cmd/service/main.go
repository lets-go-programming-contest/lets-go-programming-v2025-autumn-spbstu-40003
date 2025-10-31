package main

import (
	"flag"
	"fmt"

	"github.com/Artem-Hack/task-3/internal/config"
	"github.com/Artem-Hack/task-3/internal/parser"
	"github.com/Artem-Hack/task-3/internal/writer"
)

func main() {
	cfgPath := flag.String("cfg", "config.yaml", "path to config.yaml")
	flag.Parse()

	cfg, err := config.Read(*cfgPath)
	if err != nil {
		panic("Error reading configuration: " + err.Error())
	}

	valutes, err := parser.LoadXML(cfg.InputFile)
	if err != nil {
		panic("XML parsing error: " + err.Error())
	}

	if err := writer.ExportJSON(cfg.OutputFile, valutes); err != nil {
		panic("JSON write error: " + err.Error())
	}

	fmt.Printf("Saved in Json")
}

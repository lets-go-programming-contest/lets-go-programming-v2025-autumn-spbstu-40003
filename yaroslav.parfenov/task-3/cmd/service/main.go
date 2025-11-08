package main

import (
	"flag"
	"sort"

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

	sort.Slice(valutes.Valutes, func(i, j int) bool { return valutes.Valutes[i].Value > valutes.Valutes[j].Value })

	err = processfiles.OutputEncodedValutes(cfg.OutputFile, valutes.Valutes)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"flag"

	"github.com/gituser549/task-3/internal/encodevalutes"
	"github.com/gituser549/task-3/internal/parsevalutes"
	"github.com/gituser549/task-3/internal/processconfig"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "config file")
	flag.Parse()

	cfg, err := processconfig.GetConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	valutes := parsevalutes.ParseInput(cfg.InputFile)

	preparedValutes, err := encodevalutes.PrepareValutesForEncode(&valutes)
	if err != nil {
		panic(err)
	}

	err = encodevalutes.OutputEncodedValutes(cfg.OutputFile, preparedValutes)
	if err != nil {
		panic(err)
	}
}

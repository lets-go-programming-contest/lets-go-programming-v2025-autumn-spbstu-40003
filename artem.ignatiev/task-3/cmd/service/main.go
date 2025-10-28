package main

import (
	"fmt"
	"os"

	"github.com/kryjkaqq/task-3/internal/cbr"
	"github.com/kryjkaqq/task-3/internal/config"
	"github.com/kryjkaqq/task-3/internal/jsonout"
)

func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read config: %v\n", err)
		os.Exit(1)
	}

	cbrData, err := cbr.ParseCBR(cfg.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := jsonout.WriteJSON(cfg.OutputFile, cbrData); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kryjkaqq/task-3/internal/cbr"
	"github.com/kryjkaqq/task-3/internal/config"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(fmt.Errorf("cannot load config %q: %w", *cfgPath, err))
	}

	inf, err := os.Open(cfg.InputFile)
	if err != nil {
		panic(fmt.Errorf("cannot open input file %q: %w", cfg.InputFile, err))
	}
	defer inf.Close()

	currs, err := cbr.ParseXML(inf)
	if err != nil {
		panic(fmt.Errorf("cannot parse xml %q: %w", cfg.InputFile, err))
	}

	outDir := filepath.Dir(cfg.OutputFile)
	if outDir != "." && outDir != "" {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			panic(fmt.Errorf("cannot create output directory %q: %w", outDir, err))
		}
	}

	outf, err := os.Create(cfg.OutputFile)
	if err != nil {
		panic(fmt.Errorf("cannot create output file %q: %w", cfg.OutputFile, err))
	}
	defer func() {
		_ = outf.Close()
	}()

	enc := json.NewEncoder(outf)
	enc.SetIndent("", "    ")
	if err := enc.Encode(currs); err != nil {
		panic(fmt.Errorf("cannot write json to %q: %w", cfg.OutputFile, err))
	}

	fmt.Printf("OK: %d currencies written to %s\n", len(currs), cfg.OutputFile)
}

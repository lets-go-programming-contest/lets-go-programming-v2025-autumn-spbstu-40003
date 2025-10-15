package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/Vurvaa/task-3/internal/currency"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "", "Provide config fileJSON path")

	flag.Parse()

	if configFile == "" {
		panic("missing required -config path")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		panic(fmt.Errorf("read config %q: %w", configFile, err))
	}

	var config currency.Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("unmarshal config %q: %v", configFile, err)

		return
	}

	valCurs, err := currency.ReadValCurs(config.InputFile)
	if err != nil {
		fmt.Printf("read config %q: %v", configFile, err)
	}

	currency.SortValute(valCurs.Valutes)

	valutsJSON, err := json.MarshalIndent(valCurs.Valutes, "", "  ")
	if err != nil {
		fmt.Printf("marshal config %q: %v", config.InputFile, err)

		return
	}

	err = os.MkdirAll(filepath.Dir(config.OutputFile), 0755)
	if err != nil {
		panic(fmt.Errorf("make parent dirs %q: %w", config.OutputFile, err))
	}

	fileJSON, err := os.OpenFile(config.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(fmt.Errorf("open output file %q: %w", config.OutputFile, err))
	}

	_, err = fileJSON.Write(valutsJSON)
	if err != nil {
		fmt.Printf("write config %q: %v", config.OutputFile, err)

		return
	}
}

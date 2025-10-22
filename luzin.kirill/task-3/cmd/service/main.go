package main

import (
	"github.com/KiRy6A/task-3/internal/configreader"
	"github.com/KiRy6A/task-3/internal/data"
	"github.com/KiRy6A/task-3/internal/valutemanager"
)

func main() {
	config, err := configreader.Parse()
	if err != nil {
		panic("Config parsing error: " + err.Error())
	}

	valutes, err := valutemanager.Read(config.Input)
	if err != nil {
		panic("Reading error: " + err.Error())
	}

	valutes.AllValutes = data.Sort(valutes.AllValutes)

	err = valutemanager.Write(config.Output, valutes)
	if err != nil {
		panic("Writing error: " + err.Error())
	}
}

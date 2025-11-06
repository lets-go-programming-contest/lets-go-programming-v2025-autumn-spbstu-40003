package main

import (
	"github.com/alexpi3/task-3/internal/config"
	"github.com/alexpi3/task-3/internal/converter"
	"github.com/alexpi3/task-3/internal/parser"
)

func main() {
	inputPath, outputPath := config.ParseConfig()
	valCurs := parser.ParseXML(inputPath)
	results := converter.ToResult(valCurs)

	converter.SortByValueDesc(results)
	converter.SaveToJSON(outputPath, results)
}

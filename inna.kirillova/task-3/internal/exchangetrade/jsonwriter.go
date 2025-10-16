package exchangetrade

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func saveJSON(path string, data []ExchangeTrade) {
	if len(data) == 0 {
		panic("error: no data to save to JSON file")
	}

	dir := filepath.Dir(path)
	const dirPerm = 0o755
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		panic("error while creating output directory: " + err.Error())
	}

	file, err := os.Create(path)
	if err != nil {
		panic("error while creating output file: " + err.Error())
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("error while closing file: " + err.Error())
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		panic("error while encoding data to JSON: " + err.Error())
	}

	fmt.Printf("successfully saved %d exchange trades to %s\n", len(data), path)
}

func Process(configPath string) {
	config := loadConfig(configPath)
	trades := parseXML(config.InputFile)
	saveJSON(config.OutputFile, trades)

	fmt.Printf("processing completed: %d exchange trades sorted and saved\n", len(trades))
}

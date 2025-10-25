package currency

import (
	"fmt"

	"github.com/lexachanskii/task-3/internal/fileformats"
)

func GetValues(inputPath string) ([]fileformats.Val, error) {
	res, err := fileformats.ReadXML(inputPath)
	if err != nil {
		return nil, fmt.Errorf("error in reading xml %w", err)
	}

	return res, nil
}

func WriteValues(val []fileformats.Val, outputPath string) error {
	err := fileformats.BuildJSON(val, outputPath)
	if err != nil {
		return fmt.Errorf("writevalues: %w", err)
	}

	return nil
}

package currency

import (
	"fmt"

	"github.com/lexachanskii/task-3/internal/utils"
)

func GetValues(configPath string) ([]utils.Val, error) {
	cfg, err := utils.ReadYaml(configPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading yaml %w", err)
	}

	res, err := utils.ReadXML(cfg.InputFile)
	if err != nil {
		return nil, fmt.Errorf("error in reading xml %w", err)
	}

	return res, nil
}

func WriteValues(val []utils.Val, configPath string) error {
	cfg, err := utils.ReadYaml(configPath)
	if err != nil {
		return fmt.Errorf("error while reading config %w", err)
	}

	err = utils.BuildJSON(val, cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("writevalues: %w", err)
	}

	return nil
}

package temperatureoptimizer

import (
	"fmt"
)

func ProcessDepartments(numDepartments *int) error {
	var numEmployees int

	for range *numDepartments {
		_, err := fmt.Scanln(&numEmployees)
		if err != nil {
			return fmt.Errorf("invalid number of employees: %w", err)
		}

		err = ProcessEmployees(&numEmployees)
		if err != nil {
			return err
		}
	}

	return nil
}

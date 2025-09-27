package temperatureoptimizer

import (
	"fmt"
)

const errInvNumEmployees = "invalid number of employees: %w"

func ProcessDepartments(numDepartments *int) error {
	var numEmployees int

	for range *numDepartments {
		_, err := fmt.Scanln(&numEmployees)
		if err != nil {
			return fmt.Errorf(errInvNumEmployees, err)
		}

		err = ProcessEmployees(&numEmployees)
		if err != nil {
			return err
		}
	}

	return nil
}

package temperatureoptimizer

import (
	"errors"
	"fmt"
)

var errInvNumEmployees = errors.New("Invalid number of collegues")

func ProcessDepartments(numDepartments *int) error {
	var numEmployees int

	for *numDepartments > 0 {
		_, err := fmt.Scanln(&numEmployees)
		if err != nil {
			return errInvNumEmployees
		}
		err = ProcessEmployees(&numEmployees)
		if err != nil {
			return err
		}

		(*numDepartments)--
	}

	return nil
}

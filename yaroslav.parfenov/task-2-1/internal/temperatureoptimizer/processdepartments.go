package temperatureoptimizer

import (
	"errors"
	"fmt"
)

var invNumCollegues = errors.New("Invalid number of collegues")

func ProcessDepartments(numDepartments *int) error {

	var numCollegues int

	for *numDepartments > 0 {
		_, err := fmt.Scanln(&numCollegues)
		if err != nil {
			return invNumCollegues
		}

		err = ProcessEmployees(&numCollegues)

		if err != nil {
			return err
		}

		(*numDepartments)--
	}

	return nil
}

package temperatureoptimizer

import (
	"errors"
	"fmt"
)

func ProcessDepartments(numDepartments *int) error {

	var (
		invNumCollegues = errors.New("Invalid number of collegues")
		numCollegues    int
	)

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

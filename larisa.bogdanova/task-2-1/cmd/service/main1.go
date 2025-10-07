package main

import "fmt"

const (
	MinTemperature = 15
	MaxTemperature = 30
	InvalidValue   = -1
)

func main() {
	if err := run(); err != nil {
		fmt.Println(InvalidValue)
	}
}

func run() error {
	var departmentCount int
	if _, err := fmt.Scan(&departmentCount); err != nil {
		return errReadingDepartmentCount
	}

	for range makeRange(departmentCount) {
		if err := handleDepartment(); err != nil {
			fmt.Println(InvalidValue)

			continue
		}
	}

	return nil
}

func handleDepartment() error {
	var employeeCount int
	if _, err := fmt.Scan(&employeeCount); err != nil {
		return errReadingEmployeeCount
	}

	var lowerLimit, upperLimit int = MinTemperature, MaxTemperature
	var valid bool = true

	for range makeRange(employeeCount) {
		var op string
		var temp int
		if _, err := fmt.Scan(&op, &temp); err != nil {
			return errReadingEmployee
		}

		if !valid {
			fmt.Println(InvalidValue)

			continue
		}

		switch op {
		case ">=":
			if temp > lowerLimit {
				lowerLimit = temp
			}
		case "<=":
			if temp < upperLimit {
				upperLimit = temp
			}
		default:
			return errInvalidOperator
		}

		if lowerLimit <= upperLimit {
			fmt.Println(lowerLimit)
		} else {
			fmt.Println(InvalidValue)
			valid = false
		}
	}

	return nil
}

func makeRange(n int) []struct{} {
	return make([]struct{}, n)
}

var (
	errReadingDepartmentCount = fmt.Errorf("failed to read department count")
	errReadingEmployeeCount   = fmt.Errorf("failed to read employee count")
	errReadingEmployee        = fmt.Errorf("failed to read employee input")
	errInvalidOperator        = fmt.Errorf("invalid operator")
)

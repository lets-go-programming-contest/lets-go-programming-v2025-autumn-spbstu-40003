package main

import (
	"errors"
	"fmt"
)

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

	var (
		lowerLimit, upperLimit int = MinTemperature, MaxTemperature
		valid                      = true
	)

	for range makeRange(employeeCount) {
		var (
			operator string
			temp     int
		)

		if _, err := fmt.Scan(&operator, &temp); err != nil {
			return errReadingEmployee
		}

		if !valid {
			fmt.Println(InvalidValue)

			continue
		}

		switch operator {
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
	errReadingDepartmentCount = errors.New("failed to read department count")
	errReadingEmployeeCount   = errors.New("failed to read employee count")
	errReadingEmployee        = errors.New("failed to read employee input")
	errInvalidOperator        = errors.New("invalid operator")
)

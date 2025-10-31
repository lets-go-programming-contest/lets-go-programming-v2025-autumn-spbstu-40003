package main

import (
	"errors"
	"fmt"
)

var (
	ErrEmployee    = errors.New("incorrect number of employees")
	ErrDepartment  = errors.New("incorrect number of departments")
	ErrSymbol      = errors.New("inisCorrect Symbol")
	ErrTemperature = errors.New("inisCorrect temperature")
)

const (
	minTemp     = 15
	maxTemp     = 30
	incorrValue = -1
)

func main() {
	err := execute()
	if err != nil {
		fmt.Println(incorrValue)
	}
}

func execute() error {
	var departNum int
	if _, scanErr := fmt.Scan(&departNum); scanErr != nil {
		return fmt.Errorf("%s: %w", ErrDepartment.Error(), scanErr)
	}

	for range createSlice(departNum) {
		err := handleDepartment()
		if err != nil {
			fmt.Println(incorrValue)

			continue
		}
	}

	return nil
}

func handleDepartment() error {
	var employeeNumber int
	if _, scanErr := fmt.Scan(&employeeNumber); scanErr != nil {
		return fmt.Errorf("%s: %w", ErrEmployee.Error(), scanErr)
	}

	lowerBound := minTemp
	upperBound := maxTemp
	isCorrect := true

	for range createSlice(employeeNumber) {
		var (
			symbol string
			temp   int
		)

		_, scanErr := fmt.Scan(&symbol, &temp)
		if scanErr != nil {
			return fmt.Errorf("%s: %w", ErrTemperature.Error(), scanErr)
		}

		if !isCorrect {
			fmt.Println(incorrValue)

			continue
		}

		switch symbol {
		case ">=":
			if temp > lowerBound {
				lowerBound = temp
			}
		case "<=":
			if temp < upperBound {
				upperBound = temp
			}
		default:
			return ErrSymbol
		}

		if lowerBound <= upperBound {
			fmt.Println(lowerBound)
		} else {
			fmt.Println(incorrValue)

			isCorrect = false
		}
	}

	return nil
}

func createSlice(n int) []struct{} {
	return make([]struct{}, n)
}

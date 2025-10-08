package main

import (
	"errors"
	"fmt"
)

var (
	ErrNumOfDep           = errors.New("incorrect value for number of departments")
	ErrNumOfEmployees     = errors.New("incorrect value for number of employees")
	ErrInvalidTemperature = errors.New("invalid temperature value")
	ErrInvalidOperator    = errors.New("invalid operator")
)

const (
	minTemp     = 15
	maxTemp     = 30
	invalidTemp = -1
	minValue    = 1
	maxValue    = 1000
)

func main() {
	var numDep int

	_, err := fmt.Scan(&numDep)
	if err != nil {
		fmt.Println(ErrNumOfDep)

		return
	}
	for range numDep {
		if err = processDepartment(); err != nil {
			fmt.Println(err)

			return
		}
	}
}

func processDepartment() error {
	var employees int
	
	_, err := fmt.Scan(&employees)
	if err != nil {
		return fmt.Errorf("error with reading number of employees: %w", err)
	}

	if employees < minValue || employees > maxValue {
		return fmt.Errorf("%w: %d", ErrNumOfEmployees, employees)
	}

	lowerLimit := minTemp
	upperLimit := maxTemp

	for range employees {
		var (
			operator  string
			tempValue int
		)
		
		if _, err = fmt.Scan(&operator, &tempValue); err != nil {
			return fmt.Errorf("error with reading temperature preference: %w", err)
		}

		if tempValue < minTemp || tempValue > maxTemp {
			return fmt.Errorf("%w: %d", ErrInvalidTemperature, tempValue)
		}

		switch operator {
		case "<=":
			upperLimit = Min(upperLimit, tempValue)
		case ">=":
			lowerLimit = Max(lowerLimit, tempValue)
		default:
			return fmt.Errorf("%w: %s", ErrInvalidOperator, operator)
		}

		if lowerLimit <= upperLimit {
			fmt.Println(lowerLimit)
		} else {
			fmt.Println(invalidTemp)
		}
	}

	return nil
}

func Max(firstNumber, secondNumber int) int {
	if firstNumber > secondNumber {
		return firstNumber
	}

	return secondNumber
}

func Min(firstNumber, secondNumber int) int {
	if firstNumber < secondNumber {
		return firstNumber
	}

	return secondNumber
}

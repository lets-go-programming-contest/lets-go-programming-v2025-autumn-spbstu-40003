package main

import (
	"fmt"
	"errors"
)

var (
	ErrNumOfDep = errors.New("incorrect value for number of departments")
)

const (
	minTemp = 15
	maxTemp = 30 
	invalidTemp = -1
	minValue = 1
	maxValue = 1000
)

func main() {
	var numDep int

	_, err := fmt.Scan(&numDep)
	if err != nil || numDep < minValue || numDep > maxValue {
		fmt.Println(ErrNumOfDep)
		return
	}
	
	for range numDep {
		if err := processDepartment(); err != nil {
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
		return fmt.Errorf("incorrect value for number of employees: %d", employees)
	}

	lowerLimit := minTemp
	upperLimit := maxTemp

	for j := 0; j < employees; j++ {
		var (
			operator  string
			tempValue int
		)
		
		if _, err := fmt.Scan(&operator, &tempValue); err != nil {
			return fmt.Errorf("error with reading temperature preference: %w", err)
		}

		if tempValue < minTemp || tempValue > maxTemp {
			return fmt.Errorf("invalid temperature value: %d", tempValue)
		}

		switch operator {
		case "<=":
			upperLimit = Min(upperLimit, tempValue)
		case ">=":
			lowerLimit = Max(lowerLimit, tempValue)
		default:
			return fmt.Errorf("invalid operator %q", operator)
		}

		if lowerLimit <= upperLimit {
			fmt.Println(lowerLimit)
		} else {
			fmt.Println(invalidTemp)
		}
	}

	return nil
}

func Max(FirstNumber, SecondNumber int) int {
	if FirstNumber > SecondNumber {
		return FirstNumber
	}
	return SecondNumber
}

func Min(FirstNumber, SecondNumber int) int {
	if FirstNumber < SecondNumber {
		return FirstNumber
	}
	return SecondNumber
}

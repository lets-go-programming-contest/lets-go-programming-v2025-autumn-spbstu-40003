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
	var departmentNum int
	_, err := fmt.Scan(&departmentNum)
	if err != nil {

		return ErrDepartment
	}

	departments := createSlice(departmentNum)
	for range departments {
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
	_, err := fmt.Scan(&employeeNumber)
	if err != nil {

		return ErrEmployee
	}

	lowerBound := minTemp
	upperBound := maxTemp
	isCorrect := true

	employees := createSlice(employeeNumber)
	for range employees {
		var (
			symbol string
			temp   int
		)

		_, err := fmt.Scan(&symbol, &temp)
		if err != nil {
			return ErrTemperature
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

package internal

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAmount      = errors.New("amount must be greater than zero and lower than 1000")
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrInvalidTemperature = errors.New("invalid temperature")
	ErrScan               = errors.New("scan error")
)

const (
	MaxTemp = 30
	MinTemp = 15
)

func processEmployee(op string, temp int, currentMin, currentMax *int) error {
	if temp < MinTemp || temp > MaxTemp {
		return ErrInvalidTemperature
	}

	switch op {
	case ">=":
		if temp > *currentMin {
			*currentMin = temp
		}
	case "<=":
		if temp < *currentMax {
			*currentMax = temp
		}
	default:
		return ErrInvalidOperation
	}

	return nil
}

func processDepartment(employeeCount int) error {
	if employeeCount <= 0 || employeeCount > 1000 {
		return ErrInvalidAmount
	}

	currentMin := MinTemp
	currentMax := MaxTemp

	for range employeeCount {
		var (
			operation string
			temp      int
		)

		_, err := fmt.Scan(&operation, &temp)
		if err != nil {
			return ErrScan
		}

		err = processEmployee(operation, temp, &currentMin, &currentMax)
		if err != nil {
			return err
		}

		if currentMin <= currentMax {
			fmt.Println(currentMin)
		} else {
			fmt.Println(-1)
		}
	}

	return nil
}

func Temp() error {
	var departmentCount int

	_, err := fmt.Scanln(&departmentCount)
	if err != nil {
		return ErrScan
	}

	if departmentCount <= 0 || departmentCount > 1000 {
		return ErrInvalidAmount
	}

	for range departmentCount {
		var employeeCount int

		_, err := fmt.Scanln(&employeeCount)
		if err != nil {
			return ErrScan
		}

		err = processDepartment(employeeCount)
		if err != nil {
			return err
		}
	}

	return nil
}

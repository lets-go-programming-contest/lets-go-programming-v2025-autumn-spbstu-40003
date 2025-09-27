package utils

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAmount      = errors.New("amount must be greater than zero and lower than 1000")
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrInvalidTemperature = errors.New("invalid temperature")
)

const (
	MaxTemp = 30
	MinTemp = 15
)

func processEmployee(op string, temp int, currentMin, currentMax *int) error {
	if temp < MinTemp || temp > MaxTemp {
		return fmt.Errorf("error of the range of temperature: %w", ErrInvalidTemperature)
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
		return fmt.Errorf("error while trying to read input temperature: %w", ErrInvalidOperation)
	}

	return nil
}

func processDepartment(employeeCount int) error {
	if employeeCount <= 0 || employeeCount > 1000 {
		return fmt.Errorf("error of the range of employee: %w", ErrInvalidAmount)
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
			return fmt.Errorf("error while trying to read input temperature for employee %w", err)
		}

		err = processEmployee(operation, temp, &currentMin, &currentMax)
		if err != nil {
			return fmt.Errorf("error processing employee %w", err)
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
		return fmt.Errorf("error while trying to read department count: %w", err)
	}

	if departmentCount <= 0 || departmentCount > 1000 {
		return fmt.Errorf("invalid department count %w", ErrInvalidAmount)
	}

	for range departmentCount {
		var employeeCount int

		_, err := fmt.Scanln(&employeeCount)
		if err != nil {
			return fmt.Errorf("error while trying to read employee count for department %w", err)
		}

		err = processDepartment(employeeCount)
		if err != nil {
			return fmt.Errorf("error processing department %w", err)
		}
	}

	return nil
}

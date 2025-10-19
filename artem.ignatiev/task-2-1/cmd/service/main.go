package main

import (
	"fmt"
)

const (
	minAllowed = 15
	maxAllowed = 30
	noSolution = -1
)

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println(noSolution)

		return
	}

	for range departments {
		if err := processDepartment(); err != nil {
			fmt.Println(noSolution)

			return
		}
	}
}

func processDepartment() error {
	var employees int
	if _, err := fmt.Scan(&employees); err != nil {
		return fmt.Errorf("error reading number of employees: %w", err)
	}

	low := minAllowed
	high := maxAllowed
	stillPossible := true

	for range employees {
		var operator string

		var temp int

		if _, err := fmt.Scan(&operator, &temp); err != nil {
			return fmt.Errorf("error reading temperature preference: %w", err)
		}

		if !stillPossible {
			fmt.Println(noSolution)

			continue
		}

		switch operator {
		case ">=":
			if temp > low {
				low = temp
			}
		case "<=":
			if temp < high {
				high = temp
			}
		default:
			fmt.Println(noSolution)

			stillPossible = false

			continue
		}

		if low <= high {
			fmt.Println(low)
		} else {
			fmt.Println(noSolution)

			stillPossible = false
		}
	}

	return nil
}

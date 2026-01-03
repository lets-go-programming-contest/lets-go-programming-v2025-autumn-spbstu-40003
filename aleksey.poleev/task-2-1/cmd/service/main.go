package main

import "fmt"

const (
	minAllowed = 15
	maxAllowed = 30
	noSolution = -1
)

func main() {
	var dep int
	if _, err := fmt.Scan(&dep); err != nil {
		fmt.Println(noSolution)

		return
	}

	for range dep {
		if err := tempCalc(); err != nil {
			fmt.Println(noSolution)

			return
		}
	}
}

func tempCalc() error {
	var workers int
	if _, err := fmt.Scan(&workers); err != nil {
		return fmt.Errorf("error reading number of employees: %w", err)
	}

	low := minAllowed
	high := maxAllowed
	isPos := true

	for range workers {
		var operator string
		var temp int

		if _, err := fmt.Scan(&operator, &temp); err != nil {
			return fmt.Errorf("error reading temperature preference: %w", err)
		}

		if !isPos {
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

			isPos = false

			continue
		}
		if low <= high {
			fmt.Println(low)
		} else {
			fmt.Println(noSolution)

			isPos = false
		}
	}

	return nil
}
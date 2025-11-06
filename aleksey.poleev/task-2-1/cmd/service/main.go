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
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		fmt.Println(noSolution)

		return
	}

	for range n {
		if err := tempCalculator(); err != nil {
			fmt.Println(noSolution)

			return
		}
	}
}

func tempCalculator() error {
	var k int
	if _, err := fmt.Scan(&k); err != nil {
		return fmt.Errorf("error reading number of employers: %w", err)
	}

	low := minAllowed
	high := maxAllowed
	possible := true

	for range k {
		var op string

		var t int

		if _, err := fmt.Scan(&op, &t); err != nil {
			return fmt.Errorf("error reading temperature preference: %w", err)
		}

		if !possible {
			fmt.Println(noSolution)

			continue
		}

		switch op {
		case ">=":
			if t > low {
				low = t
			}
		case "<=":
			if t < high {
				high = t
			}
		default:
			fmt.Println(noSolution)

			possible = false

			continue
		}

		if low <= high {
			fmt.Println(low)
		} else {
			fmt.Println(noSolution)

			possible = false
		}
	}

	return nil
}

// </3

package main

import "fmt"

const (
	minLimit   = 15
	maxLimit   = 30
	errorValue = -1
)

func main() {
	var depCount int

	if _, err := fmt.Scan(&depCount); err != nil {
		fmt.Println(errorValue)
		return
	}

	for dep := 0; dep < depCount; dep++ {
		if err := processDepartment(); err != nil {
			fmt.Println(errorValue)
			return
		}
	}
}

func processDepartment() error {
	var empCount int

	if _, err := fmt.Scan(&empCount); err != nil {
		return fmt.Errorf("failed to read employee count: %w", err)
	}

	lower, upper := minLimit, maxLimit
	valid := true

	for i := 0; i < empCount; i++ {
		var op string
		var temp int

		if _, err := fmt.Scan(&op, &temp); err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		if !valid {
			continue
		}

		switch op {
		case ">=":
			if temp > lower {
				lower = temp
			}
		case "<=":
			if temp < upper {
				upper = temp
			}
		}

		if lower > upper {
			fmt.Println(errorValue)
			valid = false
		}
	}

	if valid {
		fmt.Println(lower)
	}

	return nil
}

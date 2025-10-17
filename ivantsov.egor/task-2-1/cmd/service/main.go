package main

import "fmt"

const (
	minTemp  = 15
	maxTemp  = 30
	failCode = -1
)

func main() {
	var depCount int
	if _, err := fmt.Scan(&depCount); err != nil {
		fmt.Println(failCode)
		return
	}

	for i := 0; i < depCount; i++ {
		if err := evaluateDepartment(); err != nil {
			fmt.Println(failCode)
			return
		}
	}
}

func evaluateDepartment() error {
	var empCount int
	if _, err := fmt.Scan(&empCount); err != nil {
		return fmt.Errorf("cannot read employee count: %w", err)
	}

	minLimit, maxLimit := minTemp, maxTemp

	for i := 0; i < empCount; i++ {
		var sign string
		var t int

		if _, err := fmt.Scan(&sign, &t); err != nil {
			return fmt.Errorf("cannot read preference: %w", err)
		}

		switch sign {
		case ">=":
			if t > minLimit {
				minLimit = t
			}
		case "<=":
			if t < maxLimit {
				maxLimit = t
			}
		}
	}

	if minLimit <= maxLimit {
		fmt.Println(maxLimit)
	} else {
		fmt.Println(failCode)
	}

	return nil
}

package main

import "fmt"

const (
	minBase  = 15
	maxBase  = 30
	errValue = -1
)

func main() {
	var depCount int
	if _, err := fmt.Scan(&depCount); err != nil {
		fmt.Println(errValue)
		return
	}

	handleAllDepartments(depCount)
}

func handleAllDepartments(depCount int) {
	for i := 0; i < depCount; i++ {
		if err := analyzeDepartment(); err != nil {
			fmt.Println(errValue)
			return
		}
	}
}

func analyzeDepartment() error {
	var empCount int
	if _, err := fmt.Scan(&empCount); err != nil {
		return fmt.Errorf("department input error: %w", err)
	}

	minLimit := minBase
	maxLimit := maxBase
	isOk := true

	for i := 0; i < empCount; i++ {
		var sign string
		var value int

		if _, err := fmt.Scan(&sign, &value); err != nil {
			return fmt.Errorf("temperature read error: %w", err)
		}

		if !isOk {
			fmt.Println(errValue)
			continue
		}

		if sign == ">=" && value > minLimit {
			minLimit = value
		} else if sign == "<=" && value < maxLimit {
			maxLimit = value
		}

		if minLimit <= maxLimit {
			fmt.Println(maxLimit)
		} else {
			fmt.Println(errValue)
			isOk = false
		}
	}

	return nil
}

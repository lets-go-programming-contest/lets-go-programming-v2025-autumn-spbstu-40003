package main

import (
	"errors"
	"fmt"
)

var (
	ErrDep  = errors.New("incorrect number of departments")
	ErrEmp  = errors.New("incorrect number of employees")
	ErrTemp = errors.New("invalid temperature")
	ErrSign = errors.New("invalid sign")
)

const (
	minT     = 15
	maxT     = 30
	invalidT = -1
)

func main() {
	if err := run(); err != nil {
		fmt.Println(invalidT)
	}
}

func run() error {
	var depCount int
	if _, err := fmt.Scan(&depCount); err != nil {
		return ErrDep
	}

	for range loopCount(depCount) {
		if err := processDepartment(); err != nil {
			fmt.Println(invalidT)

			continue
		}
	}

	return nil
}

func processDepartment() error {
	var empCount int
	if _, err := fmt.Scan(&empCount); err != nil {
		return ErrEmp
	}

	depLow := minT
	depHigh := maxT
	valid := true

	for range loopCount(empCount) {
		var (
			sign string
			temp int
		)

		if _, err := fmt.Scan(&sign, &temp); err != nil {
			return ErrTemp
		}

		if !valid {
			fmt.Println(invalidT)

			continue
		}

		switch sign {
		case ">=":
			if temp > depLow {
				depLow = temp
			}
		case "<=":
			if temp < depHigh {
				depHigh = temp
			}
		default:
			return ErrSign
		}

		if depLow <= depHigh {
			fmt.Println(depLow)
		} else {
			fmt.Println(invalidT)
			valid = false
		}
	}

	return nil
}

func loopCount(n int) []struct{} {
	return make([]struct{}, n)
}

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
	minV     = 1
	maxV     = 1000
	invalidT = -1
)

func main() {
	var depCount int

	_, err := fmt.Scan(&depCount)
	if err != nil {
		fmt.Println(ErrDep)
		return
	}

	if depCount < minV || depCount > maxV {
		fmt.Println(invalidT)
		return
	}

	for range loopCount(depCount) {
		var empCount int

		_, err := fmt.Scan(&empCount)
		if err != nil {
			fmt.Println(invalidT)
			continue
		}

		if empCount < minV || empCount > maxV {
			fmt.Println(invalidT)
			continue
		}

		err = processDepartment(empCount)
		if err != nil {
			fmt.Println(invalidT)
		}
	}
}

func processDepartment(empCount int) error {
	if empCount < minV || empCount > maxV {
		return ErrEmp
	}

	depLow := minT
	depHigh := maxT
	valid := true

	for range loopCount(empCount) {
		var sign string
		var temp int

		_, err := fmt.Scan(&sign, &temp)
		if err != nil {
			return ErrTemp
		}

		if temp < minT || temp > maxT {
			return ErrTemp
		}

		if !valid {
			continue
		}

		switch sign {
		case "<=":
			depHigh = minInt(depHigh, temp)
		case ">=":
			depLow = maxInt(depLow, temp)
		default:
			return ErrSign
		}

		if depLow > depHigh {
			valid = false
		}
	}

	if valid {
		fmt.Println(depLow)
	} else {
		fmt.Println(invalidT)
	}

	return nil
}

func loopCount(total int) []struct{} {
	return make([]struct{}, total)
}

func maxInt(first, second int) int {
	if first > second {
		return first
	}

	return second
}

func minInt(first, second int) int {
	if first < second {
		return first
	}

	return second
}

package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Range struct {
	low, high int
}

var ErrNoOverlap = errors.New("the intervals do not overlap")

/*
	В задании нет четкого определения "оптимальной температуры".
	Считаю, что	«оптимальная температура» — это минимально возможная
	температура	в текущем допустимом диапазоне, но если старое значение
	ещё подходит, то оно сохраняется
*/

func intersection(r1 *Range, r2 *Range) error {
	if r1.high < r2.low || r1.low > r2.high {
		return ErrNoOverlap
	}

	r1.low = max(r1.low, r2.low)
	r1.high = min(r1.high, r2.high)

	return nil
}

func readInt() (int, error) {
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func processEmployee(currRange *Range, optTemp *int) error {
	var sign string
	var value int

	_, err := fmt.Scan(&sign, &value)
	if err != nil {
		return fmt.Errorf("invalid temperature range")
	}

	var newRange Range
	switch sign {
	case ">=":
		newRange = Range{value, 30}
	case "<=":
		newRange = Range{15, value}
	default:
		return fmt.Errorf("invalid sign: %s", sign)
	}

	err = intersection(currRange, &newRange)
	if err != nil {
		fmt.Println(-1)

		return nil
	}

	if *optTemp >= currRange.low && *optTemp <= currRange.high {
		fmt.Println(*optTemp)
	} else {
		*optTemp = currRange.low
		fmt.Println(*optTemp)
	}

	return nil
}

func processDepartment() error {
	emplNum, err := readInt()
	if err != nil {
		return fmt.Errorf("can`t read emplNum")
	}

	currRange := Range{15, 30}
	optTemp := -1

	for emplNum > 0 {
		if err := processEmployee(&currRange, &optTemp); err != nil {
			fmt.Println("Error:", err)

			return err
		}

		emplNum--
	}

	return nil
}

func main() {
	depNum, err := readInt()
	if err != nil {
		fmt.Println("Error: can`t read depNum")

		return
	}

	for depNum > 0 {
		if err := processDepartment(); err != nil {
			return
		}

		depNum--
	}
}

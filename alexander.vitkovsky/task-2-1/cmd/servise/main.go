package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Range struct {
	low, high int
}

var (
	ErrNoOverlap   = errors.New("the intervals do not overlap")
	ErrInvalidTemp = errors.New("invalid temperature range")
	ErrInvalidSign = errors.New("invalid sign")
	ErrReadEmplNum = errors.New("can't read emplNum")
	ErrReadDepNum  = errors.New("can't read depNum")
)

/*
	В задании нет четкого определения "оптимальной температуры".
	Считаю, что	«оптимальная температура» — это минимально возможная
	температура	в текущем допустимом диапазоне, но если старое значение
	ещё подходит, то оно сохраняется
*/

func intersection(first *Range, second *Range) error {
	if first.high < second.low || first.low > second.high {
		return ErrNoOverlap
	}

	first.low = max(first.low, second.low)
	first.high = min(first.high, second.high)

	return nil
}

func readInt() (int, error) {
	var input string

	_, err := fmt.Scan(&input)
	if err != nil {
		return 0, fmt.Errorf("read value: %w", err)
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("atoi: %w", err)
	}

	return val, nil
}

func processEmployee(currRange *Range, optTemp *int) error {
	var sign string
	var value int

	_, err := fmt.Scan(&sign, &value)
	if err != nil {
		return ErrInvalidTemp
	}

	var newRange Range

	switch sign {
	case ">=":
		newRange = Range{value, 30}
	case "<=":
		newRange = Range{15, value}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidSign, sign)
	}

	err = intersection(currRange, &newRange)
	if err != nil {
		fmt.Println(-1)

		return err
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
		return fmt.Errorf("%w: %w", ErrReadEmplNum, err)
	}

	currRange := Range{15, 30}
	optTemp := -1

	for emplNum > 0 {
		if err := processEmployee(&currRange, &optTemp); err != nil && !errors.Is(err, ErrNoOverlap) {
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
		fmt.Println("Error:", ErrReadDepNum, err)

		return
	}

	for depNum > 0 {
		if err := processDepartment(); err != nil {
			return
		}

		depNum--
	}
}

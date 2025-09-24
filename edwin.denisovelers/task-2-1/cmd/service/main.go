package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	errParse              = errors.New("parse error")
	errImpossibleSolution = errors.New("impossible solution")
)

type ConstraintType int

const (
	GreaterEqual ConstraintType = iota
	LessEqual
)

type PreferredTemperature struct {
	value int
	kind  ConstraintType
}

func parsePreferredTemperature(s string) (PreferredTemperature, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return PreferredTemperature{}, errParse
	}

	op, numStr := parts[0], parts[1]
	var constraint ConstraintType
	switch op {
	case ">=":
		constraint = GreaterEqual
	case "<=":
		constraint = LessEqual
	default:
		return PreferredTemperature{}, errParse
	}

	value, err := strconv.Atoi(numStr)
	if err != nil {
		return PreferredTemperature{}, err
	}

	return PreferredTemperature{value: value, kind: constraint}, nil
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func readInt(r *bufio.Reader) (int, error) {
	var v int
	_, err := fmt.Fscanln(r, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func applyPreference(minT, maxT int, pref PreferredTemperature) (int, int, error) {
	if pref.kind == GreaterEqual {
		if pref.value > maxT {
			return 0, 0, errImpossibleSolution
		}
		minT = max(minT, pref.value)
	} else {
		if pref.value < minT {
			return 0, 0, errImpossibleSolution
		}
		maxT = min(maxT, pref.value)
	}

	if minT > maxT {
		return 0, 0, errImpossibleSolution
	}

	return minT, maxT, nil
}

func processDepartment(r *bufio.Reader, lowerBound, upperBound int) error {
	employeeCount, err := readInt(r)
	minT, maxT := lowerBound, upperBound
	if err != nil {
		return err
	}

	for e := 0; e < employeeCount; e++ {
		rawTemperature, err := readLine(r)
		if err != nil {
			return err
		}

		pref, err := parsePreferredTemperature(rawTemperature)
		if err != nil {
			return err
		}

		minT, maxT, err = applyPreference(minT, maxT, pref)
		if err != nil {
			if errors.Is(err, errImpossibleSolution) {
				fmt.Println(-1)
				return nil
			}
			return err
		}

		fmt.Println(minT)
	}
	return nil
}

func solve(r *bufio.Reader, lowerBound int, upperBound int) error {
	departmentCount, err := readInt(r)
	if err != nil {
		return err
	}

	for d := 0; d < departmentCount; d++ {
		if err := processDepartment(r, lowerBound, upperBound); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if err := solve(in, 15, 30); err != nil {
		fmt.Println(err)
	}
}

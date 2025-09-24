package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	errParse              = errors.New("parse error")
	errImpossibleSolution = errors.New("impossible constraint")
)

type PreferredTemperature struct {
	value     int
	isGreater bool
}

func parsePreferredTemperature(s string) (PreferredTemperature, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return PreferredTemperature{}, errParse
	}

	op, numStr := parts[0], parts[1]
	var isGreater bool
	switch op {
	case ">=":
		isGreater = true
	case "<=":
		isGreater = false
	default:
		return PreferredTemperature{}, errParse
	}

	value, err := strconv.Atoi(numStr)
	if err != nil {
		return PreferredTemperature{}, err
	}

	return PreferredTemperature{value: value, isGreater: isGreater}, nil
}

func readLine(r io.Reader) (string, error) {
	line, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func readInt(r io.Reader) (int, error) {
	var v int
	_, err := fmt.Fscanln(r, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func applyPreference(minT, maxT int, pref PreferredTemperature) (int, int, error) {
	if pref.isGreater {
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

func processDepartment(r io.Reader, lowerBound, upperBound int) error {
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
			fmt.Println(-1)
			return nil
		}

		fmt.Println(minT)
	}
	return nil
}

func processDepartments(r io.Reader, lowerBound int, upperBound int) error {
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
	if err := processDepartments(in, 15, 30); err != nil {
		fmt.Println(err)
	}
}

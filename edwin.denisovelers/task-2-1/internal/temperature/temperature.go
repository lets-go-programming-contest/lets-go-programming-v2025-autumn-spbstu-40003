package temperature

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	errParse              = errors.New("parse error")
	errImpossibleSolution = errors.New("impossible solution")
)

type constraintType int

const (
	greaterEqual constraintType = iota
	lessEqual
)

const (
	DefaultLowerBound = 15
	DefaultUpperBound = 30
	expectedParts     = 2
)

type preferredTemperature struct {
	value int
	kind  constraintType
}

func parsePreferredTemperature(s string) (preferredTemperature, error) {
	parts := strings.Fields(s)
	if len(parts) != expectedParts {
		return preferredTemperature{}, errParse
	}

	op, numStr := parts[0], parts[1]

	var constraint constraintType

	switch op {
	case ">=":
		constraint = greaterEqual
	case "<=":
		constraint = lessEqual
	default:
		return preferredTemperature{}, errParse
	}

	value, err := strconv.Atoi(numStr)
	if err != nil {
		return preferredTemperature{}, fmt.Errorf("invalid number %q: %w", numStr, err)
	}

	return preferredTemperature{value: value, kind: constraint}, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error %w", err)
	}

	return strings.TrimSpace(line), nil
}

func readInt(reader *bufio.Reader) (int, error) {
	var value int

	_, err := fmt.Fscanln(reader, &value)
	if err != nil {
		return 0, fmt.Errorf("error %w", err)
	}

	return value, nil
}

func applyPreference(minT, maxT int, pref preferredTemperature) (int, int, error) {
	if pref.kind == greaterEqual {
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

func processDepartment(reader *bufio.Reader, lowerBound, upperBound int) error {
	employeeCount, err := readInt(reader)
	minT, maxT := lowerBound, upperBound
	corrupted := false

	if err != nil {
		return err
	}

	for range employeeCount {
		rawTemperature, err := readLine(reader)
		if err != nil {
			return err
		}

		if corrupted {
			fmt.Println(-1)

			continue
		}

		pref, err := parsePreferredTemperature(rawTemperature)
		if err != nil {
			return err
		}

		minT, maxT, err = applyPreference(minT, maxT, pref)
		if err != nil {
			if errors.Is(err, errImpossibleSolution) {
				fmt.Println(-1)

				corrupted = true

				continue
			}

			return err
		}

		fmt.Println(minT)
	}

	return nil
}

func Solve(reader *bufio.Reader, lowerBound int, upperBound int) error {
	departmentCount, err := readInt(reader)
	if err != nil {
		return err
	}

	for range departmentCount {
		if err := processDepartment(reader, lowerBound, upperBound); err != nil {
			continue
		}
	}

	return nil
}

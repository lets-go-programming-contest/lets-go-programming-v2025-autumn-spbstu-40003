package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	low, high int
}

var (
	ErrNoOverlap    = errors.New("the intervals do not overlap")
	ErrInvalidToken = errors.New("invalid token")
)

func intersection(first *Range, second *Range) error {
	if first.high < second.low || first.low > second.high {
		return ErrNoOverlap
	}

	if first.low < second.low {
		first.low = second.low
	}

	if first.high > second.high {
		first.high = second.high
	}

	return nil
}

func readInt(reader *bufio.Reader) (int, error) {
	var n int

	_, err := fmt.Fscan(reader, &n)
	if err != nil {
		return 0, fmt.Errorf("scan int: %w", err)
	}

	return n, nil
}

func scanConstraint(reader *bufio.Reader) (string, int, error) {
	var tok string

	if _, err := fmt.Fscan(reader, &tok); err != nil {
		return "", 0, fmt.Errorf("scan token: %w", err)
	}

	if strings.HasPrefix(tok, ">=") || strings.HasPrefix(tok, "<=") {
		sign := tok[:2]
		num := strings.TrimSpace(tok[2:])

		if num == "" {
			v, err := readInt(reader)
			if err != nil {
				return "", 0, fmt.Errorf("scan glued number: %w", err)
			}

			return sign, v, nil
		}

		v, err := strconv.Atoi(num)
		if err != nil {
			return "", 0, fmt.Errorf("atoi glued number: %w", err)
		}

		return sign, v, nil
	}

	if tok == ">=" || tok == "<=" {
		v, err := readInt(reader)
		if err != nil {
			return "", 0, fmt.Errorf("scan spaced number: %w", err)
		}

		return tok, v, nil
	}

	return "", 0, fmt.Errorf("%w: %q", ErrInvalidToken, tok)
}

func processEmployee(reader *bufio.Reader, writer *bufio.Writer, currRange *Range, optTemp *int) error {
	sign, value, err := scanConstraint(reader)
	if err != nil {
		return fmt.Errorf("read constraint: %w", err)
	}

	var newRange Range

	switch sign {
	case ">=":
		newRange = Range{value, 30}
	case "<=":
		newRange = Range{15, value}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidToken, sign)
	}

	if err := intersection(currRange, &newRange); err != nil {
		if _, werr := fmt.Fprintln(writer, -1); werr != nil {
			return fmt.Errorf("write -1: %w", werr)
		}

		currRange.low, currRange.high = 1, 0
		*optTemp = -1

		return nil
	}

	out := currRange.low
	if *optTemp >= currRange.low && *optTemp <= currRange.high {
		out = *optTemp
	} else {
		*optTemp = currRange.low
	}

	if _, err := fmt.Fprintln(writer, out); err != nil {
		return fmt.Errorf("write result: %w", err)
	}

	return nil
}

func processDepartment(reader *bufio.Reader, writer *bufio.Writer) error {
	emplNum, err := readInt(reader)
	if err != nil {
		return fmt.Errorf("read employees count: %w", err)
	}

	currRange := Range{15, 30}
	optTemp := -1

	for range emplNum {
		if err := processEmployee(reader, writer, &currRange, &optTemp); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	defer func() {
		if err := writer.Flush(); err != nil {
			_ = err
		}
	}()

	depNum, err := readInt(reader)
	if err != nil {
		return
	}

	for range depNum {
		if err := processDepartment(reader, writer); err != nil {
			return
		}
	}
}

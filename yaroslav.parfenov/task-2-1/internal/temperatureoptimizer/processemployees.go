package temperatureoptimizer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	errInvTemperature = errors.New("Invalid temperature")
	errInvSign        = errors.New("Invalid sign")
)

func ProcessEmployees(numEmployees *int) error {
	const (
		minTemperature = 15
		maxTemperature = 30
	)

	var (
		curBorder   int
		err         error
		leftBorder  = minTemperature
		rightBorder = maxTemperature
	)

	for *numEmployees > 0 {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		sign := scanner.Text()[0:2]

		switch sign {
		case "<=":
			curBorder, err = strconv.Atoi(scanner.Text()[3:])
			if curBorder <= rightBorder {
				rightBorder = curBorder
			}
		case ">=":
			curBorder, err = strconv.Atoi(scanner.Text()[3:])
			if curBorder >= leftBorder {
				leftBorder = curBorder
			}
		default:
			return errInvSign
		}

		if err != nil || curBorder < minTemperature || curBorder > maxTemperature {
			return errInvTemperature
		}

		if leftBorder <= rightBorder {
			fmt.Println(leftBorder)
		} else {
			fmt.Println(-1)
		}

		(*numEmployees)--
	}

	return nil
}

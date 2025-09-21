package temperatureoptimizer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func ProcessEmployees(numEmployees *int) error {

	const (
		minTemperature = 15
		maxTemperature = 30
	)

	var (
		invTemperature = errors.New("Invalid temperature")
		invSign        = errors.New("Invalid sign")
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
			return invSign
		}

		if err != nil || curBorder < minTemperature || curBorder > maxTemperature {
			return invTemperature
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

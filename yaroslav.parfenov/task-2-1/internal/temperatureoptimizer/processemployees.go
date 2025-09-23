package temperatureoptimizer

import (
	"errors"
	"fmt"
)

var (
	errInvRecord      = errors.New("invalid record format")
	errInvTemperature = errors.New("invalid temperature")
	errInvSign        = errors.New("invalid sign")
)

func ProcessEmployees(numEmployees *int) error {
	const (
		minTemperature = 15
		maxTemperature = 30
	)

	var (
		sign        string
		curBorder   int
		leftBorder  = minTemperature
		rightBorder = maxTemperature
	)

	for *numEmployees > 0 {
		_, err := fmt.Scanln(&sign, &curBorder)
		if err != nil {
			return errInvRecord
		}

		switch sign {
		case "<=":
			if curBorder <= rightBorder {
				rightBorder = curBorder
			}
		case ">=":
			if curBorder >= leftBorder {
				leftBorder = curBorder
			}
		default:
			return errInvSign
		}

		if curBorder < minTemperature || curBorder > maxTemperature {
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

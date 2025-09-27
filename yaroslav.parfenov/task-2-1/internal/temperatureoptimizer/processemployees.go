package temperatureoptimizer

import (
	"fmt"
)

const (
	errInvRecord      = "%w: invalid record format"
	errInvTemperature = "%w: invalid temperature"
	errInvSign        = "%w: invalid sign"
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

	for range *numEmployees {
		_, err := fmt.Scanln(&sign, &curBorder)
		if err != nil {
			return fmt.Errorf(errInvRecord, err)
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
			return fmt.Errorf(errInvSign, err)
		}

		if curBorder < minTemperature || curBorder > maxTemperature {
			return fmt.Errorf(errInvTemperature, err)
		}

		if leftBorder <= rightBorder {
			fmt.Println(leftBorder)
		} else {
			fmt.Println(-1)
		}
	}

	return nil
}

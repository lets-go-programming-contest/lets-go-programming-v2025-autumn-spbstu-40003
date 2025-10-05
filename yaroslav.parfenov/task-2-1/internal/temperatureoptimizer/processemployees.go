package temperatureoptimizer

import (
	"fmt"
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
			return fmt.Errorf("invalid record format: %w", err)
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
			return fmt.Errorf("invalid sign: %w", err)
		}

		if curBorder < minTemperature || curBorder > maxTemperature {
			return fmt.Errorf("invalid temperature: %w", err)
		}

		if leftBorder <= rightBorder {
			fmt.Println(leftBorder)
		} else {
			fmt.Println(-1)
		}
	}

	return nil
}

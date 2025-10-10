package temperature

import (
	"fmt"
)

const (
	tempSize   = 16
	tempOffset = 15
)

func applyLessEqual(difference int, temperature []int) ([]int, bool) {
	switch {
	case difference < 0:
		fmt.Println(-1)

		return temperature, false
	case difference >= len(temperature):
		fmt.Println(temperature[0])

		return temperature, true
	default:
		temperature = temperature[:difference+1]

		fmt.Println(temperature[0])

		return temperature, true
	}
}

func applyGreaterEqual(difference int, temperature []int) ([]int, bool) {
	switch {
	case difference >= len(temperature):
		fmt.Println(-1)

		return temperature, false
	case difference <= 0:
		fmt.Println(temperature[0])

		return temperature, true
	default:
		temperature = temperature[difference:]

		fmt.Println(temperature[0])

		return temperature, true
	}
}

func CheckRange(employeeCount int) error {
	var currentTemperature TemperatureRange

	temperature := make([]int, tempSize)

	for index := range temperature {
		temperature[index] = tempOffset + index
	}

	isValid := true

	for employeeCount > 0 {
		_, err := fmt.Scan(&currentTemperature.Range)
		if err != nil {
			return fmt.Errorf("error while scanning range: %w", err)
		}

		_, err = fmt.Scan(&currentTemperature.Value)
		if err != nil {
			return fmt.Errorf("error while scanning value: %w", err)
		}

		if !isValid {
			fmt.Println(-1)

			employeeCount--

			continue
		}

		difference := currentTemperature.Value - temperature[0]

		switch currentTemperature.Range {
		case "<=":
			temperature, isValid = applyLessEqual(difference, temperature)
		case ">=":
			temperature, isValid = applyGreaterEqual(difference, temperature)
		default:
			fmt.Println(-1)

			isValid = false
		}

		employeeCount--
	}

	return nil
}

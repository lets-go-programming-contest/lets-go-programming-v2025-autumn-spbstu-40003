package temperature

import (
	"fmt"
)

const (
	tempSize   = 16
	tempOffset = 15
)

func CheckRange(employeeCount int) {
	var (
		currentTemperature TemperatureRange
		difference         int
	)

	temperature := make([]int, tempSize)

	for index := range temperature {
		temperature[index] = tempOffset + index
	}

	isValid := true

	for employeeCount > 0 {
		_, err := fmt.Scan(&currentTemperature.Range)
		if err != nil {
			return
		}

		_, err = fmt.Scan(&currentTemperature.Value)
		if err != nil {
			return
		}

		if !isValid {
			fmt.Println(-1)
			employeeCount--
			continue
		}

		difference = currentTemperature.Value - temperature[0]

		switch currentTemperature.Range {
		case "<=":
			if difference < 0 {
				fmt.Println(-1)
				isValid = false
			} else if difference >= len(temperature) {
				fmt.Println(temperature[0])
			} else {
				temperature = temperature[:difference+1]
				fmt.Println(temperature[0])
			}

			employeeCount--

		case ">=":
			if difference >= len(temperature) {
				fmt.Println(-1)
				isValid = false
			} else if difference < 0 {
				fmt.Println(temperature[0])
			} else {
				temperature = temperature[difference:]
				fmt.Println(temperature[0])
			}
			employeeCount--

		default:
			return
		}
	}
}

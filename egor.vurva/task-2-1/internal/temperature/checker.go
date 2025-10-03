package temperature

import (
	"fmt"
)

func CheckRange(employeeCount int) {
	var currentTemperature TemperatureRange
	var difference int
	temperature := make([]int, 16)

	for index := 0; index < 16; index++ {
		temperature[index] = 15 + index
	}

	for employeeCount > 0 {
		_, err := fmt.Scan(&currentTemperature.Range)
		if err != nil {
			return
		}
		_, err = fmt.Scan(&currentTemperature.Value)
		if err != nil {
			return
		}

		difference = currentTemperature.Value - temperature[0]

		switch currentTemperature.Range {
		case "<=":
			if difference < 0 {
				fmt.Println(-1)
				return
			}
			temperature = temperature[:difference]
			fmt.Println(temperature[0])
			employeeCount--

		case ">=":
			if difference > len(temperature) {
				fmt.Println(-1)
				return
			}
			temperature = temperature[difference:]
			fmt.Println(temperature[0])
			employeeCount--

		default:
			return
		}
	}
}

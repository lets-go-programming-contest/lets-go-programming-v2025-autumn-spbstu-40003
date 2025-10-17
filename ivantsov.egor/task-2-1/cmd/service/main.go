package main

import (
	"fmt"
)

const (
	minTemp      = 15
	maxTemp      = 30
	invalidValue = -1
)

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println(invalidValue)
		return
	}

	for d := 0; d < departments; d++ {
		var employees int
		if _, err := fmt.Scan(&employees); err != nil {
			fmt.Println(invalidValue)
			return
		}

		minBound, maxBound := minTemp, maxTemp

		for e := 0; e < employees; e++ {
			var sign string
			var temp int
			if _, err := fmt.Scan(&sign, &temp); err != nil {
				fmt.Println(invalidValue)
				return
			}

			switch sign {
			case ">=":
				if temp > minBound {
					minBound = temp
				}
			case "<=":
				if temp < maxBound {
					maxBound = temp
				}
			}
		}

		if minBound <= maxBound {
			fmt.Println(minBound)
		} else {
			fmt.Println(invalidValue)
		}
	}
}

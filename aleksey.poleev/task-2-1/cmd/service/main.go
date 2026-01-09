package main

import "fmt"

const (
	minTemp = 15
	maxTemp = 30
)

func main() {
	var depCount int

	if _, scanError := fmt.Scan(&depCount); scanError != nil {
		fmt.Printf("Reading error: %v\n", scanError)
		fmt.Println(-1)

		return
	}

	for range depCount {
		var count int

		if _, scanError := fmt.Scan(&count); scanError != nil {
			fmt.Println(-1)

			return
		}

		tempCalc(count)
	}
}

func tempCalc(count int) {
	currMin := minTemp
	currMax := maxTemp
	isPoss := true

	for range count {
		var operator string

		var targetTemp int

		if _, err := fmt.Scan(&operator, &targetTemp); err != nil {
			fmt.Println(-1)

			return
		}

		if !isPoss {
			fmt.Println(-1)

			continue
		}

		switch operator {
		case ">=":
			if targetTemp > currMin {
				currMin = targetTemp
			}
		case "<=":
			if targetTemp < currMax {
				currMax = targetTemp
			}
		default:
			fmt.Println(-1)

			isPoss = false

			continue
		}

		if currMin <= currMax {
			fmt.Println(currMin)
		} else {
			fmt.Println(-1)

			isPoss = false
		}
	}
}

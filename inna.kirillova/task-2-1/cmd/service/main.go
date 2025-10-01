package main

import "fmt"

const (
	minTemp     = 15
	maxTemp     = 30
	InvalidTemp = -1
)

func main() {
	var numOfDep int

	_, err := fmt.Scan(&numOfDep)
	if err != nil {
		return
	}

	for range numOfDep {
		handleDepartment()
	}
}

func handleDepartment() {
	var numOfEmpl int

	_, err := fmt.Scan(&numOfEmpl)
	if err != nil {
		return
	}

	lowerLimit, upperLimit := minTemp, maxTemp
	flag := true

	for range numOfEmpl {
		var (
			sign string
			temp int
		)

		_, err := fmt.Scan(&sign, &temp)
		if err != nil {
			return
		}

		if !flag {
			fmt.Println(InvalidTemp)

			continue
		}

		if sign == ">=" {
			if temp > lowerLimit {
				lowerLimit = temp
			}
		} else {
			if temp < upperLimit {
				upperLimit = temp
			}
		}

		if lowerLimit <= upperLimit {
			fmt.Println(lowerLimit)
		} else {
			fmt.Println(InvalidTemp)

			flag = false
		}
	}
}

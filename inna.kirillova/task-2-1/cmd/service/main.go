package main

import "fmt"

const (
	minTemp     = 15
	maxTemp     = 30
	InvalidTemp = -1
)

func main() {
	var N int

	_, err := fmt.Scan(&N)
	if err != nil {
		return
	}

	for i := 0; i < N; i++ {
		handleDepartment()
	}
}

func handleDepartment() {
	var K int

	_, err := fmt.Scan(&K)
	if err != nil {
		return
	}

	lowerLimit, upperLimit := minTemp, maxTemp
	flag := true

	for j := 0; j < K; j++ {
		var sign string
		var temp int

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

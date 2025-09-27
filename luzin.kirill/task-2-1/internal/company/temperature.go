package company

import (
	"errors"
	"fmt"
)

var (
	errLogicInput       = errors.New("wrong input for >= and <=")
	errTemperatureInput = errors.New("wrong input for temperature")
)

func OptimizeTemperature(cEmployee int) error {
	var (
		minT, maxT, tempT, optimalT = 15, 30, 0, 0
		input                       = ""
	)

	for range cEmployee {
		_, err := fmt.Scan(&input)

		if err != nil {
			return errLogicInput
		}

		_, err = fmt.Scan(&tempT)

		if err != nil {
			return errTemperatureInput
		}

		switch input {
		case ">=":
			if tempT <= maxT && optimalT != -1 {
				if minT < tempT {
					minT = tempT
				}
			} else {
				optimalT = -1
			}

		case "<=":
			if tempT >= minT && optimalT != -1 {
				if maxT > tempT {
					maxT = tempT
				}
			} else {
				optimalT = -1
			}

		default:
			return errLogicInput
		}

		if optimalT != -1 {
			optimalT = minT
		}

		fmt.Println(optimalT)
	}

	return nil
}

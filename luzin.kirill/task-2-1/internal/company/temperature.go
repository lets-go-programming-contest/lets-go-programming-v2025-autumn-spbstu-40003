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

		minT, maxT, err = chooseLogic(minT, maxT, input, tempT)

		if err != nil {
			return err
		}

		if optimalT != -1 {
			optimalT = minT
		}

		fmt.Println(optimalT)
	}

	return nil
}

func chooseLogic(lowerT int, upperT int, str string, temp int) (int, int, error) {
	switch str {
	case ">=":
		if temp <= upperT {
			if lowerT < temp {
				lowerT = temp
			}
		} else {
			return -1, -1, nil
		}

	case "<=":
		if temp >= lowerT {
			if upperT > temp {
				upperT = temp
			}
		} else {
			return -1, -1, nil
		}

	default:
		return -1, -1, errLogicInput
	}

	return lowerT, upperT, nil
}

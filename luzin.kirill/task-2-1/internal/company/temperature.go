package company

import (
	"errors"
	"fmt"
)

var errWrongInput = errors.New("wrong input for >= and <=")

func OptimizeTemperature(cEmployee int) error {
	var (
		minT, maxT, tempT, optimalT int    = 15, 30, 0, 0
		input                       string = ""
	)

	for i := 0; i < cEmployee; i++ {
		_, err := fmt.Scan(&input)

		if err != nil {
			return err
		}

		_, err = fmt.Scan(&tempT)

		if err != nil {
			return err
		}

		switch input {
		case ">=":
			if optimalT != -1 {
				if tempT <= maxT {

					if minT < tempT {
						minT = tempT
					} else {
						tempT = minT
					}
					optimalT = tempT

				} else {
					optimalT = -1
				}
			}

		case "<=":
			if optimalT != -1 {
				if tempT >= minT {

					if maxT > tempT {
						maxT = tempT
					} else {
						tempT = maxT
					}
					optimalT = tempT

				} else {
					optimalT = -1
				}
			}

		default:
			return errWrongInput
		}

		fmt.Println(optimalT)
	}

	return nil
}

package solution

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func ProcessEmployees(numDepartments *int) error {

	const (
		minTemperature = 15
		maxTemperature = 30
	)

	var numCollegues int

	for *numDepartments > 0 {
		_, err := fmt.Scanln(&numCollegues)
		if err != nil {
			return errors.New("Invalid number of collegues")
		}

		var (
			leftBorder  = minTemperature
			rightBorder = maxTemperature
		)

		for numCollegues > 0 {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()

			sign := scanner.Text()[0:2]
			var (
				curBorder int
				err       error
			)

			if sign == "<=" {
				curBorder, err = strconv.Atoi(scanner.Text()[3:])

				if curBorder <= rightBorder {
					rightBorder = curBorder
				}
			} else if sign == ">=" {
				curBorder, err = strconv.Atoi(scanner.Text()[3:])

				if curBorder >= leftBorder {
					leftBorder = curBorder
				}
			} else {
				return errors.New("Incorrect sign")
			}

			if err != nil || curBorder < minTemperature || curBorder > maxTemperature {
				return errors.New("Invalid temperature")
			}

			if leftBorder <= rightBorder {
				fmt.Println(leftBorder)
			} else {
				fmt.Println(-1)
			}

			numCollegues--
		}

		(*numDepartments)--
	}
	return nil
}

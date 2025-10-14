package temperature

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	errDepartments = errors.New("failed to read the number of departments")
	errEmployees   = errors.New("failed to read the number of colleagues")
	errTemperature = errors.New("failed to read temperature constraint")
)

const (
	lowerTemperatureBound  = 15
	higherTemperatureBound = 30
	undefinedTemperature   = -1
)

type Temperature struct {
	lowTemperature  int
	highTemperature int
	isBroken        bool
}

func updateTemperature(input string, temperature *Temperature) error {
	var (
		operator string
		num      int
	)

	if _, err := fmt.Sscanf(input, "%s %d", &operator, &num); err != nil {
		return errTemperature
	}

	if !temperature.isBroken {
		switch operator {
		case "<=":
			if num < temperature.highTemperature {
				temperature.highTemperature = num
			}
		case ">=":
			if num > temperature.lowTemperature {
				temperature.lowTemperature = num
			}
		default:
			return errTemperature
		}

		rightTempCondition := temperature.lowTemperature <= temperature.highTemperature

		if !rightTempCondition {
			temperature.lowTemperature = undefinedTemperature
			temperature.isBroken = true
		}
	}

	return nil
}

func TemperatureControl() error {
	var (
		departments int
		employees   int
		temperature Temperature
		input       string
		err         error
	)

	reader := bufio.NewReader(os.Stdin)

	if _, err := fmt.Scanln(&departments); err != nil {
		return errDepartments
	}

	for i := 0; i < departments; i++ {
		if _, err := fmt.Scanln(&employees); err != nil {
			return errEmployees
		}

		temperature = Temperature{
			lowerTemperatureBound,
			higherTemperatureBound,
			false,
		}

		for j := 0; j < employees; j++ {
			input, err = reader.ReadString('\n')
			if err != nil {
				return errTemperature
			}

			if err := updateTemperature(strings.TrimSpace(input), &temperature); err != nil {
				return err
			}

			fmt.Println(temperature.lowTemperature)
		}
	}

	return nil
}

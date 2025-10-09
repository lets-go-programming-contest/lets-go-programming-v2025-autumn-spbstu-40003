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
	LOWER_TEMPERATURE_BOUND  = 15
	HIGHER_TEMPERATURE_BOUND = 30
	UNDEFINED_TEMPERATURE    = -1
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
			temperature.lowTemperature = UNDEFINED_TEMPERATURE
			temperature.isBroken = true
		}
	}

	return nil
}

func TemperatureControl() error {
	var (
		departments int
		employees   int
	)

	reader := bufio.NewReader(os.Stdin)

	if _, err := fmt.Scanln(&departments); err != nil {
		return errDepartments
	}

	for range departments {
		if _, err := fmt.Scanln(&employees); err != nil {
			return errEmployees
		}

		temperature := Temperature{
			LOWER_TEMPERATURE_BOUND,
			HIGHER_TEMPERATURE_BOUND,
			false,
		}

		for range employees {
			input, err := reader.ReadString('\n')
			if err != nil {
				return errTemperature
			}
			input = strings.TrimSpace(input)

			updateTemperature(input, &temperature)
			fmt.Println(temperature.lowTemperature)
		}
	}

	return nil
}

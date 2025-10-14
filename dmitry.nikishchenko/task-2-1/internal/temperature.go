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

func EmployeeTemperatureControl(
	employees int,
	reader *bufio.Reader,
	writer *bufio.Writer,
	temperature *Temperature) error {

	for range employees {
		input, err := reader.ReadString('\n')
		if err != nil {
			return errTemperature
		}

		if err := updateTemperature(strings.TrimSpace(input), temperature); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(writer, temperature.lowTemperature); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}

	return nil
}

func TemperatureControl() error {
	var (
		departments int
		employees   int
		temperature Temperature
	)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "flush error: %v\n", err)
		}
	}()

	if _, err := fmt.Scanln(&departments); err != nil {
		return errDepartments
	}

	for range departments {
		if _, err := fmt.Scanln(&employees); err != nil {
			return errEmployees
		}

		temperature = Temperature{
			lowerTemperatureBound,
			higherTemperatureBound,
			false,
		}

		err := EmployeeTemperatureControl(employees, reader, writer, &temperature)
		if err != nil {
			return err
		}
	}

	return nil
}

package conditioner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	errInvalidDepartmentFormat  = errors.New("invalid department count format")
	errInvalidDepartmentCount   = errors.New("department count must be greater than 0")
	errInvalidEmployeeFormat    = errors.New("invalid employee count format")
	errInvalidEmployeeCount     = errors.New("employee count must be greater than 0")
	errCommandRead              = errors.New("could not read command")
	errParseTemperature         = errors.New("could not parse temperature")
	errDataPrint                = errors.New("error printing data")
	errInvalidArgumentCount     = errors.New("invalid argument count")
	errInvalidTemperatureFormat = errors.New("invalid temperature value format")
	errTemperatureOutOfRange    = errors.New("temperature value out of range")
	errInvalidOperator          = errors.New("invalid operator")
)

const (
	expectedArgumentCount = 2
	minimalTemperature    = 15
	maximalTemperature    = 30
	operatorGreater       = ">="
	operatorLess          = "<="
)

type DepartmentProcessor struct {
	minimalSetTemperature int
	maximalSetTemperature int
}

func NewDepartmentProcessor() *DepartmentProcessor {
	return &DepartmentProcessor{
		minimalSetTemperature: minimalTemperature,
		maximalSetTemperature: maximalTemperature,
	}
}

func (dp *DepartmentProcessor) reset() {
	dp.minimalSetTemperature = minimalTemperature
	dp.maximalSetTemperature = maximalTemperature
}

func ProcessDepartments(reader io.Reader, writer io.Writer) error {
	processor := NewDepartmentProcessor()

	return processor.processDepartments(reader, writer)
}

func (dp *DepartmentProcessor) processDepartments(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		return errInvalidDepartmentFormat
	}

	departmentCountStr := scanner.Text()

	departmentCount, err := strconv.Atoi(departmentCountStr)
	if err != nil {
		return errInvalidDepartmentFormat
	}

	if departmentCount <= 0 {
		return errInvalidDepartmentCount
	}

	allResults := make([][]int, 0, departmentCount)

	for range departmentCount {
		if !scanner.Scan() {
			return errInvalidEmployeeFormat
		}

		employeeCountStr := scanner.Text()
		employeeCount, err := strconv.Atoi(employeeCountStr)

		if err != nil {
			return errInvalidEmployeeFormat
		}

		if employeeCount <= 0 {
			return errInvalidEmployeeCount
		}

		dp.reset()
		departmentResults := make([]int, 0, employeeCount)

		for range employeeCount {
			if !scanner.Scan() {
				return errCommandRead
			}

			command := scanner.Text()

			err = dp.parseTemperature(command)
			if err != nil {
				return errParseTemperature
			}

			if dp.minimalSetTemperature <= dp.maximalSetTemperature {
				departmentResults = append(departmentResults, dp.minimalSetTemperature)
			} else {
				departmentResults = append(departmentResults, -1)
			}
		}

		allResults = append(allResults, departmentResults)
	}

	for _, departmentResults := range allResults {
		for _, result := range departmentResults {
			_, err = fmt.Fprintln(writer, result)
			if err != nil {
				return errDataPrint
			}
		}
	}

	return nil
}

func (dp *DepartmentProcessor) parseTemperature(raw string) error {
	arguments := strings.Fields(raw)
	if len(arguments) != expectedArgumentCount {
		return errInvalidArgumentCount
	}

	operator := arguments[0]

	value, err := strconv.Atoi(arguments[1])
	if err != nil {
		return errInvalidTemperatureFormat
	}

	if value > maximalTemperature || value < minimalTemperature {
		return errTemperatureOutOfRange
	}

	switch operator {
	case operatorGreater:
		if value > dp.minimalSetTemperature {
			dp.minimalSetTemperature = value
		}
	case operatorLess:
		if value < dp.maximalSetTemperature {
			dp.maximalSetTemperature = value
		}
	default:
		return errInvalidOperator
	}

	return nil
}

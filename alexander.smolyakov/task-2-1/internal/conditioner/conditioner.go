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
    errInvalidEmployeeFormat = errors.New("Invalid employee count format.")
    errInvalidEmployeeCount = errors.New("Employee count must be greater than 0.")
    errCommandRead = errors.New("Could not read command.")
    errParseTemperature = errors.New("Could not parse temperature.")
    errDataPrint = errors.New("Error printing data.")
    errInvalidArgumentCount = errors.New("Invalid argument count.")
    errInvalidTemperatureFormat = errors.New("Invalid temperature value format.")
    errTemperatureOutOfRange = errors.New("Temperature value out of range.")
    errInvalidOperator = errors.New("Invalid operator.")
)

var (
    minimalSetTemperature int
    maximalSetTemperature int
)

const (
    expectedArgumentCount = 2
    minimalTemperature = 15
    maximalTemperature = 30
    operatorGreater = ">="
    operatorLess = "<="
)

func reset() {
    minimalSetTemperature = minimalTemperature
    maximalSetTemperature = maximalTemperature
}

func ProcessDepartment(reader io.Reader, writer io.Writer) error {
    bufReader := bufio.NewReader(reader)

    var employees int
    _, error := fmt.Fscanln(reader, &employees)
    if error != nil {
        return errInvalidEmployeeFormat
    }
    if employees <= 0 {
        return errInvalidEmployeeCount
    }

    reset()
    for i := 0; i < employees; i++ {
        command, error := bufReader.ReadString('\n')
        if error != nil {
            return errCommandRead
        }

        error = parseTemperature(command)
        if error != nil {
            return errParseTemperature 
        }

	if minimalSetTemperature <= maximalSetTemperature {
	    _, error := fmt.Fprintln(writer, minimalSetTemperature)
	    if error != nil {
	        return errDataPrint
	    }
	} else {
	    _, error := fmt.Fprintln(writer, -1)
	    if error != nil {
	        return errDataPrint
	    }
	}
    }
    return nil
}

func parseTemperature(raw string) error {
    arguments := strings.Fields(raw)
    if len(arguments) != expectedArgumentCount {
        return errInvalidArgumentCount
    }

    operator := arguments[0]
    value, error := strconv.Atoi(arguments[1])
    if error != nil {
        return errInvalidTemperatureFormat
    }

    if value > maximalTemperature || value < minimalTemperature {
        return errTemperatureOutOfRange
    }

    switch operator {
    case operatorGreater:
        minimalSetTemperature = value
    case operatorLess:
        maximalSetTemperature = value
    default:
        return errInvalidOperator
    }

    return nil
}

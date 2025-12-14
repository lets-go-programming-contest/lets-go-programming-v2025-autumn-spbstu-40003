package conditioner

import (
    "fmt"
    "errors"
    "strings"
    "strconv"
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

func Reset() {
    minimalSetTemperature = minimalTemperature
    maximalSetTemperature = maximalTemperature
}

func parseTemperature(raw string) error {
    arguments := strings.Fields(raw)
    if len(arguments) != expectedArgumentCount {
        errorMessage := fmt.Sprintf(
	    "Invalid argument count. Expected %d but %d given",
	    expectedArgumentCount, len(arguments))
        return errors.New(errorMessage)
    }

    operator := arguments[0]
    value, error := strconv.Atoi(arguments[1])
    if error != nil {
        return errors.New("Invalid temperature value format.")
    }

    if value > maximalTemperature || value < minimalTemperature {
        return errors.New("Temperature value out of range.")
    }

    switch operator {
    case operatorGreater:
        minimalSetTemperature = value
    case operatorLess:
	maximalSetTemperature = value
    default:
	return errors.New("Invalid operator.")
    }

    return nil
}

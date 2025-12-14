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

func ProcessDepartment(reader io.Reader, writer io.Writer) error {
    bufReader := bufio.NewReader(reader)

    var employees int
    _, error := fmt.Fscanln(reader, &employees)
    if error != nil {
        return errors.New("Invalid employee count format.")
    }
    if employees <= 0 {
        return errors.New("Employees count must be greater than 0.")
    }

    for i := 0; i < employees; i++ {
        command, error := bufReader.ReadString('\n')
        if error != nil {
            return errors.New("Could not read command.")
        }

        error = parseTemperature(command)
        if error != nil {
            return fmt.Errorf("Could not parse temperature: %w", error) 
        }
    }
    return nil
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

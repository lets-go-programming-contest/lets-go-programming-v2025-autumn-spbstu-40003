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
    errInvalidEmployeeFormat    = errors.New("invalid employee count format.")
    errInvalidEmployeeCount     = errors.New("employee count must be greater than 0.")
    errCommandRead              = errors.New("could not read command.")
    errParseTemperature         = errors.New("could not parse temperature.")
    errDataPrint                = errors.New("error printing data.")
    errInvalidArgumentCount     = errors.New("invalid argument count.")
    errInvalidTemperatureFormat = errors.New("invalid temperature value format.")
    errTemperatureOutOfRange    = errors.New("temperature value out of range.")
    errInvalidOperator          = errors.New("invalid operator.")
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

func ProcessDepartment(reader io.Reader, writer io.Writer) error {
    processor := NewDepartmentProcessor()
    return processor.processDepartment(reader, writer)
}

func (dp *DepartmentProcessor) processDepartment(reader io.Reader, writer io.Writer) error {
    bufReader := bufio.NewReader(reader)

    var employees int
    _, err := fmt.Fscanln(reader, &employees)
    if err != nil {
        return errInvalidEmployeeFormat
    }
    if employees <= 0 {
        return errInvalidEmployeeCount
    }

    dp.reset()
    
    for i := 0; i < employees; i++ {
        command, err := bufReader.ReadString('\n')
        if err != nil {
            return errCommandRead
        }

        err = dp.parseTemperature(command)
        if err != nil {
            return errParseTemperature
        }

        if dp.minimalSetTemperature <= dp.maximalSetTemperature {
            _, err := fmt.Fprintln(writer, dp.minimalSetTemperature)
            if err != nil {
                return errDataPrint
            }
        } else {
            _, err := fmt.Fprintln(writer, -1)
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
        dp.minimalSetTemperature = value
    case operatorLess:
        dp.maximalSetTemperature = value
    default:
        return errInvalidOperator
    }

    return nil
}

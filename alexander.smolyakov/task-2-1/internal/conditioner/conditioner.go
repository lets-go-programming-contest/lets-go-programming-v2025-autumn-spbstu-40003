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
	
	departmentCount, err := dp.readDepartmentCount(scanner)
	if err != nil {
		return err
	}
	
	allResults, err := dp.processAllDepartments(scanner, departmentCount)
	if err != nil {
		return err
	}
	
	return dp.printAllResults(writer, allResults)
}

func (dp *DepartmentProcessor) readDepartmentCount(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, errInvalidDepartmentFormat
	}
	
	departmentCountStr := scanner.Text()
	departmentCount, err := strconv.Atoi(departmentCountStr)
	if err != nil {
		return 0, errInvalidDepartmentFormat
	}
	
	if departmentCount <= 0 {
		return 0, errInvalidDepartmentCount
	}
	
	return departmentCount, nil
}

func (dp *DepartmentProcessor) processAllDepartments(scanner *bufio.Scanner, departmentCount int) ([][]int, error) {
	allResults := make([][]int, 0, departmentCount)
	
	for range departmentCount {
		departmentResults, err := dp.processSingleDepartment(scanner)
		if err != nil {
			return nil, err
		}
		
		allResults = append(allResults, departmentResults)
	}
	
	return allResults, nil
}

func (dp *DepartmentProcessor) processSingleDepartment(scanner *bufio.Scanner) ([]int, error) {
	employeeCount, err := dp.readEmployeeCount(scanner)
	if err != nil {
		return nil, err
	}
	
	dp.reset()
	departmentResults := make([]int, 0, employeeCount)
	
	for range employeeCount {
		result, err := dp.processEmployeeCommand(scanner)
		if err != nil {
			return nil, err
		}
		
		departmentResults = append(departmentResults, result)
	}
	
	return departmentResults, nil
}

func (dp *DepartmentProcessor) readEmployeeCount(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, errInvalidEmployeeFormat
	}
	
	employeeCountStr := scanner.Text()
	employeeCount, err := strconv.Atoi(employeeCountStr)
	if err != nil {
		return 0, errInvalidEmployeeFormat
	}
	
	if employeeCount <= 0 {
		return 0, errInvalidEmployeeCount
	}
	
	return employeeCount, nil
}

func (dp *DepartmentProcessor) processEmployeeCommand(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, errCommandRead
	}
	
	command := scanner.Text()
	err := dp.parseTemperature(command)
	if err != nil {
		return 0, errParseTemperature
	}
	
	if dp.minimalSetTemperature <= dp.maximalSetTemperature {
		return dp.minimalSetTemperature, nil
	}
	
	return -1, nil
}

func (dp *DepartmentProcessor) printAllResults(writer io.Writer, allResults [][]int) error {
	for _, departmentResults := range allResults {
		for _, result := range departmentResults {
			_, err := fmt.Fprintln(writer, result)
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

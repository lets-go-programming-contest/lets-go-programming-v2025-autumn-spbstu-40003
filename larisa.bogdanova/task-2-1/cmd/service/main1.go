package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MinTemperature = 15
	MaxTemperature = 30
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	reader := bufio.NewReader(os.Stdin)

	var departmentCount int
	_, err := fmt.Fscanln(reader, &departmentCount)
	if err != nil {
		return fmt.Errorf("failed to read department count: %w", err)
	}

	for deptIndex := 0; deptIndex < departmentCount; deptIndex++ {
		var employeeCount int
		_, err := fmt.Fscanln(reader, &employeeCount)
		if err != nil {
			return fmt.Errorf("failed to read employee count for department %d: %w", deptIndex+1, err)
		}

		lowerLimit := MinTemperature
		upperLimit := MaxTemperature
		hadFormatError := false

		for empIndex := 0; empIndex < employeeCount; {
			line, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input for employee %d in department %d: %w", empIndex+1, deptIndex+1, err)
			}
			line = strings.TrimSpace(line)

			if !strings.Contains(line, " ") {
				hadFormatError = true
				return fmt.Errorf("input format error: missing space between operator and number")
			}

			parts := strings.Fields(line)

			if hadFormatError {
				return fmt.Errorf("previous input format error — program terminated")
			}

			if len(parts) != 2 {
				hadFormatError = true
				return fmt.Errorf("invalid input format, expected operator and number separated by space")
			}

			operator := parts[0]
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				hadFormatError = true
				return fmt.Errorf("invalid number format: %w", err)
			}

			switch operator {
			case ">=":
				if value > upperLimit {
					fmt.Println(-1)
					return nil
				}
				if value > lowerLimit {
					lowerLimit = value
				}
			case "<=":
				if value < lowerLimit {
					fmt.Println(-1)
					return nil
				}
				if value < upperLimit {
					upperLimit = value
				}
			default:
				hadFormatError = true
				return fmt.Errorf("unknown operator %q", operator)
			}

			if lowerLimit > upperLimit {
				fmt.Println(-1)
				return nil
			}

			fmt.Println(lowerLimit)
			empIndex++
		}
	}

	return nil
}

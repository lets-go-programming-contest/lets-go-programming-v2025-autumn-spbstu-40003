package calc

import (
	"errors"
	"fmt"
)

func InputOperandsAndOperation(firOperand *int, secOperand *int, operation *byte) error {
	_, err := fmt.Scanln(firOperand)

	if err != nil {
		return errors.New("Invalid first operand")
	}

	_, err = fmt.Scanln(secOperand)

	if err != nil {
		return errors.New("Invalid second operand")
	}

	_, err = fmt.Scanf("%c", operation)

	if err != nil {
		return errors.New("Invalid operation")
	}

	if *operation == '/' && *secOperand == 0 {
		return errors.New("Division by zero")
	}

	return nil
}

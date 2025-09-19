package internal

import "fmt"

func Calculate() {
	var (
		firOperand int
		secOperand int
		operation  byte
		err        error
	)

	err = InputOperandsAndOperation(&firOperand, &secOperand, &operation)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		switch operation {
		case '+':
			fmt.Println(float64(firOperand + secOperand))
		case '-':
			fmt.Println(float64(firOperand - secOperand))
		case '*':
			fmt.Println(float64(firOperand * secOperand))
		case '/':
			fmt.Println(float64(firOperand / secOperand))
		}
	}
}

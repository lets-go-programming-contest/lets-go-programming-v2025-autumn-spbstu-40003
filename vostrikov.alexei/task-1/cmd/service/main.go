package main

import (
	"fmt"

	"github.com/lexachanskii/task-1/internal/calc"
)

func main() {
	res, err := calc.Calculate()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}

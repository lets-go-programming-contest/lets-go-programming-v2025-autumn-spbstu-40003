package main

import (
	"fmt"
	"lab_1/internal"
)

func main() {
	result, err := internal.Calculate()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

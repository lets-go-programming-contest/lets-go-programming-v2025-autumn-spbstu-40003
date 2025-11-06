package main

import (
	"fmt"

	"github.com/lexachanskii/task-2-2/internal/table"
)

func main() {
	dish, err := table.Table()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println(dish)
}

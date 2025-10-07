package main

import (
	"fmt"

	guests "github.com/lexachanskii/task-2-1/internal/office"
)

func main() {
	err := guests.ClimateControl()
	if err != nil {
		fmt.Println(err)
	}
}

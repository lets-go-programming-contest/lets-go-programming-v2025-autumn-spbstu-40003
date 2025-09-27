package main

import (
	"fmt"

	"github.com/ArtttNik/task-2-1/internal/utils"
)

func main() {
	err := utils.Temp()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}
}

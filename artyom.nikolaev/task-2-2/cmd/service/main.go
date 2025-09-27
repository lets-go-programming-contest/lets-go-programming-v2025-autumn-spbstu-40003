package main

import (
	"fmt"

	"github.com/ArtttNik/task-2-2/internal/utils"
)

func main() {
	err := utils.FindKDish()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}
}

package main

import (
	"fmt"

	"github.com/KiRy6A/task-2-2/internal/dishesMenu"
)

func main() {
	var (
		dishes       dishesMenu.Dishes
		selectedDish int
	)

	err := dishes.WriteMenu()
	if err != nil {
		fmt.Println("error writing menu of dishes:", err.Error())
		return
	}

	selectedDish, err = dishes.SelectDishe()
	if err != nil {
		fmt.Println("error choosing dish:", err.Error())
		return
	}

	fmt.Println(selectedDish)
}

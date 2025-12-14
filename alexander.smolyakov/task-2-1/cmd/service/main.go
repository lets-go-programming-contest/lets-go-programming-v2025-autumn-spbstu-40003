package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Ignitron1/task-2-1/internal/conditioner"
)

func main() {

	var writer io.Writer = os.Stdout
	var reader io.Reader = os.Stdin
	var departments int

	_, err := fmt.Fscan(reader, &departments)
	if err != nil {
		fmt.Println("Could not read number of departments.")
	}

	for i := 0; i < departments; i++ {
		err = conditioner.ProcessDepartment(reader, writer)
		if err != nil {
			fmt.Println(err)
		}
	}
}

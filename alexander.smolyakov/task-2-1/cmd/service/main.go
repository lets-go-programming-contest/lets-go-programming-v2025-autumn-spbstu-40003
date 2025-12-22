package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Ignitron1/task-2-1/internal/conditioner"
)

func main() {
	var (
		writer io.Writer = os.Stdout
		reader io.Reader = os.Stdin
	)

	err := conditioner.ProcessDepartments(reader, writer)
	if err != nil {
		fmt.Println(err)
	}
}

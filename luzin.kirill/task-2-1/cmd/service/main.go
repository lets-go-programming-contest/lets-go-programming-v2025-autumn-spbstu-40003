package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/KiRy6A/task-2-1/internal/company"
)

var (
	errDepartament = errors.New("going over limits of possible departament values")
	errEmployee    = errors.New("going over limits of possible employee values")
)

func main() {
	var cDepartament, cEmployee int

	_, err := fmt.Scan(&cDepartament)
	if err != nil {
		fmt.Println("error reading departament count:", err.Error())
		os.Exit(0)
	}

	if cDepartament < 1 || cDepartament > 1000 {
		fmt.Println(errDepartament.Error())
		os.Exit(0)
	}

	for range cDepartament {
		_, err = fmt.Scan(&cEmployee)
		if err != nil {
			fmt.Println("error reading employee count:", err.Error())
			os.Exit(0)
		}

		if cEmployee < 1 || cEmployee > 1000 {
			fmt.Println(errEmployee.Error())
			os.Exit(0)
		}

		err = company.OptimizeTemperature(cEmployee)
		if err != nil {
			fmt.Println("error temperature optimization:", err)
			os.Exit(0)
		}
	}
}

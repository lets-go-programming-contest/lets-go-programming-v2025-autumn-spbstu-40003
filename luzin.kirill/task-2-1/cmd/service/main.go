package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/KiRy6A/task-2-1/internal/company"
)

var (
	errDepartament = errors.New("wrong number of departaments")
	errEmployee    = errors.New("wrong number of employees")
)

func main() {

	var cDepartament, cEmployee = 0, 0

	_, err := fmt.Scan(&cDepartament)
	if err != nil {
		println(errDepartament.Error())
		os.Exit(0)
	}

	if cDepartament < 1 || cDepartament > 1000 {
		println(errDepartament.Error())
		os.Exit(0)
	}

	for range cDepartament {
		_, err = fmt.Scan(&cEmployee)
		if err != nil {
			println(errEmployee.Error())
			os.Exit(0)
		}

		if cEmployee < 1 || cEmployee > 1000 {
			println(errEmployee.Error())
			os.Exit(0)
		}

		err = company.OptimizeTemperature(cEmployee)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
}

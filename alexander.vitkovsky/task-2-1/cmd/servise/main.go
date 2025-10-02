package main

import (
	"fmt"
	"strconv"
)

type Range struct {
	low, high int
}

func intersection(r1 *Range, r2 *Range) error {
	if r1.high < r2.low || r1.low > r2.high { // не пересекаются
		return fmt.Errorf("the intervals do not overlap")
	} else { // как-то пересекаются
		r1.low = max(r1.low, r2.low)
		r1.high = min(r1.high, r2.high)
		return nil
	}
}


func main() {
	var (
		input, sign string
		value       int
	)

	_, err := fmt.Scan(&input)
	if err != nil {
		fmt.Println("Error: can`t read depNum")
		return
	}
	depNum, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Error: depNum is NAN")
		return
	}

	for depNum > 0 { // цикл по отделам
		_, err = fmt.Scan(&input)
		if err != nil {
			fmt.Println("Error: can`t read emplNum")
			return
		}

		emplNum, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Error: emplNum is NAN")
			return
		}

		currRange := Range{15, 30}
		optTemp := -1

		for emplNum > 0 { // цикл по сотрудникам внутри отдела
			_, err = fmt.Scan(&sign, &value)
			if err != nil {
				fmt.Println("Error: invalid temperature range")
				return
			}

			var newRange Range
			switch sign {
			case ">=":
				newRange = Range{value, 30}
			case "<=":
				newRange = Range{15, value}
			default:
				fmt.Printf("Error: invalid sign: %s\n", sign)
			}
			err = intersection(&currRange, &newRange)
			if err != nil {
				fmt.Println(-1) // диапазоны не пересеклись
			} else {
				/*
					В задании нет четкого определения "оптимальной температуры".
					Считаю, что	«оптимальная температура» — это минимально возможная
					температура	в текущем допустимом диапазоне, но если старое значение
					ещё подходит, то оно сохраняется
				*/
				if optTemp >= currRange.low && optTemp <= currRange.high {
					fmt.Println(optTemp) // старое значение подходит
				} else {
					optTemp = currRange.low // берём минимально возможное
					fmt.Println(optTemp)
				}
			}

			emplNum--
		}
		depNum--
	}
}
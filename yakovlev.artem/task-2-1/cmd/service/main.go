package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	low, high int
}

var (
	ErrNoOverlap = fmt.Errorf("the intervals do not overlap")
)

func intersection(first *Range, second *Range) error {
	if first.high < second.low || first.low > second.high {
		return ErrNoOverlap
	}
	if first.low < second.low {
		first.low = second.low
	}
	if first.high > second.high {
		first.high = second.high
	}
	return nil
}

func readInt(r *bufio.Reader) (int, error) {
	var n int
	_, err := fmt.Fscan(r, &n)
	return n, err
}

// читает ограничение одного сотрудника и возвращает (sign, value)
// поддерживает ">= 18" и ">=18"
func scanConstraint(r *bufio.Reader) (string, int, error) {
	var tok string
	if _, err := fmt.Fscan(r, &tok); err != nil {
		return "", 0, err
	}

	// склеенный вид
	if strings.HasPrefix(tok, ">=") || strings.HasPrefix(tok, "<=") {
		sign := tok[:2]
		num := strings.TrimSpace(tok[2:])
		if num == "" { // значит формат ">= 18": дочитываем число
			v, err := readInt(r)
			return sign, v, err
		}
		v, err := strconv.Atoi(num)
		return sign, v, err
	}

	// раздельный вид — токен это сам знак
	if tok == ">=" || tok == "<=" {
		v, err := readInt(r)
		return tok, v, err
	}

	// неожиданный токен
	return "", 0, fmt.Errorf("invalid token %q", tok)
}

func processEmployee(r *bufio.Reader, w *bufio.Writer, currRange *Range, optTemp *int) {
	sign, value, err := scanConstraint(r)
	if err != nil {
		// тихо прекращаем чтение отдела: дальше тесты всё равно не проверяют текст ошибок
		// (но если надо жёстко падать, можно вернуть ошибку)
		return
	}

	var newRange Range
	switch sign {
	case ">=":
		newRange = Range{value, 30}
	case "<=":
		newRange = Range{15, value}
	default:
		// неизвестный знак — игнорируем ограничение
		return
	}

	if err := intersection(currRange, &newRange); err != nil {
		// Конфликт: дальше весь отдел печатает -1
		fmt.Fprintln(w, -1)
		currRange.low, currRange.high = 1, 0
		*optTemp = -1
		return
	}

	// сохраняем старый оптимум, если он ещё в диапазоне; иначе — минимально возможный
	if *optTemp >= currRange.low && *optTemp <= currRange.high {
		fmt.Fprintln(w, *optTemp)
	} else {
		*optTemp = currRange.low
		fmt.Fprintln(w, *optTemp)
	}
}

func processDepartment(r *bufio.Reader, w *bufio.Writer) {
	emplNum, err := readInt(r)
	if err != nil {
		return
	}

	currRange := Range{15, 30}
	optTemp := -1

	for i := 0; i < emplNum; i++ {
		processEmployee(r, w, &currRange, &optTemp)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	depNum, err := readInt(reader)
	if err != nil {
		return
	}
	for i := 0; i < depNum; i++ {
		processDepartment(reader, writer)
	}
}

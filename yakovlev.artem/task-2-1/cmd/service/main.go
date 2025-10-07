package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	minAge = 0
	maxAge = 120
)

type AgeMultiset struct {
	counts [maxAge + 1]int
	size   int
}

func (s *AgeMultiset) add(age int) {
	if age < minAge || age > maxAge {
		return
	}

	s.counts[age]++
	s.size++
}

func (s *AgeMultiset) del(age int) {
	if age < minAge || age > maxAge {
		return
	}

	if s.counts[age] > 0 {
		s.counts[age]--
		s.size--
	}
}

func (s *AgeMultiset) kth(k int) int {
	if k <= 0 || k > s.size {
		return -1
	}

	remain := k

	for age := minAge; age <= maxAge; age++ {
		if s.counts[age] >= remain {
			return age
		}
		remain -= s.counts[age]
	}

	return -1
}

func readInitialEmployees(reader *bufio.Reader, set *AgeMultiset) {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	for range n {
		var age int
		if _, err := fmt.Fscan(reader, &age); err != nil {
			break
		}

		set.add(age)
	}
}

func processStream(reader *bufio.Reader, writer *bufio.Writer, set *AgeMultiset) {
	for {
		var operation string
		var value int

		_, err := fmt.Fscan(reader, &operation, &value)
		if err != nil {
			break
		}

		switch operation {
		case "add", "hire", "+":
			set.add(value)
		case "del", "fire", "-":
			set.del(value)
		case "get", "query", "?", "kth":
			result := set.kth(value)
			_, _ = fmt.Fprintln(writer, result)
		default:
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer func() { _ = writer.Flush() }()

	var ages AgeMultiset

	readInitialEmployees(reader, &ages)
	processStream(reader, writer, &ages)
}

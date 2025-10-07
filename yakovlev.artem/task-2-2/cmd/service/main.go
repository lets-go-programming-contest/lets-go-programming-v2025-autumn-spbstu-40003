package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type IntMinHeap []int

func (h *IntMinHeap) Len() int           { return len(*h) }
func (h *IntMinHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntMinHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntMinHeap) Push(x any) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *IntMinHeap) Pop() any {
	old := *h

	size := len(old)
	if size == 0 {
		return 0
	}

	val := old[size-1]

	*h = old[:size-1]

	return val
}

func readInt(r *bufio.Reader) (int, error) {
	var tok string
	if _, err := fmt.Fscan(r, &tok); err != nil {
		return 0, fmt.Errorf("scan int: %w", err)
	}

	num, err := strconv.Atoi(tok)
	if err != nil {
		return 0, fmt.Errorf("atoi: %w", err)
	}

	return num, nil
}

func kthPreferred(scores []int, kth int) int {
	if kth <= 0 || len(scores) == 0 {
		return 0
	}

	minHeap := &IntMinHeap{}
	heap.Init(minHeap)

	for _, score := range scores {
		heap.Push(minHeap, score)

		if minHeap.Len() > kth {
			heap.Pop(minHeap)
		}
	}

	return (*minHeap)[0]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	defer func() {
		if err := writer.Flush(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "flush:", err)
		}
	}()

	count, err := readInt(reader)
	if err != nil && !errors.Is(err, io.EOF) {
		_, _ = fmt.Fprintln(os.Stderr, "read N:", err)

		return
	}

	values := make([]int, count)
	for idx := range values {
		val, rerr := readInt(reader)
		if rerr != nil {
			if !errors.Is(rerr, io.EOF) {
				_, _ = fmt.Fprintln(os.Stderr, "read value:", rerr)

				return
			}

			val = 0
		}

		values[idx] = val
	}

	kth, err := readInt(reader)
	if err != nil && !errors.Is(err, io.EOF) {
		_, _ = fmt.Fprintln(os.Stderr, "read k:", err)

		return
	}

	ans := kthPreferred(values, kth)

	if _, werr := fmt.Fprintln(writer, ans); werr != nil {
		_, _ = fmt.Fprintln(os.Stderr, "write:", werr)
	}
}

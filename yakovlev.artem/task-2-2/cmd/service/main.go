package main

import (
	"bufio"
	"container/heap"
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
	v, ok := x.(int)
	if !ok {
		return
	}
	*h = append(*h, v)
}

func (h *IntMinHeap) Pop() any {
	old := *h
	size := len(old)
	if size == 0 {
		return 0
	}
	v := old[size-1]
	*h = old[:size-1]
	return v
}

func readInt(r *bufio.Reader) (int, error) {
	var tok string
	if _, err := fmt.Fscan(r, &tok); err != nil {
		return 0, fmt.Errorf("scan int: %w", err)
	}

	v, err := strconv.Atoi(tok)
	if err != nil {
		return 0, fmt.Errorf("atoi: %w", err)
	}

	return v, nil
}

func kthPreferred(scores []int, kth int) int {
	if kth <= 0 || len(scores) == 0 {
		return 0
	}

	h := &IntMinHeap{}
	heap.Init(h)

	for _, score := range scores {
		heap.Push(h, score)
		if h.Len() > kth {
			heap.Pop(h)
		}
	}

	return (*h)[0]
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
	if err != nil && err != io.EOF {
		_, _ = fmt.Fprintln(os.Stderr, "read N:", err)
		return
	}

	values := make([]int, 0, count)
	for idx := 0; idx < count; idx++ {
		val, err := readInt(reader)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "read value:", err)
			return
		}
		values = append(values, val)
	}

	kth, err := readInt(reader)
	if err != nil && err != io.EOF {
		_, _ = fmt.Fprintln(os.Stderr, "read k:", err)
		return
	}

	answer := kthPreferred(values, kth)
	if _, err := fmt.Fprintln(writer, answer); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "write:", err)
	}
}

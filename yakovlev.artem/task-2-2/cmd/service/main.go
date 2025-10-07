package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"os"
)

type IntMinHeap []int

func (h IntMinHeap) Len() int           { return len(h) }
func (h IntMinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntMinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntMinHeap) Pop() any {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]

	return val
}

func readInt(reader *bufio.Reader) (int, error) {
	var value int

	_, err := fmt.Fscan(reader, &value)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, fmt.Errorf("unexpected EOF: %w", err)
		}

		return 0, fmt.Errorf("scan int: %w", err)
	}

	return value, nil
}

func writeInt(writer *bufio.Writer, value int) error {
	_, err := fmt.Fprintln(writer, value)
	if err != nil {
		return fmt.Errorf("write int: %w", err)
	}

	return nil
}

func kthPreferred(scores []int, k int) int {
	minH := &IntMinHeap{}

	heap.Init(minH)

	for _, v := range scores {
		if minH.Len() < k {
			heap.Push(minH, v)

			continue
		}

		if v > (*minH)[0] {
			heap.Pop(minH)
			heap.Push(minH, v)
		}
	}

	return (*minH)[0]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	defer func() { _ = writer.Flush() }()

	n, err := readInt(reader)
	if err != nil {
		return
	}

	scores := make([]int, 0, n)

	for range n {
		val, scanErr := readInt(reader)
		if scanErr != nil {
			return
		}

		scores = append(scores, val)
	}

	k, err := readInt(reader)
	if err != nil {
		return
	}

	if k < 1 || k > n {
		return
	}

	answer := kthPreferred(scores, k)

	if err := writeInt(writer, answer); err != nil {
		return
	}
}

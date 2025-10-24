package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	length := len(old)
	x := old[length-1]
	*h = old[0 : length-1]

	return x
}

func main() {
	var elementCount, kIndex int

	if _, err := fmt.Scan(&elementCount); err != nil {
		fmt.Printf("Error reading input (context: elementCount): %v\n", err)

		return
	}

	arr := make([]int, elementCount)
	for i := range arr {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			fmt.Printf("Error reading input (context: array element %d): %v\n", i, err)

			return
		}
	}

	if _, err := fmt.Scan(&kIndex); err != nil {
		fmt.Printf("Error reading input (context: kIndex): %v\n", err)

		return
	}

	heapData := &IntHeap{}
	heap.Init(heapData)

	for _, v := range arr {
		heap.Push(heapData, v)

		if heapData.Len() > kIndex {
			heap.Pop(heapData)
		}
	}

	if heapData.Len() > 0 {
		fmt.Println((*heapData)[0])
	}
}

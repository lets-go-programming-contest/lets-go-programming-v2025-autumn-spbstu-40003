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
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var n, kIndex int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	arr := make([]int, n)
	for i := range arr {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			return
		}
	}

	if _, err := fmt.Scan(&kIndex); err != nil {
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

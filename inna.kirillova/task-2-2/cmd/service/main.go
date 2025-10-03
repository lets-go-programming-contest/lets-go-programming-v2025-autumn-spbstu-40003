package main

import (
	"container/heap"
	"fmt"
)

const errorValue = 0

type IntHeap []int

func (heap IntHeap) Len() int {
	return len(heap)
}

func (heap IntHeap) Less(i, j int) bool {
	return heap[i] < heap[j]
}

func (heap IntHeap) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}

func (heap *IntHeap) Push(value interface{}) {
	value, flag := value.(int)
	if !flag {
		return
	}
	*heap = append(*heap, value.(int))
}

func (heap *IntHeap) Pop() interface{} {
	currentHeap := *heap
	length := len(currentHeap)
	lastElement := currentHeap[length-1]
	*heap = currentHeap[0 : length-1]

	return lastElement
}

func main() {
	var (
		countOfDish      int
		preferenceOfDish int
	)

	_, err := fmt.Scan(&countOfDish)
	if err != nil || countOfDish < 1 {
		fmt.Println(errorValue)
		return
	}

	dishRate := make([]int, countOfDish)
	for i := range countOfDish {
		_, err := fmt.Scan(&dishRate[i])
		if err != nil {
			fmt.Println(errorValue)
			return
		}
	}

	_, err = fmt.Scan(&preferenceOfDish)
	if err != nil || preferenceOfDish < 1 || preferenceOfDish > countOfDish {
		fmt.Println(errorValue)
		return
	}

	result := getPreference(dishRate, preferenceOfDish)
	fmt.Println(result)
}

func getPreference(dishes []int, preference int) int {
	heapData := &IntHeap{}
	heap.Init(heapData)

	for _, rating := range dishes {
		heap.Push(heapData, rating)

		if heapData.Len() > preference {
			heap.Pop(heapData)
		}
	}

	return (*heapData)[0]
}

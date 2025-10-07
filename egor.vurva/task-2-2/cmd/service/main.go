package main

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/Vurvaa/task-2-2/internal/maxheap"
)

func main() {
	maxHeap := &maxheap.MaxIntHeap{}
	heap.Init(maxHeap)

	var countFood, idxFood int

	_, err := fmt.Scan(&countFood)
	if err != nil {
		fmt.Println("Failed to read countFood: ", err)

		return
	}

	for countFood > 0 {
		var someFood string
		_, err = fmt.Scan(&someFood)
		if err != nil {
			fmt.Println("Failed to read food sequence: ", err)

			return
		}

		currFood, err := strconv.Atoi(someFood)
		if err != nil {
			fmt.Println("Failed to parse: ", err)

			return
		}

		heap.Push(maxHeap, currFood)

		countFood--
	}

	_, err = fmt.Scan(&idxFood)
	if err != nil {
		fmt.Println("Failed to read idxFood: ", err)

		return
	}

	for idxFood > 1 {
		heap.Pop(maxHeap)

		idxFood--
	}

	fmt.Println(heap.Pop(maxHeap))
}

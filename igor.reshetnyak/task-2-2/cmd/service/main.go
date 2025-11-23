package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrDishes       = errors.New("wrong num of dishes")
	ErrRaitingValue = errors.New("wrong dish raiting value")
	ErrRaitingNum   = errors.New("wrong raiting num")
	ErrAnyToInt     = errors.New("failed to convert interface{} to int")
)

const (
	minDishValue = 1
	maxValue     = 10000
	minValue     = -10000
)

type Heap []int

func (h *Heap) Len() int {
	return len(*h)
}

func (h *Heap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *Heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *Heap) Push(num interface{}) {
	intNum, ok := num.(int)
	if !ok {
		fmt.Println(ErrAnyToInt)

		return
	}

	*h = append(*h, intNum)
}

func (h *Heap) Pop() interface{} {
	oldHeap := *h
	oldLen := len(oldHeap)
	num := oldHeap[oldLen-1]
	*h = oldHeap[0 : oldLen-1]

	return num
}

func main() {
	var dishes, raitingNum int

	if _, err := fmt.Scan(&dishes); err != nil ||
		dishes < minDishValue || dishes > maxValue {
		fmt.Println(ErrDishes, err)

		return
	}

	dishRatings := make([]int, dishes)

	for index := range dishes {
		if _, err := fmt.Scan(&dishRatings[index]); err != nil ||
			dishRatings[index] < minValue || dishRatings[index] > maxValue {
			fmt.Println(ErrRaitingValue, err)

			return
		}
	}

	if _, err := fmt.Scan(&raitingNum); err != nil {
		fmt.Println(ErrRaitingNum, err)

		return
	}

	optimalRaiting := findOptimalRaiting(dishRatings, raitingNum)

	fmt.Println(optimalRaiting)
}

func findOptimalRaiting(dishRatings []int, raitingNum int) int {
	dishHeap := &Heap{}
	heap.Init(dishHeap)

	for _, dish := range dishRatings {
		if dishHeap.Len() < raitingNum {
			heap.Push(dishHeap, dish)
		} else if dish > (*dishHeap)[0] {
			heap.Pop(dishHeap)
			heap.Push(dishHeap, dish)
		}
	}

	return (*dishHeap)[0]
}

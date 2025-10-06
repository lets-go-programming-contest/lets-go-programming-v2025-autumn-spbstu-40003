package main

import (
	"container/heap"
	"errors"
	"fmt"
	"strconv"
)

type MaxHeap []int

func (h *MaxHeap) Len() int { return len(*h) }

func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }

func (h *MaxHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

var (
	ErrReadDishNum    = errors.New("can't read dishNum")
	ErrReadDishes     = errors.New("can't read dishes")
	ErrReadPreference = errors.New("can't read preference")
)

func readInt() (int, error) {
	var input string

	_, err := fmt.Scan(&input)
	if err != nil {
		return 0, fmt.Errorf("read value: %w", err)
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("atoi: %w", err)
	}

	return val, nil
}

func findPreference(nums []int, preference int) int {
	dishes := &MaxHeap{}

	for _, dish := range nums {
		heap.Push(dishes, dish)
	}

	heap.Init(dishes)

	var chosen int
	for range make([]struct{}, preference) {
		v, ok := heap.Pop(dishes).(int)
		if ok {
			chosen = v
		}
	}

	return chosen
}

func main() {
	dishNum, err := readInt()
	if err != nil {
		fmt.Println("Error:", ErrReadDishNum, err)

		return
	}

	nums := make([]int, dishNum)

	for index := 0; index < dishNum; index++ {
		nums[index], err = readInt()
		if err != nil {
			fmt.Println("Error:", ErrReadDishes, err)

			return
		}
	}

	preference, err := readInt()
	if err != nil {
		fmt.Println("Error:", ErrReadPreference, err)

		return
	}

	fmt.Println(findPreference(nums, preference))
}

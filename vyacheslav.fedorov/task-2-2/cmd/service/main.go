package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrItRat  = errors.New("invalid rating value")
	ErrSel    = errors.New("invalid k value")
	ErrICount = errors.New("invalid number of menu items")
)

const (
	lRat  = -10000
	hRat  = 10000
	minIt = 1
)

type Score []int

func (d *Score) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*d = append(*d, val)
	}
}

func (d *Score) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

func (d *Score) Less(i, j int) bool {
	return (*d)[i] > (*d)[j]
}

func (d *Score) Len() int {
	return len(*d)
}

func (d *Score) Pop() interface{} {
	current := *d
	n := len(current)
	elem := current[n-1]
	*d = current[:n-1]

	return elem
}

func main() {
	var (
		iCount    int
		selection int
	)

	_, scanErr1 := fmt.Scan(&iCount)
	if scanErr1 != nil || iCount < minIt || iCount > hRat {
		fmt.Println(ErrICount)

		return
	}

	mData, err := collectScores(iCount, lRat, hRat)
	if err != nil {
		fmt.Println(err)

		return
	}

	_, scanErr2 := fmt.Scan(&selection)
	if scanErr2 != nil || selection < 1 || selection > iCount {
		fmt.Println(ErrSel)

		return
	}

	remaining := selection
	for remaining > 0 {
		result, Tr := heap.Pop(mData).(int)
		if !Tr {
			fmt.Println(-1)

			return
		}

		remaining--
		if remaining == 0 {
			fmt.Println(result)
		}
	}
}

func collectScores(count, minB, maxB int) (*Score, error) {
	data := &Score{}
	heap.Init(data)

	for range count {
		var score int
		_, scanErr := fmt.Scan(&score)

		if scanErr != nil || score < minB || score > maxB {
			return nil, ErrItRat
		}

		heap.Push(data, score)
	}

	return data, nil
}

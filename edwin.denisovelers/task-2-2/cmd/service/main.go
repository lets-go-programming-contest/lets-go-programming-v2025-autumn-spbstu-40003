package main

import (
	"container/heap"
	"fmt"
	"os"

	"github.com/wedwincode/task-2-2/internal/intheap"
)

func main() {
	var size int

	_, err := fmt.Fscan(os.Stdin, &size)
	if err != nil {
		return
	}

	arr := make([]int, size)
	for i := range size {
		_, err := fmt.Fscan(os.Stdin, &arr[i])
		if err != nil {
			return
		}
	}

	var kValue int

	_, err = fmt.Fscan(os.Stdin, &kValue)
	if err != nil {
		return
	}

	minHeap := &intheap.IntHeap{}

	for _, v := range arr {
		heap.Push(minHeap, v)

		if minHeap.Len() > kValue {
			heap.Pop(minHeap)
		}
	}

	fmt.Println((*minHeap)[0])
}

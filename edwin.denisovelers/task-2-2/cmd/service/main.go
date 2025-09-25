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
	for i := 0; i < size; i++ {
		_, err := fmt.Fscan(os.Stdin, &arr[i])
		if err != nil {
			return
		}
	}

	var k int

	_, err = fmt.Fscan(os.Stdin, &k)
	if err != nil {
		return
	}

	h := &intheap.IntHeap{}

	for _, v := range arr {
		heap.Push(h, v)
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	fmt.Println((*h)[0])
}

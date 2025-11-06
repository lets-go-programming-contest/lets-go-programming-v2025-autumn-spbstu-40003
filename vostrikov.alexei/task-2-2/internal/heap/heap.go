package heap

type HeapInt []int

func (heap *HeapInt) Len() int { return len(*heap) }

func (heap *HeapInt) Less(i int, j int) bool { return (*heap)[i] > (*heap)[j] }

func (heap *HeapInt) Swap(i int, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *HeapInt) Push(x interface{}) {
	if x, ok := x.(int); ok {
		*heap = append(*heap, x)
	}
}

func (heap *HeapInt) Pop() interface{} {
	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]

	return x
}

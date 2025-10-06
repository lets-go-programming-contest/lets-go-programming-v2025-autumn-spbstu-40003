package main

import "testing"

func TestFindPreference(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{3, 2, 1, 5, 6, 4}, 2, 5},     // классический пример
		{[]int{7, 10, 4, 3, 20, 15}, 3, 10}, // проверка для k=3
		{[]int{1, 2, 3, 4, 5}, 1, 5},        // максимум
		{[]int{1, 2, 3, 4, 5}, 5, 1},        // минимум
		{[]int{5, 5, 5, 5}, 2, 5},           // одинаковые элементы
	}

	for _, tt := range tests {
		got := findPreference(tt.nums, tt.k)
		if got != tt.want {
			t.Errorf("findKthLargest(%v, %d) = %d; want %d", tt.nums, tt.k, got, tt.want)
		}
	}
}

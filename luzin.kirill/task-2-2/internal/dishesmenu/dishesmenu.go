package dishesmenu

import (
	"container/heap"
	"errors"
	"fmt"

	ownHeap "github.com/KiRy6A/task-2-2/internal/heap"
)

var (
	errLimit = errors.New("going over acceptable values")
	errType  = errors.New("expected int type")
)

type Dishes struct {
	menu ownHeap.IntHeap
}

func (dishes *Dishes) WriteMenu() error {
	var cDishes, rating int

	_, err := fmt.Scan(&cDishes)
	if err != nil {
		return fmt.Errorf("error scanning counter of dishes: %w", err)
	}

	for range cDishes {
		_, err := fmt.Scan(&rating)
		if err != nil {
			return fmt.Errorf("error scanning counter of dishes: %w", err)
		}

		heap.Push(&dishes.menu, rating)
	}

	return nil
}

func (dishes *Dishes) SelectDishe() (int, error) {
	var selectedNumber, foundedDish int

	_, err := fmt.Scan(&selectedNumber)
	if err != nil {
		return 0, fmt.Errorf("error scanning number of selected dish: %w", err)
	}

	if selectedNumber < 1 || selectedNumber > dishes.menu.Len() {
		return 0, fmt.Errorf("error limit selected dish: %w", errLimit)
	}

	for range selectedNumber - 1 {
		heap.Pop(&dishes.menu)
	}

	foundedDish, ok := heap.Pop(&dishes.menu).(int)
	if !ok {
		return 0, fmt.Errorf("error type when popping: %w", errType)
	}

	return foundedDish, nil
}

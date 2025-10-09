package dishesMenu

import (
	"container/heap"
	"errors"
	"fmt"

	ownHeap "github.com/KiRy6A/task-2-2/internal/heap"
)

var errLimit = errors.New("going over acceptable values")

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
	var k, foundedDish int

	_, err := fmt.Scan(&k)
	if err != nil {
		return 0, fmt.Errorf("error scanning number of selected dish: %w", err)
	}

	if k < 1 || k > dishes.menu.Len() {
		return 0, fmt.Errorf("error limit selected dish: %w", errLimit)
	}

	for range k {
		foundedDish = heap.Pop(&dishes.menu).(int)
	}

	return foundedDish, nil
}

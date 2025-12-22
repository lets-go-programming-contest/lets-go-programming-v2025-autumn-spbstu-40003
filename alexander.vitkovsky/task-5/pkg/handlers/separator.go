package handlers

import (
	"context"
	"errors"
)

var ErrNoOutputs = errors.New("no outputs found")

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return ErrNoOutputs
	}

	defer func() {
		for _, output := range outputs {
			close(output)
		}
	}()

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case message, exists := <-input:
			if !exists {
				return nil
			}
			select {
			case <-ctx.Done():
				return nil
			case outputs[index%len(outputs)] <- message:
				index++
			}
		}
	}
}

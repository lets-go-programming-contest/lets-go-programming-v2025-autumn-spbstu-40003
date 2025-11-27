package handlers

import (
	"context"
	"errors"
)

var errNoOutputsFound = errors.New("no outputs found")

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	currentIndex := 0
	outputCount := len(outputs)

	if outputCount == 0 {
		return errNoOutputsFound
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[currentIndex%outputCount] <- data:
				currentIndex++
			}
		}
	}
}

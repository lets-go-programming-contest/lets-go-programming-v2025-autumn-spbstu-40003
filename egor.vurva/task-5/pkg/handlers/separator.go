package handlers

import (
	"context"
	"errors"
)

var errNoOutputChannels = errors.New("there are no output channels")

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	output []chan string,
) error {
	if len(output) == 0 {
		return errNoOutputChannels
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			outChan := output[index%len(output)]

			select {
			case <-ctx.Done():
				return nil

			case outChan <- data:
				index++
			}
		}
	}
}

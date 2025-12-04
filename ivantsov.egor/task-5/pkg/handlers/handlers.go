package handlers

import (
	"context"
	"errors"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

type FanInHandler func(ctx context.Context, inputs []chan string, output chan string) error

func FanIn(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, inputChannel := range inputs {
		currentChannel := inputChannel

		go func(inChan chan string) {
			defer waitGroup.Done()

			for {
				select {
				case value, ok := <-inChan:
					if !ok {
						return
					}

					output <- value
				case <-ctx.Done():
					return
				}
			}
		}(currentChannel)
	}

	waitGroup.Wait()

	close(output)

	return nil
}

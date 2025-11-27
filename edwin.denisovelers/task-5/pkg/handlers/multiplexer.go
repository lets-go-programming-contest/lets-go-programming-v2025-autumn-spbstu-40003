package handlers

import (
	"context"
	"strings"
	"sync"
)

const noMultiplexer = "no multiplexer"

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	defer close(output)

	waitGroup := sync.WaitGroup{}

	readInput := func(input chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case data, ok := <-input:
				if !ok {
					return
				}

				if strings.Contains(data, noMultiplexer) {
					continue
				}

				select {
				case <-ctx.Done():
					return
				case output <- data:
				}
			}
		}
	}

	for _, ch := range inputs {
		waitGroup.Add(1)

		go readInput(ch)
	}

	waitGroup.Wait()

	return nil
}

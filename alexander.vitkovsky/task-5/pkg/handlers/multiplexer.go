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

	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	readFunc := func(ch chan string) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case message, ok := <-ch:
				if !ok {
					return
				}

				if strings.Contains(message, noMultiplexer) {
					continue
				}

				select {
				case <-ctx.Done():
					return
				case output <- message:
				}
			}
		}
	}

	for _, input := range inputs {
		go readFunc(input)
	}

	wg.Wait()
	return nil
}

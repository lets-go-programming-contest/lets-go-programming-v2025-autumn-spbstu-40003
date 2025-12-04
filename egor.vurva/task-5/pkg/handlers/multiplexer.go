package handlers

import (
	"context"
	"strings"
	"sync"
)

const mulFlag = "no multiplexer"

func MultiplexerFunc(
	ctx context.Context,
	input []chan string,
	output chan string,
) error {
	waitGroup := sync.WaitGroup{}

	for _, channel := range input {
		waitGroup.Add(1)

		go func(input chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, ok := <-input:
					if !ok {
						return
					}

					if !strings.Contains(value, mulFlag) {
						select {
						case <-ctx.Done():
							return
						case output <- value:
						}
					}
				}
			}
		}(channel)
	}

	waitGroup.Wait()

	return nil
}

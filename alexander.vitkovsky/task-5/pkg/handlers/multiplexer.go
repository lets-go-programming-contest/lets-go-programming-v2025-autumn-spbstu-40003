package handlers

import (
	"context"
	"strings"
	"sync"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, ch := range inputs {
		go func(in chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-in:
					if !ok {
						return
					}
					if strings.Contains(msg, "no multiplexer") {
						continue
					}
					select {
					case <-ctx.Done():
						return
					case output <- msg:
					}
				}
			}
		}(ch)
	}

	wg.Wait()
	close(output)
	return nil
}

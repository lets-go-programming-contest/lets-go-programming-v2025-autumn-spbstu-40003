package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(item, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}

			prefix := "decorated: "
			if !strings.HasPrefix(item, prefix) {
				item = prefix + item
			}

			select {
			case output <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return nil
	}

	var i int
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}

			targetCh := outputs[i%len(outputs)]
			i++

			select {
			case targetCh <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var wg sync.WaitGroup

	readCh := func(ch chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-ch:
				if !ok {
					return
				}

				if strings.Contains(item, "no multiplexer") {
					continue
				}

				select {
				case output <- item:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, ch := range inputs {
		wg.Add(1)
		go readCh(ch)
	}

	wg.Wait()
	return nil
}

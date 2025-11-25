package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(item, "no decorator") {
				return ErrCantDecorate
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

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	readCh := func(c chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-c:
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

	for _, channel := range inputs {
		waitGroup.Add(1)

		go readCh(channel)
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}

			targetIndex := counter % len(outputs)
			targetCh := outputs[targetIndex]
			counter++

			select {
			case targetCh <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

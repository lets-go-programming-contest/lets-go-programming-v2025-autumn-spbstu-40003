package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

const (
	noDecorator   = "no decorator"
	noMultiplexer = "no multiplexer"
	decorated     = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case word, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(word, noDecorator) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(word, decorated) {
				word = decorated + word
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- word:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	outputChanNum := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case word, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[outputChanNum] <- word:
				outputChanNum = (outputChanNum + 1) % len(outputs)
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	waitGroup := sync.WaitGroup{}
	for _, inputChan := range inputs {
		waitGroup.Add(1)

		go func(inputChan chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case word, ok := <-inputChan:
					if !ok {
						return
					}

					if strings.Contains(word, noMultiplexer) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- word:
					}
				}
			}
		}(inputChan)
	}

	waitGroup.Wait()

	return nil
}

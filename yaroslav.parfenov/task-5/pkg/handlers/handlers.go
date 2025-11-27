package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
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

			if strings.Contains(word, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(word, "decorated: ") {
				word = "decorated: " + word
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
	var outputChanNum = 0
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
	wg := sync.WaitGroup{}
	for _, inputChan := range inputs {
		wg.Add(1)

		go func(inputChan chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case word, ok := <-inputChan:
					if !ok {
						return
					}

					if strings.Contains(word, "no multiplexer") {
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

	wg.Wait()
	return nil
}

package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
	ErrEmptyOutputs    = errors.New("empty outputs")
)

const (
	blockDecorator   = "no decorator"
	blockMultiplexer = "no multiplexer"
	prefix           = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case msg, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(msg, blockDecorator) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(msg, prefix) {
				msg = prefix + msg
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- msg:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case msg, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[idx] <- msg:
				idx = (idx + 1) % len(outputs)
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	for _, ch := range inputs {
		wg.Add(1)
		src := ch

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case msg, ok := <-src:
					if !ok {
						return
					}

					if strings.Contains(msg, blockMultiplexer) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- msg:
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

package handlers

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	if len(inputs) == 0 {
		return nil
	}

	done := make(chan struct{})
	defer close(done)

	for _, inputChan := range inputs {
		go func(in chan string) {
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					case <-done:
						return
					}
				}
			}
		}(inputChan)
	}

	<-ctx.Done()

	return ctx.Err()
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return nil
	}

	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := int(atomic.AddUint64(&counter, 1)-1) % len(outputs)

			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

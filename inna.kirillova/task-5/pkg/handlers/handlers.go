package handlers

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/sync/errgroup"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
	ErrEmptyOutputs    = errors.New("outputs must not be empty")
)

const (
	noDecorator     = "no decorator"
	noMultiplexer   = "no multiplexer"
	decoratedPrefix = "decorated:"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(data, noDecorator) {
				return ErrCantBeDecorated
			}
			if !strings.HasPrefix(data, decoratedPrefix) {
				data = decoratedPrefix + data
			}
			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			select {
			case outputs[index] <- data:
			case <-ctx.Done():
				return nil
			}
			index = (index + 1) % len(outputs)
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, in := range inputs {
		in := in
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case data, ok := <-in:
					if !ok {
						return nil
					}
					if strings.Contains(data, noMultiplexer) {
						continue
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return nil
					}
				}
			}
		})
	}

	return g.Wait()
}

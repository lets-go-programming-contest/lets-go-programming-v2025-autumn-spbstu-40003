package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrCannotBeDecorated             = errors.New("can't be decorated")
	ErrNoOutputChannelsForSeparator  = errors.New("no output channels provided for separator")
	ErrNoInputChannelsForMultiplexer = errors.New("no input channels provided for multiplexer")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", ctx.Err())

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("%w", ErrCannotBeDecorated)
			}

			if !strings.HasPrefix(data, prefix) {
				select {
				case <-ctx.Done():
					return fmt.Errorf("%w", ctx.Err())
				case output <- prefix + data:
				}
			} else {
				select {
				case <-ctx.Done():
					return fmt.Errorf("%w", ctx.Err())
				case output <- data:
				}
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return fmt.Errorf("%w", ErrNoOutputChannelsForSeparator)
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", ctx.Err())

		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return fmt.Errorf("%w", ctx.Err())
			case outputs[index%len(outputs)] <- data:
				index++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("%w", ErrNoInputChannelsForMultiplexer)
	}

	var waitGroup sync.WaitGroup
	errCh := make(chan error, len(inputs))

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		processChannel := func(inputChan chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case data, ok := <-inputChan:
					if !ok {
						return
					}

					if !strings.Contains(data, "no multiplexer") {
						select {
						case <-ctx.Done():
							return
						case output <- data:
						}
					}
				}
			}
		}

		go processChannel(inputChannel)
	}

	go func() {
		waitGroup.Wait()
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", ctx.Err())

	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("multiplexer error: %w", err)
		}

		return nil
	}
}

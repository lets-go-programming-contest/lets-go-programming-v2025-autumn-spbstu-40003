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

func processChannelData(ctx context.Context, data string, output chan string) bool {
	if strings.Contains(data, "no multiplexer") {
		return true
	}

	select {
	case <-ctx.Done():
		return false
	case output <- data:
		return true
	}
}

func processInputChannel(ctx context.Context, inputChan chan string, output chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case data, ok := <-inputChan:
			if !ok {
				return
			}

			if !processChannelData(ctx, data, output) {
				return
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("%w", ErrNoInputChannelsForMultiplexer)
	}

	var waitGroup sync.WaitGroup
	errorChannel := make(chan error, len(inputs))

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		go processInputChannel(ctx, inputChannel, output, &waitGroup)
	}

	go func() {
		waitGroup.Wait()
		close(errorChannel)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", ctx.Err())

	case err := <-errorChannel:
		if err != nil {
			return fmt.Errorf("multiplexer error: %w", err)
		}

		return nil
	}
}

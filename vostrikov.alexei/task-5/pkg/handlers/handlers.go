package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

const (
	DecoratedPrefix       = "decorated: "
	SkipMultiplexerMarker = "no multiplexer"
	SkipDecoratorMarker   = "no decorator"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	inputChannel chan string,
	outputChannel chan string,
) error {
	defer close(outputChannel)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("prefix decorator context canceled: %w", ctx.Err())
		case value, hasValue := <-inputChannel:
			if !hasValue {
				return nil
			}

			if strings.Contains(value, SkipDecoratorMarker) {
				return fmt.Errorf("%w: %s", ErrCantBeDecorated, value)
			}

			if !strings.HasPrefix(value, DecoratedPrefix) {
				value = DecoratedPrefix + value
			}

			select {
			case <-ctx.Done():
				return fmt.Errorf("prefix decorator context canceled: %w", ctx.Err())
			case outputChannel <- value:
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	inputChannel chan string,
	outputChannels []chan string,
) error {
	if len(outputChannels) == 0 {
		return drainInput(ctx, inputChannel)
	}

	return distributeRoundRobin(ctx, inputChannel, outputChannels)
}

func drainInput(
	ctx context.Context,
	inputChannel chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("separator canceled: %w", ctx.Err())

		case _, hasValue := <-inputChannel:
			if !hasValue {
				return nil
			}
		}
	}
}

func distributeRoundRobin(
	ctx context.Context,
	inputChannel chan string,
	outputChannels []chan string,
) error {
	defer closeAll(outputChannels)

	currentIndex := 0
	outputCount := len(outputChannels)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("separator canceled: %w", ctx.Err())

		case value, hasValue := <-inputChannel:
			if !hasValue {
				return nil
			}

			targetChannel := nextChannel(outputChannels, &currentIndex, outputCount)

			if err := sendValue(ctx, targetChannel, value); err != nil {
				return err
			}
		}
	}
}

func nextChannel(
	outputChannels []chan string,
	currentIndex *int,
	outputCount int,
) chan string {
	channel := outputChannels[*currentIndex%outputCount]
	*currentIndex++

	return channel
}

func sendValue(
	ctx context.Context,
	targetChannel chan string,
	value string,
) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("separator canceled: %w", ctx.Err())
	case targetChannel <- value:
		return nil
	}
}

func closeAll(channels []chan string) {
	for _, ch := range channels {
		close(ch)
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputChannels []chan string,
	outputChannel chan string,
) error {
	defer close(outputChannel)

	if len(inputChannels) == 0 {
		return nil
	}

	mergedChannel := mergeInputs(ctx, inputChannels)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("multiplexer context canceled: %w", ctx.Err())
		case value, hasValue := <-mergedChannel:
			if !hasValue {
				return nil
			}

			if strings.Contains(value, SkipMultiplexerMarker) {
				continue
			}

			select {
			case <-ctx.Done():
				return fmt.Errorf("multiplexer context canceled: %w", ctx.Err())
			case outputChannel <- value:
			}
		}
	}
}

func mergeInputs(
	ctx context.Context,
	inputChannels []chan string,
) chan string {
	mergedChannel := make(chan string)

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputChannels))

	for _, inputChannel := range inputChannels {
		localInputChannel := inputChannel

		mergeWorker := func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case value, hasValue := <-localInputChannel:
					if !hasValue {
						return
					}

					select {
					case <-ctx.Done():
						return
					case mergedChannel <- value:
					}
				}
			}
		}

		go mergeWorker()
	}

	waiter := func() {
		waitGroup.Wait()
		close(mergedChannel)
	}

	go waiter()

	return mergedChannel
}

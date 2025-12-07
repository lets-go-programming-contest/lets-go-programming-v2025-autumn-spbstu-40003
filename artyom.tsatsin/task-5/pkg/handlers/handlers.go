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

func PrefixDecoratorFunc(
	ctx context.Context,
	inputChannel chan string,
	outputChannel chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case message, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(message, blockDecorator) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(message, prefix) {
				message = prefix + message
			}

			select {
			case <-ctx.Done():
				return nil

			case outputChannel <- message:
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
		return ErrEmptyOutputs
	}

	currentIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case message, ok := <-inputChannel:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil

			case outputChannels[currentIndex] <- message:
				currentIndex = (currentIndex + 1) % len(outputChannels)
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputChannels []chan string,
	outputChannel chan string,
) error {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputChannels))

	for _, inputChannel := range inputChannels {
		source := inputChannel

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case message, ok := <-source:
					if !ok {
						return
					}

					if strings.Contains(message, blockMultiplexer) {
						continue
					}

					select {
					case <-ctx.Done():
						return

					case outputChannel <- message:
					}
				}
			}
		}()
	}

	waitGroup.Wait()

	return nil
}

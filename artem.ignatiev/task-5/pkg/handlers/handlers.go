package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	inputChannel chan string,
	outputChannel chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case item, ok := <-inputChannel:
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
			case outputChannel <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	inputChannel chan string,
	outputsChannels []chan string,
) error {
	if len(outputsChannels) == 0 {
		return nil
	}

	var index int

	for {
		select {
		case <-ctx.Done():
			return nil

		case item, ok := <-inputChannel:
			if !ok {
				return nil
			}

			targetChannel := outputsChannels[index%len(outputsChannels)]
			index++

			select {
			case targetChannel <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputsChannels []chan string,
	outputChannel chan string,
) error {
	var waitGroup sync.WaitGroup

	readFunc := func(channel chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case item, ok := <-channel:
				if !ok {
					return
				}

				if strings.Contains(item, "no multiplexer") {
					continue
				}

				select {
				case outputChannel <- item:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, channel := range inputsChannels {
		waitGroup.Add(1)

		go readFunc(channel)
	}

	waitGroup.Wait()

	return nil
}

package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantDecorate           = errors.New("can't decorate value")
	ErrContextDoneInDecorator = errors.New("decorator context canceled")
	ErrContextDoneInSeparator = errors.New("separator context canceled")
)

const (
	StrDecorated   = "decorated: "
	StrNoDecorator = "no decorator"
	StrNoMult      = "no multiplexer"
)

func PrefixDecoratorFunc(
	internalContext context.Context,
	inputChannel chan string,
	outputChannel chan string,
) error {
	for {
		select {
		case receivedValue, channelOpen := <-inputChannel:
			if !channelOpen {
				return nil
			}

			if strings.Contains(receivedValue, StrNoDecorator) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(receivedValue, StrDecorated) {
				receivedValue = StrDecorated + receivedValue
			}

			select {
			case outputChannel <- receivedValue:
			case <-internalContext.Done():
				return ErrContextDoneInDecorator
			}

		case <-internalContext.Done():
			return ErrContextDoneInDecorator
		}
	}
}

func SeparatorFunc(
	internalContext context.Context,
	inputChannel chan string,
	outputChannels []chan string,
) error {
	channelIndex := 0

	for {
		select {
		case receivedValue, channelOpen := <-inputChannel:
			if !channelOpen {
				return nil
			}

			select {
			case outputChannels[channelIndex] <- receivedValue:
			case <-internalContext.Done():
				return ErrContextDoneInSeparator
			}

			channelIndex++
			if channelIndex >= len(outputChannels) {
				channelIndex = 0
			}

		case <-internalContext.Done():
			return ErrContextDoneInSeparator
		}
	}
}

func MultiplexerFunc(
	internalContext context.Context,
	inputChannels []chan string,
	outputChannel chan string,
) error {
	var workersGroup sync.WaitGroup
	stopChannel := make(chan struct{})

	for _, inputChannel := range inputChannels {
		workersGroup.Add(1)

		go func(sourceChannel chan string) {
			defer workersGroup.Done()

			for {
				select {
				case receivedValue, channelOpen := <-sourceChannel:
					if !channelOpen {
						return
					}

					if strings.Contains(receivedValue, StrNoMult) {
						continue
					}

					select {
					case outputChannel <- receivedValue:
					case <-internalContext.Done():
						return
					case <-stopChannel:
						return
					}

				case <-internalContext.Done():
					return
				case <-stopChannel:
					return
				}
			}
		}(inputChannel)
	}

	<-internalContext.Done()
	close(stopChannel)
	workersGroup.Wait()

	return nil
}

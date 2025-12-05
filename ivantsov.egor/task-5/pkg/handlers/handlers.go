package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantDecorate           = errors.New("can't be decorated")
	ErrContextDoneInDecorator = errors.New("context done in decorator")
	ErrContextDoneInSeparator = errors.New("context done in separator")
)

const (
	StrNoDecorator = "no decorator"
	StrNoMult      = "no multiplexer"
	StrDecorated   = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, inputChannel chan string, outputChannel chan string) error {
	for {
		select {
		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(data, StrNoDecorator) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(data, StrDecorated) {
				data = StrDecorated + data
			}

			select {
			case outputChannel <- data:
			case <-ctx.Done():
				return errors.Join(ErrContextDoneInDecorator, ctx.Err())
			}
		case <-ctx.Done():
			return errors.Join(ErrContextDoneInDecorator, ctx.Err())
		}
	}
}

func SeparatorFunc(ctx context.Context, inputChannel chan string, outputChannels []chan string) error {
	outputIndex := 0

	for {
		select {
		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			select {
			case outputChannels[outputIndex] <- data:
			case <-ctx.Done():
				return errors.Join(ErrContextDoneInSeparator, ctx.Err())
			}

			outputIndex++
			if outputIndex == len(outputChannels) {
				outputIndex = 0
			}
		case <-ctx.Done():
			return errors.Join(ErrContextDoneInSeparator, ctx.Err())
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	waitGroup := &sync.WaitGroup{}
	done := make(chan struct{})

	for _, inputChannel := range inputChannels {
		channelCopy := inputChannel

		waitGroup.Add(1)

		go func(ch chan string) {
			defer waitGroup.Done()

			for {
				select {
				case data, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(data, StrNoMult) {
						continue
					}

					select {
					case outputChannel <- data:
					case <-ctx.Done():
						return
					case <-done:
						return
					}
				case <-ctx.Done():
					return
				case <-done:
					return
				}
			}
		}(channelCopy)
	}

	<-ctx.Done()
	close(done)
	waitGroup.Wait()

	return nil
}

package handlers

import (
	"context"
	"errors"
	"fmt"
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
		case payload, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(payload, StrNoDecorator) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(payload, StrDecorated) {
				payload = StrDecorated + payload
			}

			select {
			case outputChannel <- payload:
			case <-ctx.Done():
				return fmt.Errorf("%w: %w", ErrContextDoneInDecorator, ctx.Err())
			}

		case <-ctx.Done():
			return fmt.Errorf("%w: %w", ErrContextDoneInDecorator, ctx.Err())
		}
	}
}

func SeparatorFunc(ctx context.Context, inputChannel chan string, outputGroups []chan string) error {
	current := 0

	for {
		select {
		case message, ok := <-inputChannel:
			if !ok {
				return nil
			}

			select {
			case outputGroups[current] <- message:
			case <-ctx.Done():
				return fmt.Errorf("%w: %w", ErrContextDoneInSeparator, ctx.Err())
			}

			current = (current + 1) % len(outputGroups)

		case <-ctx.Done():
			return fmt.Errorf("%w: %w", ErrContextDoneInSeparator, ctx.Err())
		}
	}
}

func MultiplexerFunc(ctx context.Context, sourceChannels []chan string, outputChannel chan string) error {
	var group sync.WaitGroup
	termination := make(chan struct{})

	for _, channel := range sourceChannels {
		group.Add(1)

		go func(localChannel chan string) {
			defer group.Done()

			for {
				select {
				case msg, ok := <-localChannel:
					if !ok {
						return
					}

					if strings.Contains(msg, StrNoMult) {
						continue
					}

					select {
					case outputChannel <- msg:
					case <-ctx.Done():
						return
					case <-termination:
						return
					}

				case <-ctx.Done():
					return
				case <-termination:
					return
				}
			}
		}(channel)
	}

	<-ctx.Done()
	close(termination)
	group.Wait()

	return nil
}

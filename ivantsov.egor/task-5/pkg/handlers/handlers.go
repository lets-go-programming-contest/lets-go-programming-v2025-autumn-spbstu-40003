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

func PrefixDecoratorFunc(
	ctx context.Context,
	in chan string,
	out chan string,
) error {
	for {
		select {
		case value, ok := <-in:
			if !ok {
				return nil
			}

			if strings.Contains(value, StrNoDecorator) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(value, StrDecorated) {
				value = StrDecorated + value
			}

			select {
			case out <- value:
			case <-ctx.Done():
				return errors.Join(
					ErrContextDoneInDecorator,
					ctx.Err(),
				)
			}

		case <-ctx.Done():
			return errors.Join(
				ErrContextDoneInDecorator,
				ctx.Err(),
			)
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	in chan string,
	outs []chan string,
) error {
	pos := 0
	count := len(outs)

	for {
		select {
		case value, ok := <-in:
			if !ok {
				return nil
			}

			select {
			case outs[pos] <- value:
			case <-ctx.Done():
				return errors.Join(
					ErrContextDoneInSeparator,
					ctx.Err(),
				)
			}

			pos++

			if pos >= count {
				pos = 0
			}

		case <-ctx.Done():
			return errors.Join(
				ErrContextDoneInSeparator,
				ctx.Err(),
			)
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	ins []chan string,
	out chan string,
) error {
	var waitGroup sync.WaitGroup
	stopSignal := make(chan struct{})

	for _, source := range ins {
		waitGroup.Add(1)

		go func(ch chan string) {
			defer waitGroup.Done()

			for {
				select {
				case value, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(value, StrNoMult) {
						continue
					}

					select {
					case out <- value:
					case <-ctx.Done():
						return
					case <-stopSignal:
						return
					}

				case <-ctx.Done():
					return
				case <-stopSignal:
					return
				}
			}
		}(source)
	}

	<-ctx.Done()

	close(stopSignal)
	waitGroup.Wait()

	return nil
}

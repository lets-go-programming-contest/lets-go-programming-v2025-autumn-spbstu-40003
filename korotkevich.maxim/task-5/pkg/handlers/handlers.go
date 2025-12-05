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

func PrefixDecoratorFunc(ctx context.Context, inputChan chan string, outputCh chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-inputChan:
			if !ok {
				return nil
			}

			if strings.Contains(msg, blockDecorator) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(msg, prefix) {
				msg = prefix + msg
			}

			select {
			case <-ctx.Done():
				return nil
			case outputCh <- msg:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, inputChan chan string, outputChans []chan string) error {
	if len(outputChans) == 0 {
		return ErrEmptyOutputs
	}

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-inputChan:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputChans[idx] <- msg:
				idx = (idx + 1) % len(outputChans)
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChans []chan string, outputCh chan string) error {
	if len(inputChans) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputChans))

	for _, ch := range inputChans {
		src := ch
		go func(conv chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-conv:
					if !ok {
						return
					}

					if strings.Contains(msg, blockMultiplexer) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case outputCh <- msg:
					}
				}
			}
		}(src)
	}

	waitGroup.Wait()

	return nil
}

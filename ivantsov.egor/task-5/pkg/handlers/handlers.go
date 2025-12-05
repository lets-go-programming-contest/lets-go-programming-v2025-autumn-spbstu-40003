package handlers

import (
	"context"
	"strings"
	"sync"
)

type ProcessorFunc func(context.Context, chan string, chan string) error

func PrefixDecoratorFunc(prefix string) ProcessorFunc {
	return func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error {

		for {
			select {
			case <-ctx.Done():
				return nil

			case value, ok := <-input:
				if !ok {
					return nil
				}

				output <- prefix + value
			}
		}
	}
}

func MultiplexerFunc() func(
	context.Context,
	[]chan string,
	chan string,
) error {

	return func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error {

		var waitGroup sync.WaitGroup

		for _, inputChannel := range inputs {
			waitGroup.Add(1)

			go func(channel chan string) {
				defer waitGroup.Done()

				for {
					select {
					case <-ctx.Done():
						return

					case value, ok := <-channel:
						if !ok {
							return
						}

						output <- value
					}
				}

			}(inputChannel)
		}

		waitGroup.Wait()
		return nil
	}
}

func SeparatorFunc(separator string) func(
	context.Context,
	chan string,
	[]chan string,
) error {

	return func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error {

		for {
			select {
			case <-ctx.Done():
				return nil

			case value, ok := <-input:
				if !ok {
					return nil
				}

				parts := strings.Split(value, separator)

				for index, part := range parts {
					outputs[index] <- part
				}
			}
		}
	}
}

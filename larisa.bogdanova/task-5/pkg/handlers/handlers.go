package handlers

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var errCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, inputCh chan string, outputCh chan string) error {
	defer close(outputCh)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-inputCh:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("%w", errCannotBeDecorated)
			}
			if !strings.HasPrefix(data, "decorated:") {
				data = "decorated: " + data
			}

			select {
			case outputCh <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, inputCh chan string, outputs []chan string) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-inputCh:
			if !ok {
				return nil
			}

			target := outputs[index%len(outputs)]

			select {
			case target <- data:
				index++
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, outputCh chan string) error {
	defer close(outputCh)

	cases := make([]reflect.SelectCase, len(inputs)+1)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())}

	for i, ch := range inputs {
		cases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	aliveInputs := len(inputs)
	for aliveInputs > 0 {
		chosen, value, ok := reflect.Select(cases)

		if chosen == 0 {
			return ctx.Err()
		}

		if !ok {
			cases[chosen].Chan = reflect.ValueOf(nil)
			aliveInputs--
			continue
		}

		data := value.String()
		if strings.Contains(data, "no multiplexer") {
			continue
		}

		select {
		case outputCh <- data:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

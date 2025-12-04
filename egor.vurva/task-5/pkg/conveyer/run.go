package conveyer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var errHandlersNotRegistered = errors.New("no handlers found")

func (conv *Conv) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	hasHandlers := len(conv.decorators)+len(conv.multiplexers)+len(conv.separators) != 0
	if !hasHandlers {
		return errHandlersNotRegistered
	}

	for _, decorator := range conv.decorators {
		dec := decorator.fn
		input := conv.inputs[decorator.input]
		output := conv.outputs[decorator.output]

		group.Go(func() error {
			return dec(groupCtx, input, output)
		})
	}

	for _, separator := range conv.separators {
		sep := separator.fn
		input, outputs := conv.getSeparatorChannels(separator)

		group.Go(func() error {
			return sep(groupCtx, input, outputs)
		})
	}

	for _, multiplexer := range conv.multiplexers {
		mul := multiplexer.fn
		inputs, output := conv.getMultiplexerChannels(multiplexer)

		group.Go(func() error {
			return mul(groupCtx, inputs, output)
		})
	}

	err := group.Wait()

	conv.closeAllChannels()

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (conv *Conv) closeAllChannels() {
	res := make(map[chan string]string)

	for id, ch := range conv.inputs {
		res[ch] = id
	}

	for id, ch := range conv.outputs {
		if _, ok := res[ch]; !ok {
			res[ch] = id
		}
	}

	for ch := range res {
		close(ch)
	}
}

func (conv *Conv) getSeparatorChannels(sep SeparatorStorage) (chan string, []chan string) {
	input := conv.inputs[sep.input]

	outputs := make([]chan string, 0, len(sep.output))

	for _, id := range sep.output {
		outputs = append(outputs, conv.outputs[id])
	}

	return input, outputs
}

func (conv *Conv) getMultiplexerChannels(mul MultiplexerStorage) ([]chan string, chan string) {
	output := conv.outputs[mul.output]

	input := make([]chan string, 0, len(mul.input))

	for _, id := range mul.input {
		input = append(input, conv.inputs[id])
	}

	return input, output
}

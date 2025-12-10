package conveyer

import (
	"context"
)

// registers for all handlers

func (conv *Conveyer) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	conv.CreateChannel(input, output)

	conv.AddHandler(func(ctx context.Context) error {
		return handler(ctx, conv.channels[input], conv.channels[output])
	})
}

func (conv *Conveyer) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	conv.CreateChannel(input)
	conv.CreateChannel(outputs...)

	conv.AddHandler(func(ctx context.Context) error {
		out := make([]chan string, len(outputs))
		for index, name := range outputs {
			out[index] = conv.channels[name]
		}
		return handler(ctx, conv.channels[input], out)
	})
}

func (conv *Conveyer) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	conv.CreateChannel(inputs...)
	conv.CreateChannel(output)

	conv.AddHandler(func(ctx context.Context) error {
		in := make([]chan string, len(inputs))
		for index, name := range inputs {
			in[index] = conv.channels[name]
		}
		return handler(ctx, in, conv.channels[output])
	})
}

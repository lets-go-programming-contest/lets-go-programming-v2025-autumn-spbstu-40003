package conveyer

import "context"

func (c *Conveyer) RegisterDecorator(
	handler func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	c.createChannels(input)
	c.createChannels(output)
	c.addHandler(func(ctx context.Context) error {
		return handler(ctx, c.channels[input], c.channels[output])
	})
}

func (c *Conveyer) RegisterSeparator(
	handler func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	c.createChannels(input)
	c.createChannels(outputs...)
	c.addHandler(func(ctx context.Context) error {
		outputChannels := make([]chan string, 0, len(outputs))

		for _, ch := range outputs {
			outputChannels = append(outputChannels, c.channels[ch])
		}

		return handler(ctx, c.channels[input], outputChannels)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	c.createChannels(inputs...)
	c.createChannels(output)
	c.addHandler(func(ctx context.Context) error {
		inputChannels := make([]chan string, 0, len(inputs))

		for _, ch := range inputs {
			inputChannels = append(inputChannels, c.channels[ch])
		}

		return handler(ctx, inputChannels, c.channels[output])
	})
}

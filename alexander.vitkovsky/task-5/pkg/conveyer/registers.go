package conveyer

// handlers registrars

func (conv *Conveyer) RegisterDecorator(fn decoratorFn, input string, output string) {
	conv.handlers = append(conv.handlers, handlerConfig{
		kind:     "decorator",
		function: fn,
		inputs:   []string{input},
		outputs:  []string{output},
	})

	conv.getOrCreateChannel(input)
	conv.getOrCreateChannel(output)
}

func (conv *Conveyer) RegisterMultiplexer(fn multiplexerFn, inputs []string, output string) {
	conv.handlers = append(conv.handlers, handlerConfig{
		kind:     "multiplexer",
		function: fn,
		inputs:   inputs,
		outputs:  []string{output},
	})

	for _, input := range inputs {
		conv.getOrCreateChannel(input)
	}
	conv.getOrCreateChannel(output)
}

func (conv *Conveyer) RegisterSeparator(fn separatorFn, input string, outputs []string) {
	conv.handlers = append(conv.handlers, handlerConfig{
		kind:     "separator",
		function: fn,
		inputs:   []string{input},
		outputs:  outputs,
	})

	conv.getOrCreateChannel(input)
	for _, output := range outputs {
		conv.getOrCreateChannel(output)
	}
}

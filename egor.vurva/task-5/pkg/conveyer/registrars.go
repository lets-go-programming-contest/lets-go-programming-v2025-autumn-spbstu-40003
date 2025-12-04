package conveyer

import "context"

func (conv *Conv) RegisterDecorator(
	decFunc func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	conv.ensureChan(input)
	conv.ensureChan(output)

	conv.decorators = append(conv.decorators,
		DecoratorStorage{
			fn:     decFunc,
			input:  input,
			output: output,
		})
}

func (conv *Conv) RegisterSeparator(
	sepFunc func(
		ctx context.Context,
		input chan string,
		output []chan string,
	) error,
	input string,
	output []string,
) {
	conv.ensureChan(input)

	for _, outChan := range output {
		conv.ensureChan(outChan)
	}

	conv.separators = append(conv.separators,
		SeparatorStorage{
			fn:     sepFunc,
			input:  input,
			output: output,
		})
}

func (conv *Conv) RegisterMultiplexer(
	mulFunc func(
		ctx context.Context,
		input []chan string,
		output chan string,
	) error,
	input []string,
	output string,
) {
	for _, inChan := range input {
		conv.ensureChan(inChan)
	}

	conv.ensureChan(output)

	conv.multiplexers = append(conv.multiplexers,
		MultiplexerStorage{
			fn:     mulFunc,
			input:  input,
			output: output,
		})
}

func (conv *Conv) ensureChan(name string) {
	if ch, ok := conv.inputs[name]; ok {
		if _, exists := conv.outputs[name]; !exists {
			conv.outputs[name] = ch
		}

		return
	}

	if ch, ok := conv.outputs[name]; ok {
		if _, exists := conv.inputs[name]; !exists {
			conv.inputs[name] = ch
		}

		return
	}

	ch := make(chan string, conv.chanSize)
	conv.inputs[name] = ch
	conv.outputs[name] = ch
}

package conveyer

import (
	"context"
	"fmt"

	"github.com/gituser549/task-5/pkg/chansyncmap"
	"golang.org/x/sync/errgroup"
)

type Conveyer struct {
	size int

	channelsById *chansyncmap.ChanSyncMap
	decorators   []Decorator
	multiplexers []Multiplexer
	separators   []Separator
}

type Decorator struct {
	function func(ctx context.Context, input chan string, output chan string) error
	input    chan string
	output   chan string
}

type Multiplexer struct {
	function func(ctx context.Context, inputs []chan string, output chan string) error
	inputs   []chan string
	output   chan string
}

type Separator struct {
	function func(ctx context.Context, input chan string, outputs []chan string) error
	input    chan string
	outputs  []chan string
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:         size,
		channelsById: chansyncmap.New(size),
		decorators:   make([]Decorator, 0),
		multiplexers: make([]Multiplexer, 0),
		separators:   make([]Separator, 0),
	}
}

func (conv *Conveyer) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChan := conv.channelsById.GetOrCreateChan(input)

	outputChan := conv.channelsById.GetOrCreateChan(output)

	conv.decorators = append(conv.decorators, Decorator{fn, inputChan, outputChan})
}

func (conv *Conveyer) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error,
	input []string,
	output string,
) {
	var inputChans []chan string
	for _, input := range input {
		inputChans = append(inputChans, conv.channelsById.GetOrCreateChan(input))
	}

	outputChan := conv.channelsById.GetOrCreateChan(output)

	conv.multiplexers = append(conv.multiplexers, Multiplexer{fn, inputChans, outputChan})
}

func (conv *Conveyer) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := conv.channelsById.GetOrCreateChan(input)

	var outputChans []chan string
	for _, output := range outputs {
		outputChans = append(outputChans, conv.channelsById.GetOrCreateChan(output))
	}

	conv.separators = append(conv.separators, Separator{fn, inputChan, outputChans})
}

func (conv *Conveyer) Send(input string, data string) error {
	if inputChan, ok := conv.channelsById.GetChan(input); ok {
		inputChan <- data

		return nil
	}

	return fmt.Errorf("chan not found")
}

func (conv *Conveyer) Recv(output string) (string, error) {
	if outputChan, ok := conv.channelsById.GetChan(output); ok {
		data, ok := <-outputChan
		if !ok {
			return "undefined", nil
		}

		return data, nil
	}

	return "", fmt.Errorf("chan not found")
}

func (conv *Conveyer) Run(ctx context.Context) error {
	defer conv.channelsById.CloseAllChans()

	group, groupCtx := errgroup.WithContext(ctx)

	for _, decorator := range conv.decorators {
		fn, inputChan, outputChan := decorator.function, decorator.input, decorator.output

		group.Go(func() error {
			return fn(groupCtx, inputChan, outputChan)
		})
	}

	for _, multiplexer := range conv.multiplexers {
		group.Go(func() error {
			return multiplexer.function(groupCtx, multiplexer.inputs, multiplexer.output)
		})
	}

	for _, separator := range conv.separators {
		group.Go(func() error {
			return separator.function(groupCtx, separator.input, separator.outputs)
		})
	}

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("conveyer was shut down with error: %v", err)
	}

	return nil
}

package conveyer

import (
	"context"
	"errors"
	"fmt"

	"github.com/gituser549/task-5/pkg/chansyncmap"
	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type Conveyer struct {
	size int

	channelsByID *chansyncmap.ChanSyncMap
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
		channelsByID: chansyncmap.New(size),
		decorators:   make([]Decorator, 0),
		multiplexers: make([]Multiplexer, 0),
		separators:   make([]Separator, 0),
	}
}

func (conv *Conveyer) RegisterDecorator(function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChan := conv.channelsByID.GetOrCreateChan(input)

	outputChan := conv.channelsByID.GetOrCreateChan(output)

	conv.decorators = append(conv.decorators, Decorator{function, inputChan, outputChan})
}

func (conv *Conveyer) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	input []string,
	output string,
) {
	inputChans := make([]chan string, 0, len(input))
	for _, input := range input {
		inputChans = append(inputChans, conv.channelsByID.GetOrCreateChan(input))
	}

	outputChan := conv.channelsByID.GetOrCreateChan(output)

	conv.multiplexers = append(conv.multiplexers, Multiplexer{function, inputChans, outputChan})
}

func (conv *Conveyer) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := conv.channelsByID.GetOrCreateChan(input)

	outputChans := make([]chan string, 0, len(outputs))
	for _, output := range outputs {
		outputChans = append(outputChans, conv.channelsByID.GetOrCreateChan(output))
	}

	conv.separators = append(conv.separators, Separator{function, inputChan, outputChans})
}

func (conv *Conveyer) Send(input string, data string) error {
	if inputChan, ok := conv.channelsByID.GetChan(input); ok {
		inputChan <- data

		return nil
	}

	return ErrChanNotFound
}

func (conv *Conveyer) Recv(output string) (string, error) {
	if outputChan, ok := conv.channelsByID.GetChan(output); ok {
		data, ok := <-outputChan
		if !ok {
			return undefined, nil
		}

		return data, nil
	}

	return "", ErrChanNotFound
}

func (conv *Conveyer) Run(ctx context.Context) error {
	defer conv.channelsByID.CloseAllChans()

	group, groupCtx := errgroup.WithContext(ctx)

	for _, decorator := range conv.decorators {
		group.Go(func() error {
			return decorator.function(groupCtx, decorator.input, decorator.output)
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
		return fmt.Errorf("conveyer was shut down with error: %w", err)
	}

	return nil
}

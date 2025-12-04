package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrorChanNotFound = errors.New("chan not found")

type MyConveyer struct {
	size         int
	channels     map[string]chan string
	mutex        sync.RWMutex
	decorators   []decorator
	multiplexers []multiplexer
	separators   []separator
}

func New(size int) *MyConveyer {
	return &MyConveyer{
		size:         size,
		channels:     make(map[string]chan string),
		mutex:        sync.RWMutex{},
		decorators:   make([]decorator, 0),
		multiplexers: make([]multiplexer, 0),
		separators:   make([]separator, 0),
	}
}

type decorator struct {
	function func(ctx context.Context, input chan string, output chan string) error
	input    string
	output   string
}

type multiplexer struct {
	function func(ctx context.Context, inputs []chan string, output chan string) error
	inputs   []string
	output   string
}

type separator struct {
	function func(ctx context.Context, input chan string, outputs []chan string) error
	input    string
	outputs  []string
}

func (conv *MyConveyer) getChannel(id string) chan string {
	conv.mutex.Lock()
	defer conv.mutex.Unlock()

	if channel, exists := conv.channels[id]; exists {
		return channel
	}

	channel := make(chan string, conv.size)
	conv.channels[id] = channel

	return channel
}

func (conv *MyConveyer) RegisterDecorator(
	fn func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	conv.getChannel(input)
	conv.getChannel(output)

	conv.decorators = append(conv.decorators, decorator{fn, input, output})
}

func (conv *MyConveyer) RegisterMultiplexer(
	fn func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	for _, input := range inputs {
		conv.getChannel(input)
	}
	conv.getChannel(output)

	conv.multiplexers = append(conv.multiplexers, multiplexer{fn, inputs, output})
}

func (conv *MyConveyer) RegisterSeparator(
	fn func(
		ctx context.Context,
		input chan string,
		outputs []chan string) error,
	input string,
	outputs []string,
) {
	for _, output := range outputs {
		conv.getChannel(output)
	}
	conv.getChannel(input)

	conv.separators = append(conv.separators, separator{fn, input, outputs})
}

func (conv *MyConveyer) close() {
	conv.mutex.Lock()
	defer conv.mutex.Unlock()

	for _, channel := range conv.channels {
		close(channel)
	}
}

func (conv *MyConveyer) Run(ctx context.Context) error {
	defer conv.close()

	group, groupContext := errgroup.WithContext(ctx)

	for _, decorator := range conv.decorators {
		input := conv.getChannel(decorator.input)
		output := conv.getChannel(decorator.output)

		group.Go(func() error {
			return decorator.function(groupContext, input, output)
		})
	}

	for _, multiplexer := range conv.multiplexers {
		output := conv.getChannel(multiplexer.output)
		inputs := make([]chan string, len(multiplexer.inputs))

		for index, name := range multiplexer.inputs {
			inputs[index] = conv.getChannel(name)
		}

		group.Go(func() error {
			return multiplexer.function(groupContext, inputs, output)
		})
	}

	for _, separator := range conv.separators {
		input := conv.getChannel(separator.input)
		outputs := make([]chan string, len(separator.outputs))

		for index, name := range separator.outputs {
			outputs[index] = conv.getChannel(name)
		}

		group.Go(func() error {
			return separator.function(groupContext, input, outputs)
		})
	}

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("running failed: %w", err)
	}

	return nil
}

func (conv *MyConveyer) Send(input string, data string) error {
	ch, exists := conv.channels[input]
	if !exists {
		return ErrorChanNotFound
	}

	ch <- data

	return nil
}

func (conv *MyConveyer) Recv(output string) (string, error) {
	ch, exists := conv.channels[output]
	if !exists {
		return "", ErrorChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

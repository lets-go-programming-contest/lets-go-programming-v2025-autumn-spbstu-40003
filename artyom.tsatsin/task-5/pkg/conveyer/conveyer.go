package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrChanFull     = errors.New("chan is full")
)

const (
	undefined = "undefined"
)

type DecoratorFunc func(
	context.Context,
	chan string,
	chan string,
) error

type MultiplexerFunc func(
	context.Context,
	[]chan string,
	chan string,
) error

type SeparatorFunc func(
	context.Context,
	chan string,
	[]chan string,
) error

type decoratorTask struct {
	decorator DecoratorFunc
	input     string
	output    string
}

type multiplexerTask struct {
	multiplexer MultiplexerFunc
	inputs      []string
	output      string
}

type separatorTask struct {
	separator SeparatorFunc
	input     string
	outputs   []string
}

type Conveyer struct {
	size         int
	mutex        sync.RWMutex
	channels     map[string]chan string
	decorators   []decoratorTask
	multiplexers []multiplexerTask
	separators   []separatorTask
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:         size,
		mutex:        sync.RWMutex{},
		channels:     make(map[string]chan string),
		decorators:   []decoratorTask{},
		multiplexers: []multiplexerTask{},
		separators:   []separatorTask{},
	}
}

func (conveyer *Conveyer) getOrCreate(channelID string) chan string {
	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	channel, exists := conveyer.channels[channelID]
	if exists {
		return channel
	}

	channel = make(chan string, conveyer.size)
	conveyer.channels[channelID] = channel

	return channel
}

func (conveyer *Conveyer) get(channelID string) (chan string, bool) {
	conveyer.mutex.RLock()
	defer conveyer.mutex.RUnlock()

	channel, exists := conveyer.channels[channelID]

	return channel, exists
}

func (conveyer *Conveyer) RegisterDecorator(
	function DecoratorFunc,
	input string,
	output string,
) {
	conveyer.getOrCreate(input)
	conveyer.getOrCreate(output)

	task := decoratorTask{
		decorator: function,
		input:     input,
		output:    output,
	}

	conveyer.decorators = append(conveyer.decorators, task)
}

func (conveyer *Conveyer) RegisterMultiplexer(
	function MultiplexerFunc,
	inputs []string,
	output string,
) {
	for _, inputName := range inputs {
		conveyer.getOrCreate(inputName)
	}

	conveyer.getOrCreate(output)

	task := multiplexerTask{
		multiplexer: function,
		inputs:      inputs,
		output:      output,
	}

	conveyer.multiplexers = append(conveyer.multiplexers, task)
}

func (conveyer *Conveyer) RegisterSeparator(
	function SeparatorFunc,
	input string,
	outputs []string,
) {
	conveyer.getOrCreate(input)

	for _, outputName := range outputs {
		conveyer.getOrCreate(outputName)
	}

	task := separatorTask{
		separator: function,
		input:     input,
		outputs:   outputs,
	}

	conveyer.separators = append(conveyer.separators, task)
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
	group, groupContext := errgroup.WithContext(ctx)

	for _, decorator := range conveyer.decorators {
		inputChannel := conveyer.getOrCreate(decorator.input)
		outputChannel := conveyer.getOrCreate(decorator.output)
		decoratorFunction := decorator.decorator

		group.Go(func() error {
			return decoratorFunction(groupContext, inputChannel, outputChannel)
		})
	}

	for _, multiplexer := range conveyer.multiplexers {
		outputChannel := conveyer.getOrCreate(multiplexer.output)

		inputChannels := make([]chan string, len(multiplexer.inputs))
		for index, inputName := range multiplexer.inputs {
			inputChannels[index] = conveyer.getOrCreate(inputName)
		}

		multiplexerFunction := multiplexer.multiplexer

		group.Go(func() error {
			return multiplexerFunction(groupContext, inputChannels, outputChannel)
		})
	}

	for _, separator := range conveyer.separators {
		inputChannel := conveyer.getOrCreate(separator.input)

		outputChannels := make([]chan string, len(separator.outputs))
		for index, outputName := range separator.outputs {
			outputChannels[index] = conveyer.getOrCreate(outputName)
		}

		separatorFunction := separator.separator

		group.Go(func() error {
			return separatorFunction(groupContext, inputChannel, outputChannels)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("run error: %w", err)
	}

	conveyer.closeAll()

	return nil
}

func (conveyer *Conveyer) Send(input string, data string) error {
	channel, exists := conveyer.get(input)
	if !exists {
		return ErrChanNotFound
	}

	select {
	case channel <- data:
		return nil

	default:
		return ErrChanFull
	}
}

func (conveyer *Conveyer) Recv(output string) (string, error) {
	channel, exists := conveyer.get(output)
	if !exists {
		return "", ErrChanNotFound
	}

	value, received := <-channel
	if !received {
		return undefined, nil
	}

	return value, nil
}

func (conveyer *Conveyer) closeAll() {
	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	for _, channel := range conveyer.channels {
		close(channel)
	}
}

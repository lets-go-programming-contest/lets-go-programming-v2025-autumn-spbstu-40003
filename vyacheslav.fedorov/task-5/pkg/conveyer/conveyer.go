package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrStreamNotFound = errors.New("chan not found")

const undef = "undefined"

type Pipeline struct {
	bufferSize   int
	streams      map[string]chan string
	mutex        sync.RWMutex
	decorators   []decoratorEntry
	multiplexers []multiplexerEntry
	separators   []separatorEntry
}

type decoratorEntry struct {
	function func(ctx context.Context, input chan string, output chan string) error
	input    string
	output   string
}

type multiplexerEntry struct {
	function func(ctx context.Context, inputs []chan string, output chan string) error
	inputs   []string
	output   string
}

type separatorEntry struct {
	function func(ctx context.Context, input chan string, outputs []chan string) error
	input    string
	outputs  []string
}

func New(size int) *Pipeline {
	return &Pipeline{
		bufferSize:   size,
		streams:      make(map[string]chan string),
		mutex:        sync.RWMutex{},
		decorators:   make([]decoratorEntry, 0),
		multiplexers: make([]multiplexerEntry, 0),
		separators:   make([]separatorEntry, 0),
	}
}

func (p *Pipeline) Send(input string, data string) error {
	ch, exists := p.getStream(input)
	if !exists {
		return ErrStreamNotFound
	}

	ch <- data

	return nil
}

func (p *Pipeline) Recv(output string) (string, error) {
	ch, exists := p.getStream(output)
	if !exists {
		return "", ErrStreamNotFound
	}

	data, ok := <-ch
	if !ok {
		return undef, nil
	}

	return data, nil
}

func (p *Pipeline) getOrCreateStream(name string) chan string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if ch, exists := p.streams[name]; exists {
		return ch
	}

	ch := make(chan string, p.bufferSize)
	p.streams[name] = ch

	return ch
}

func (p *Pipeline) getStream(name string) (chan string, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	ch, exists := p.streams[name]

	return ch, exists
}

func (p *Pipeline) closeAllStreams() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, ch := range p.streams {
		close(ch)
	}
}

func (p *Pipeline) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	p.getOrCreateStream(input)
	p.getOrCreateStream(output)

	p.decorators = append(p.decorators, decoratorEntry{
		function: function,
		input:    input,
		output:   output,
	})
}

func (p *Pipeline) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, name := range inputs {
		p.getOrCreateStream(name)
	}

	p.getOrCreateStream(output)

	p.multiplexers = append(p.multiplexers, multiplexerEntry{
		function: function,
		inputs:   inputs,
		output:   output,
	})
}

func (p *Pipeline) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	p.getOrCreateStream(input)

	for _, name := range outputs {
		p.getOrCreateStream(name)
	}

	p.separators = append(p.separators, separatorEntry{
		function: function,
		input:    input,
		outputs:  outputs,
	})
}

func (p *Pipeline) Run(ctx context.Context) error {
	defer p.closeAllStreams()

	workerGroup, workerGroupCtx := errgroup.WithContext(ctx)

	for _, decorator := range p.decorators {
		inputChan := p.getOrCreateStream(decorator.input)
		outputChan := p.getOrCreateStream(decorator.output)

		workerGroup.Go(func() error {
			return decorator.function(workerGroupCtx, inputChan, outputChan)
		})
	}

	for _, multiplexer := range p.multiplexers {
		outputChan := p.getOrCreateStream(multiplexer.output)
		inputChannels := make([]chan string, len(multiplexer.inputs))

		for index, name := range multiplexer.inputs {
			inputChannels[index] = p.getOrCreateStream(name)
		}

		workerGroup.Go(func() error {
			return multiplexer.function(workerGroupCtx, inputChannels, outputChan)
		})
	}

	for _, separator := range p.separators {
		inputChan := p.getOrCreateStream(separator.input)
		outputChannels := make([]chan string, len(separator.outputs))

		for index, name := range separator.outputs {
			outputChannels[index] = p.getOrCreateStream(name)
		}

		workerGroup.Go(func() error {
			return separator.function(workerGroupCtx, inputChan, outputChannels)
		})
	}

	err := workerGroup.Wait()
	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

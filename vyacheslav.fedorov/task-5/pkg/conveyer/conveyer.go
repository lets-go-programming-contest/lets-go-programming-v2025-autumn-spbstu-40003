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
	streams      map[string]chan string
	bufferSize   int
	decorators   []decoratorEntry
	multiplexers []multiplexerEntry
	separators   []separatorEntry
	mutex        sync.Mutex
}

type decoratorEntry struct {
	function func(ctx context.Context, input chan string, output chan string) error
	in       string
	out      string
}

type multiplexerEntry struct {
	function func(ctx context.Context, inputs []chan string, output chan string) error
	ins      []string
	out      string
}

type separatorEntry struct {
	function func(ctx context.Context, input chan string, outputs []chan string) error
	in       string
	outs     []string
}

func New(size int) *Pipeline {
	return &Pipeline{
		streams:      make(map[string]chan string),
		bufferSize:   size,
		decorators:   make([]decoratorEntry, 0),
		multiplexers: make([]multiplexerEntry, 0),
		separators:   make([]separatorEntry, 0),
		mutex:        sync.Mutex{},
	}
}

func (p *Pipeline) Send(input string, data string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	inputChannel, exists := p.streams[input]
	if !exists {
		return ErrStreamNotFound
	}
	select {
	case inputChannel <- 
		return nil
	default:
		inputChannel <- data

		return nil
	}
}

func (p *Pipeline) Recv(output string) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	outputChannel, exists := p.streams[output]
	if !exists {
		return "", ErrStreamNotFound
	}

	data, ok := <-outputChannel
	if !ok {
		return undef, nil
	}

	return data, nil
}

func (p *Pipeline) getOrCreateStream(name string) chan string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if channel, exists := p.streams[name]; exists {
		return channel
	}

	channel := make(chan string, p.bufferSize)
	p.streams[name] = channel

	return channel
}

func (p *Pipeline) closeAllStreams() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	closed := make(map[chan string]bool)
	for _, ch := range p.streams {
		if !closed[ch] {
			close(ch)

			closed[ch] = true
		}
	}
}

func (p *Pipeline) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	p.decorators = append(p.decorators, decoratorEntry{
		function: function,
		in:       input,
		out:      output,
	})

	p.getOrCreateStream(input)
	p.getOrCreateStream(output)
}

func (p *Pipeline) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	p.multiplexers = append(p.multiplexers, multiplexerEntry{
		function: function,
		ins:      inputs,
		out:      output,
	})

	for _, input := range inputs {
		p.getOrCreateStream(input)
	}

	p.getOrCreateStream(output)
}

func (p *Pipeline) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	p.separators = append(p.separators, separatorEntry{
		function: function,
		in:       input,
		outs:     outputs,
	})

	p.getOrCreateStream(input)

	for _, output := range outputs {
		p.getOrCreateStream(output)
	}
}

func (p *Pipeline) Run(ctx context.Context) error {
	defer p.closeAllStreams()

	workerGroup, workerGroupCtx := errgroup.WithContext(ctx)

	for _, decoratorItem := range p.decorators {
		workerGroup.Go(func() error {
			inputChan := p.getOrCreateStream(decoratorItem.in)
			outputChan := p.getOrCreateStream(decoratorItem.out)

			return decoratorItem.function(workerGroupCtx, inputChan, outputChan)
		})
	}

	for _, multiplexerItem := range p.multiplexers {
		workerGroup.Go(func() error {
			inputChannels := make([]chan string, len(multiplexerItem.ins))
			for index, input := range multiplexerItem.ins {
				inputChannels[index] = p.getOrCreateStream(input)
			}

			outputChan := p.getOrCreateStream(multiplexerItem.out)

			return multiplexerItem.function(workerGroupCtx, inputChannels, outputChan)
		})
	}

	for _, separatorItem := range p.separators {
		workerGroup.Go(func() error {
			inputChan := p.getOrCreateStream(separatorItem.in)
			outputChannels := make([]chan string, len(separatorItem.outs))
			for index, output := range separatorItem.outs {
				outputChannels[index] = p.getOrCreateStream(output)
			}

			return separatorItem.function(workerGroupCtx, inputChan, outputChannels)
		})
	}

	err := workerGroup.Wait()
	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

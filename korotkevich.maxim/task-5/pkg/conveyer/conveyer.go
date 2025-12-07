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

const undefined = "undefined"

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

type Conveyer struct {
	size         int
	channels     map[string]chan string
	mu           sync.RWMutex
	decorators   []decoratorEntry
	multiplexers []multiplexerEntry
	separators   []separatorEntry
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:         size,
		channels:     make(map[string]chan string),
		mu:           sync.RWMutex{},
		decorators:   make([]decoratorEntry, 0),
		multiplexers: make([]multiplexerEntry, 0),
		separators:   make([]separatorEntry, 0),
	}
}

func (conv *Conveyer) getOrCreateChan(name string) chan string {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	if chn, ok := conv.channels[name]; ok {
		return chn
	}

	chn := make(chan string, conv.size)
	conv.channels[name] = chn

	return chn
}

func (conv *Conveyer) getChan(name string) (chan string, bool) {
	conv.mu.RLock()
	defer conv.mu.RUnlock()

	chn, exists := conv.channels[name]

	return chn, exists
}

func (conv *Conveyer) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	conv.getOrCreateChan(input)
	conv.getOrCreateChan(output)

	conv.decorators = append(conv.decorators, decoratorEntry{
		function: function,
		input:    input,
		output:   output,
	})
}

func (conv *Conveyer) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, name := range inputs {
		conv.getOrCreateChan(name)
	}

	conv.getOrCreateChan(output)

	conv.multiplexers = append(conv.multiplexers, multiplexerEntry{
		function: function,
		inputs:   inputs,
		output:   output,
	})
}

func (conv *Conveyer) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	conv.getOrCreateChan(input)

	for _, name := range outputs {
		conv.getOrCreateChan(name)
	}

	conv.separators = append(conv.separators, separatorEntry{
		function: function,
		input:    input,
		outputs:  outputs,
	})
}

func (conv *Conveyer) closeAllChannels() {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	closed := make(map[chan string]bool)

	for _, chn := range conv.channels {
		if chn == nil {
			continue
		}

		if !closed[chn] {
			closed[chn] = true

			close(chn)
		}
	}
}

func (conv *Conveyer) Run(ctx context.Context) error {
	defer conv.closeAllChannels()

	group, gctx := errgroup.WithContext(ctx)

	startDecorator := func(decr decoratorEntry) {
		inputChan := conv.getOrCreateChan(decr.input)
		outputCh := conv.getOrCreateChan(decr.output)

		group.Go(func() error {
			return decr.function(gctx, inputChan, outputCh)
		})
	}

	startMultiplexer := func(mulp multiplexerEntry) {
		outputCh := conv.getOrCreateChan(mulp.output)
		inputChans := make([]chan string, len(mulp.inputs))

		for i, name := range mulp.inputs {
			inputChans[i] = conv.getOrCreateChan(name)
		}

		group.Go(func() error {
			return mulp.function(gctx, inputChans, outputCh)
		})
	}

	startSeparator := func(sepr separatorEntry) {
		inputChan := conv.getOrCreateChan(sepr.input)
		outputChans := make([]chan string, len(sepr.outputs))

		for i, name := range sepr.outputs {
			outputChans[i] = conv.getOrCreateChan(name)
		}

		group.Go(func() error {
			return sepr.function(gctx, inputChan, outputChans)
		})
	}

	for _, decr := range conv.decorators {
		startDecorator(decr)
	}

	for _, mulp := range conv.multiplexers {
		startMultiplexer(mulp)
	}

	for _, sepr := range conv.separators {
		startSeparator(sepr)
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

func (conv *Conveyer) Send(input string, data string) error {
	chn, exists := conv.getChan(input)
	if !exists {
		return ErrChanNotFound
	}

	select {
	case chn <- data:
		return nil
	default:
		return ErrChanFull
	}
}

func (conv *Conveyer) Recv(output string) (string, error) {
	chn, exists := conv.getChan(output)
	if !exists {
		return "", ErrChanNotFound
	}

	data, received := <-chn
	if !received {
		return undefined, nil
	}

	return data, nil
}

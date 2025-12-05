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
		decorators:   make([]decoratorEntry, 0),
		multiplexers: make([]multiplexerEntry, 0),
		separators:   make([]separatorEntry, 0),
	}
}

func (conv *Conveyer) getOrCreateChan(name string) chan string {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	if ch, ok := conv.channels[name]; ok {
		return ch
	}

	ch := make(chan string, conv.size)
	conv.channels[name] = ch
	return ch
}

func (conv *Conveyer) getChan(name string) (chan string, bool) {
	conv.mu.RLock()
	defer conv.mu.RUnlock()

	ch, ok := conv.channels[name]
	return ch, ok
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
	for _, ch := range conv.channels {
		if ch == nil {
			continue
		}
		if !closed[ch] {
			closed[ch] = true
			close(ch)
		}
	}
}

func (conv *Conveyer) Run(ctx context.Context) error {
	defer conv.closeAllChannels()

	group, gctx := errgroup.WithContext(ctx)

	startDecorator := func(d decoratorEntry) {
		inputChan := conv.getOrCreateChan(d.input)
		outputCh := conv.getOrCreateChan(d.output)

		group.Go(func() error {
			return d.function(gctx, inputChan, outputCh)
		})
	}

	startMultiplexer := func(m multiplexerEntry) {
		outputCh := conv.getOrCreateChan(m.output)
		inputChans := make([]chan string, len(m.inputs))
		for i, name := range m.inputs {
			inputChans[i] = conv.getOrCreateChan(name)
		}

		group.Go(func() error {
			return m.function(gctx, inputChans, outputCh)
		})
	}

	startSeparator := func(s separatorEntry) {
		inputChan := conv.getOrCreateChan(s.input)
		outputChans := make([]chan string, len(s.outputs))
		for i, name := range s.outputs {
			outputChans[i] = conv.getOrCreateChan(name)
		}

		group.Go(func() error {
			return s.function(gctx, inputChan, outputChans)
		})
	}

	for _, d := range conv.decorators {
		startDecorator(d)
	}

	for _, m := range conv.multiplexers {
		startMultiplexer(m)
	}

	for _, s := range conv.separators {
		startSeparator(s)
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

func (conv *Conveyer) Send(input string, data string) error {
	ch, ok := conv.getChan(input)
	if !ok {
		return ErrChanNotFound
	}

	select {
	case ch <- data:
		return nil
	default:
		return ErrChanFull
	}
}

func (conv *Conveyer) Recv(output string) (string, error) {
	ch, ok := conv.getChan(output)
	if !ok {
		return "", ErrChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return undefined, nil
	}

	return data, nil
}


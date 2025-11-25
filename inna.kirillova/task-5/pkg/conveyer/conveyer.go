package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	mu       sync.RWMutex

	decorators   []decoratorEntry
	multiplexers []multiplexerEntry
	separators   []separatorEntry
}

type decoratorEntry struct {
	fn     func(ctx context.Context, input chan string, output chan string) error
	input  string
	output string
}

type multiplexerEntry struct {
	fn     func(ctx context.Context, inputs []chan string, output chan string) error
	inputs []string
	output string
}

type separatorEntry struct {
	fn      func(ctx context.Context, input chan string, outputs []chan string) error
	input   string
	outputs []string
}

type conveyer interface {
	RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string)
	RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

func New(size int) conveyer {
	return &conveyerImpl{
		size:         size,
		channels:     make(map[string]chan string),
		decorators:   make([]decoratorEntry, 0),
		multiplexers: make([]multiplexerEntry, 0),
		separators:   make([]separatorEntry, 0),
	}
}

func (c *conveyerImpl) getOrCreateChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) getChan(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.channels[name]
	return ch, ok
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, ok := c.getChan(input)
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, ok := c.getChan(output)
	if !ok {
		return "", ErrChanNotFound
	}
	data, ok := <-ch
	if !ok {
		return undefined, nil
	}
	return data, nil
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *conveyerImpl) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string) {
	c.decorators = append(c.decorators, decoratorEntry{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *conveyerImpl) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	c.multiplexers = append(c.multiplexers, multiplexerEntry{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *conveyerImpl) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	c.separators = append(c.separators, separatorEntry{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	g, ctx := errgroup.WithContext(ctx)

	for _, d := range c.decorators {
		d := d
		in := c.getOrCreateChan(d.input)
		out := c.getOrCreateChan(d.output)
		g.Go(func() error {
			return d.fn(ctx, in, out)
		})
	}

	for _, m := range c.multiplexers {
		m := m
		out := c.getOrCreateChan(m.output)
		ins := make([]chan string, len(m.inputs))
		for i, name := range m.inputs {
			ins[i] = c.getOrCreateChan(name)
		}
		g.Go(func() error {
			return m.fn(ctx, ins, out)
		})
	}

	for _, s := range c.separators {
		s := s
		in := c.getOrCreateChan(s.input)
		outs := make([]chan string, len(s.outputs))
		for i, name := range s.outputs {
			outs[i] = c.getOrCreateChan(name)
		}
		g.Go(func() error {
			return s.fn(ctx, in, outs)
		})
	}

	return g.Wait()
}

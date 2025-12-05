package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrChanFull     = errors.New("chan is full")
)

const undefined = "undefined"

type decoratorEntry struct {
	fn     func(context.Context, chan string, chan string) error
	input  string
	output string
}

type multiplexerEntry struct {
	fn     func(context.Context, []chan string, chan string) error
	inputs []string
	output string
}

type separatorEntry struct {
	fn      func(context.Context, chan string, []chan string) error
	input   string
	outputs []string
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
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *Conveyer) getOrCreate(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) get(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[name]
	return ch, ok
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	c.getOrCreate(input)
	c.getOrCreate(output)

	c.decorators = append(c.decorators, decoratorEntry{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	for _, v := range inputs {
		c.getOrCreate(v)
	}
	c.getOrCreate(output)

	c.multiplexers = append(c.multiplexers, multiplexerEntry{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	c.getOrCreate(input)

	for _, v := range outputs {
		c.getOrCreate(v)
	}

	c.separators = append(c.separators, separatorEntry{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *Conveyer) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, gctx := errgroup.WithContext(ctx)

	for _, d := range c.decorators {
		input := c.getOrCreate(d.input)
		output := c.getOrCreate(d.output)
		fn := d.fn

		group.Go(func() error {
			return fn(gctx, input, output)
		})
	}

	for _, m := range c.multiplexers {
		var inputs []chan string

		for _, name := range m.inputs {
			inputs = append(inputs, c.getOrCreate(name))
		}

		output := c.getOrCreate(m.output)
		fn := m.fn

		group.Go(func() error {
			return fn(gctx, inputs, output)
		})
	}

	for _, s := range c.separators {
		input := c.getOrCreate(s.input)
		var outputs []chan string

		for _, name := range s.outputs {
			outputs = append(outputs, c.getOrCreate(name))
		}

		fn := s.fn

		group.Go(func() error {
			return fn(gctx, input, outputs)
		})
	}

	defer c.closeAll()

	return group.Wait()
}

func (c *Conveyer) Send(input string, data string) error {
	ch, ok := c.get(input)
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

func (c *Conveyer) Recv(output string) (string, error) {
	ch, ok := c.get(output)
	if !ok {
		return "", ErrChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return undefined, nil
	}

	return data, nil
}

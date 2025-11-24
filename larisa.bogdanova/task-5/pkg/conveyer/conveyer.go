package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	undefined       = "undefined"
)

type Conveyer interface {
	RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input, output string)
	RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyer struct {
	size  int
	mu    sync.RWMutex
	chans map[string]chan string

	decorators   []decoratorReg
	multiplexers []multiplexerReg
	separators   []separatorReg
}

type decoratorReg struct {
	fn     func(ctx context.Context, input chan string, output chan string) error
	input  string
	output string
}

type multiplexerReg struct {
	fn     func(ctx context.Context, inputs []chan string, output chan string) error
	inputs []string
	output string
}

type separatorReg struct {
	fn      func(ctx context.Context, input chan string, outputs []chan string) error
	input   string
	outputs []string
}

func New(size int) Conveyer {
	return &conveyer{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (c *conveyer) getChan(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.chans[name]
	return ch, ok
}

func (c *conveyer) getOrCreateChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.chans[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.chans[name] = ch
	return ch
}

func (c *conveyer) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input, output string) {
	c.getOrCreateChan(input)
	c.getOrCreateChan(output)
	c.decorators = append(c.decorators, decoratorReg{fn: fn, input: input, output: output})
}

func (c *conveyer) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	for _, name := range inputs {
		c.getOrCreateChan(name)
	}
	c.getOrCreateChan(output)
	c.multiplexers = append(c.multiplexers, multiplexerReg{fn: fn, inputs: inputs, output: output})
}

func (c *conveyer) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	c.getOrCreateChan(input)
	for _, name := range outputs {
		c.getOrCreateChan(name)
	}
	c.separators = append(c.separators, separatorReg{fn: fn, input: input, outputs: outputs})
}

func (c *conveyer) Send(input string, data string) error {
	ch, ok := c.getChan(input)
	if !ok {
		return ErrChanNotFound
	}
	select {
	case ch <- data:
		return nil
	default:
		return errors.New("channel is full")
	}
}

func (c *conveyer) Recv(output string) (string, error) {
	ch, ok := c.getChan(output)
	if !ok {
		return "", ErrChanNotFound
	}
	val, ok := <-ch
	if !ok {
		return undefined, nil
	}
	return val, nil
}

func (c *conveyer) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, reg := range c.decorators {
		in := c.getOrCreateChan(reg.input)
		out := c.getOrCreateChan(reg.output)
		g.Go(func() error {
			return reg.fn(ctx, in, out)
		})
	}

	for _, reg := range c.multiplexers {
		var inputs []chan string
		for _, name := range reg.inputs {
			inputs = append(inputs, c.getOrCreateChan(name))
		}
		output := c.getOrCreateChan(reg.output)
		g.Go(func() error {
			return reg.fn(ctx, inputs, output)
		})
	}

	for _, reg := range c.separators {
		input := c.getOrCreateChan(reg.input)
		var outputs []chan string
		for _, name := range reg.outputs {
			outputs = append(outputs, c.getOrCreateChan(name))
		}
		g.Go(func() error {
			return reg.fn(ctx, input, outputs)
		})
	}

	err := g.Wait()

	c.mu.Lock()
	for _, ch := range c.chans {
		close(ch)
	}
	c.chans = make(map[string]chan string)
	c.mu.Unlock()

	return err
}

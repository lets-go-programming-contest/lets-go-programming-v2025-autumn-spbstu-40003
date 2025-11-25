package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Conveyer interface {
	RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string)
	RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type channelMap struct {
	sync.RWMutex
	m map[string]chan string
}

type ConveyerImpl struct {
	channels    *channelMap
	tasks       []func(context.Context) error
	channelSize int
}

func New(size int) Conveyer {
	return &ConveyerImpl{
		channels: &channelMap{
			m: make(map[string]chan string),
		},
		channelSize: size,
		tasks:       make([]func(context.Context) error, 0),
	}
}

func (c *ConveyerImpl) getOrCreateChannel(name string) chan string {
	c.channels.Lock()
	defer c.channels.Unlock()

	if ch, exists := c.channels.m[name]; exists {
		return ch
	}

	ch := make(chan string, c.channelSize)
	c.channels.m[name] = ch
	return ch
}

func (c *ConveyerImpl) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, inputName string, outputName string) {
	inCh := c.getOrCreateChannel(inputName)
	outCh := c.getOrCreateChannel(outputName)

	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}
	c.tasks = append(c.tasks, task)
}

func (c *ConveyerImpl) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputNames []string, outputName string) {
	var inChs []chan string
	for _, name := range inputNames {
		inChs = append(inChs, c.getOrCreateChannel(name))
	}
	outCh := c.getOrCreateChannel(outputName)

	task := func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	}
	c.tasks = append(c.tasks, task)
}

func (c *ConveyerImpl) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, inputName string, outputNames []string) {
	inCh := c.getOrCreateChannel(inputName)
	var outChs []chan string
	for _, name := range outputNames {
		outChs = append(outChs, c.getOrCreateChannel(name))
	}

	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	}
	c.tasks = append(c.tasks, task)
}

func (c *ConveyerImpl) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	for _, task := range c.tasks {
		taskFunc := task
		eg.Go(func() error {
			return taskFunc(ctx)
		})
	}

	err := eg.Wait()

	c.closeAllChannels()

	return err
}

func (c *ConveyerImpl) closeAllChannels() {
	c.channels.Lock()
	defer c.channels.Unlock()
	for _, ch := range c.channels.m {
		func() {
			defer func() {
				_ = recover()
			}()
			close(ch)
		}()
	}
}

func (c *ConveyerImpl) Send(inputName string, data string) error {
	c.channels.RLock()
	ch, exists := c.channels.m[inputName]
	c.channels.RUnlock()

	if !exists {
		return errors.New("chan not found")
	}

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	ch <- data
	return nil
}

func (c *ConveyerImpl) Recv(outputName string) (string, error) {
	c.channels.RLock()
	ch, exists := c.channels.m[outputName]
	c.channels.RUnlock()

	if !exists {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return val, nil
}

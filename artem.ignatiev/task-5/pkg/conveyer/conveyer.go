package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const (
	resUndefined = "undefined"
)

type TaskFunc func(ctx context.Context) error

type Conveyer struct {
	mu          sync.Mutex
	channels    map[string]chan string
	chanSize    int
	tasks       []TaskFunc
	channelsKey []string
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:          sync.Mutex{},
		channels:    make(map[string]chan string),
		chanSize:    size,
		tasks:       make([]TaskFunc, 0),
		channelsKey: make([]string, 0),
	}
}

func (c *Conveyer) getOrMakeChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, ok := c.channels[name]; ok {
		return channel
	}

	newChannel := make(chan string, c.chanSize)
	c.channels[name] = newChannel
	c.channelsKey = append(c.channelsKey, name)

	return newChannel
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	inputChannel := c.getOrMakeChan(inputName)
	outputChannel := c.getOrMakeChan(outputName)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputsNames []string,
	outputName string,
) {
	inputChannels := make([]chan string, 0, len(inputsNames))

	for _, name := range inputsNames {
		inputChannels = append(inputChannels, c.getOrMakeChan(name))
	}

	outputChannel := c.getOrMakeChan(outputName)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputsNames []string,
) {
	inputChannel := c.getOrMakeChan(inputName)
	outputChannels := make([]chan string, 0, len(outputsNames))

	for _, name := range outputsNames {
		outputChannels = append(outputChannels, c.getOrMakeChan(name))
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, task := range c.tasks {
		currentTask := task

		group.Go(func() error {
			return currentTask(groupCtx)
		})
	}

	err := group.Wait()

	c.mu.Lock()
	for _, name := range c.channelsKey {
		close(c.channels[name])
	}
	c.mu.Unlock()

	if err != nil {
		return fmt.Errorf("conveyer run error: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(inputName string, data string) error {
	c.mu.Lock()
	channel, ok := c.channels[inputName]
	c.mu.Unlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(outputName string) (string, error) {
	c.mu.Lock()
	channel, ok := c.channels[outputName]
	c.mu.Unlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, opened := <-channel
	if !opened {
		return resUndefined, nil
	}

	return val, nil
}

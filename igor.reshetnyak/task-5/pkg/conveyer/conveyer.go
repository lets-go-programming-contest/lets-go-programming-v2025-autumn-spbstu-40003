package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const Undefined = "undefined"

var (
	ErrChanNotFound    = errors.New("chan not found")
	ErrAlreadyRunning  = errors.New("already running")
	ErrChannelFull     = errors.New("channel is full")
	ErrNoDataAvailable = errors.New("no data")
)

type Conveyer interface {
	RegisterDecorator(
		handlerFunc func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyerImpl struct {
	mu       sync.RWMutex
	channels map[string]chan string
	tasks    []func(context.Context) error
	size     int
	running  bool
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		mu:       sync.RWMutex{},
		channels: make(map[string]chan string),
		tasks:    make([]func(context.Context) error, 0),
		size:     size,
		running:  false,
	}
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *conveyerImpl) getChannel(name string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channel, ok := c.channels[name]
	if !ok {
		return nil, ErrChanNotFound
	}

	return channel, nil
}

func (c *conveyerImpl) RegisterDecorator(
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChannel := c.getOrCreateChannel(input)
	outputChannel := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannel)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *conveyerImpl) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))
	for index, name := range inputs {
		inputChannels[index] = c.getOrCreateChannel(name)
	}

	outputChannel := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannels, outputChannel)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *conveyerImpl) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChannel := c.getOrCreateChannel(input)

	outputChannels := make([]chan string, len(outputs))
	for index, name := range outputs {
		outputChannels[index] = c.getOrCreateChannel(name)
	}

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannels)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return ErrAlreadyRunning
	}

	c.running = true
	c.mu.Unlock()

	processingGroup, groupCtx := errgroup.WithContext(ctx)

	c.mu.RLock()
	tasks := make([]func(context.Context) error, len(c.tasks))
	copy(tasks, c.tasks)
	c.mu.RUnlock()

	for _, task := range tasks {
		currentTask := task

		processingGroup.Go(func() error {
			return currentTask(groupCtx)
		})
	}

	err := processingGroup.Wait()

	c.mu.Lock()
	for _, channel := range c.channels {
		close(channel)
	}

	c.running = false
	c.mu.Unlock()

	if err != nil {
		return fmt.Errorf("conveyer execution failed: %w", err)
	}

	return nil
}

func (c *conveyerImpl) Send(input string, data string) error {
	targetChannel, err := c.getChannel(input)
	if err != nil {
		return err
	}

	select {
	case targetChannel <- data:
		return nil

	default:
		return ErrChannelFull
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	targetChannel, err := c.getChannel(output)
	if err != nil {
		return "", err
	}

	select {
	case data, channelOpen := <-targetChannel:
		if !channelOpen {
			return Undefined, nil
		}

		return data, nil

	default:
		return "", ErrNoDataAvailable
	}
}

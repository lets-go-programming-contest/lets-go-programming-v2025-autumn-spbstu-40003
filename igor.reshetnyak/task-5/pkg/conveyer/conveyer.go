package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

const Undefined = "undefined"

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrSendFailed   = errors.New("send failed")
)

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
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
	cancel   context.CancelFunc
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		channels: make(map[string]chan string),
		tasks:    make([]func(context.Context) error, 0),
		size:     size,
	}
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) getChannel(name string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[name]
	if !ok {
		return nil, ErrChanNotFound
	}

	return ch, nil
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChan)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputChans[i] = c.getOrCreateChannel(name)
	}

	outputChan := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, inputChans, outputChan)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := c.getOrCreateChannel(input)

	outputChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outputChans[i] = c.getOrCreateChannel(name)
	}

	task := func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChans)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return errors.New("already running")
	}

	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.running = true
	c.mu.Unlock()

	defer func() {
		c.closeChannels()
		c.mu.Lock()
		c.running = false
		c.cancel = nil
		c.mu.Unlock()
	}()

	g, _ := errgroup.WithContext(runCtx)

	c.mu.RLock()
	tasks := make([]func(context.Context) error, len(c.tasks))
	copy(tasks, c.tasks)
	c.mu.RUnlock()

	for _, task := range tasks {
		currentTask := task
		g.Go(func() error {
			return currentTask(runCtx)
		})
	}

	select {
	case <-runCtx.Done():
		c.mu.Lock()
		if c.cancel != nil {
			c.cancel()
		}
		c.mu.Unlock()
		return runCtx.Err()
	default:
		return g.Wait()
	}
}

func (c *conveyerImpl) closeChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, err := c.getChannel(input)
	if err != nil {
		return err
	}

	c.mu.RLock()
	running := c.running
	c.mu.RUnlock()

	if !running {
		return ErrSendFailed
	}

	select {
	case ch <- data:
		return nil
	default:
		return ErrSendFailed
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, err := c.getChannel(output)
	if err != nil {
		return "", err
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return Undefined, nil
		}
		return data, nil
	default:
		return "", errors.New("no data")
	}
}

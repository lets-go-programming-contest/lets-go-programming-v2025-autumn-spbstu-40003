package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrSendTimeout  = errors.New("send timeout")
)

const defaultSendTimeout = 5 * time.Second

type Conveyer interface {
	RegisterDecorator(
		decoratorFunction func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		multiplexerFunction func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		separatorFunction func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type channelMap map[string]chan string

type handler func(ctx context.Context) error

type conveyerImpl struct {
	channelSize int
	channels    channelMap
	handlers    []handler
	mu          sync.RWMutex
	closeOnce   sync.Once
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		channelSize: size,
		channels:    make(channelMap),
		handlers:    make([]handler, 0),
		mu:          sync.RWMutex{},
		closeOnce:   sync.Once{},
	}
}

func (c *conveyerImpl) getOrCreateChannel(channelID string) chan string {
	if channel, ok := c.channels[channelID]; ok {
		return channel
	}

	channel := make(chan string, c.channelSize)
	c.channels[channelID] = channel

	return channel
}

func (c *conveyerImpl) RegisterDecorator(
	decoratorFunction func(ctx context.Context, input chan string, output chan string) error,
	inputID string,
	outputID string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(inputID)
	outputCh := c.getOrCreateChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return decoratorFunction(ctx, inputCh, outputCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	multiplexerFunction func(ctx context.Context, inputs []chan string, output chan string) error,
	inputIDs []string,
	outputID string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChs := make([]chan string, len(inputIDs))

	for i, id := range inputIDs {
		inputChs[i] = c.getOrCreateChannel(id)
	}

	outputCh := c.getOrCreateChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return multiplexerFunction(ctx, inputChs, outputCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	separatorFunction func(ctx context.Context, input chan string, outputs []chan string) error,
	inputID string,
	outputIDs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(inputID)
	outputChs := make([]chan string, len(outputIDs))

	for i, id := range outputIDs {
		outputChs[i] = c.getOrCreateChannel(id)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return separatorFunction(ctx, inputCh, outputChs)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.RLock()

	group, runCtx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		handler := h

		group.Go(func() error {
			return handler(runCtx)
		})
	}

	c.mu.RUnlock()

	err := group.Wait()

	c.closeAllChannels()

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil
		}

		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

func (c *conveyerImpl) closeAllChannels() {
	c.closeOnce.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		for _, channel := range c.channels {
			if channel != nil {
				close(channel)
			}
		}
	})
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channel, ok := c.channels[input]
	if !ok {
		return fmt.Errorf("%w", ErrChanNotFound)
	}

	select {
	case channel <- data:
		return nil
	case <-time.After(defaultSendTimeout):
		return fmt.Errorf("%w", ErrSendTimeout)
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channel, ok := c.channels[output]
	if !ok {
		return "", fmt.Errorf("%w", ErrChanNotFound)
	}

	data, open := <-channel
	if !open {
		return "undefined", nil
	}

	return data, nil
}

package conveyer

import (
	"context"
	"errors"
	"sync"
	"time"

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

type channelMap map[string]chan string

type handler func(ctx context.Context) error

type conveyerImpl struct {
	channelSize int
	channels    channelMap
	handlers    []handler
	mu          sync.RWMutex

	closeOnce sync.Once
}

func New(size int) Conveyer {
	return &conveyerImpl{
		channelSize: size,
		channels:    make(channelMap),
		handlers:    make([]handler, 0),
	}
}

func (c *conveyerImpl) getOrCreateChannel(id string) chan string {
	if ch, ok := c.channels[id]; ok {
		return ch
	}
	ch := make(chan string, c.channelSize)
	c.channels[id] = ch
	return ch
}

func (c *conveyerImpl) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, inputID string, outputID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(inputID)
	outputCh := c.getOrCreateChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputIDs []string, outputID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChs := make([]chan string, len(inputIDs))
	for i, id := range inputIDs {
		inputChs[i] = c.getOrCreateChannel(id)
	}
	outputCh := c.getOrCreateChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChs, outputCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, inputID string, outputIDs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(inputID)
	outputChs := make([]chan string, len(outputIDs))
	for i, id := range outputIDs {
		outputChs[i] = c.getOrCreateChannel(id)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputChs)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	group, runCtx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		handler := h
		group.Go(func() error {
			return handler(runCtx)
		})
	}

	err := group.Wait()

	c.closeAllChannels()

	if err != nil && errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func (c *conveyerImpl) closeAllChannels() {
	c.closeOnce.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		for _, ch := range c.channels {
			if ch != nil {
				close(ch)
			}
		}
	})
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[input]
	if !ok {
		return errors.New("chan not found")
	}

	select {
	case ch <- data:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("send timeout")
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[output]
	if !ok {
		return "", errors.New("chan not found")
	}

	data, open := <-ch
	if !open {
		return "undefined", nil
	}
	return data, nil
}

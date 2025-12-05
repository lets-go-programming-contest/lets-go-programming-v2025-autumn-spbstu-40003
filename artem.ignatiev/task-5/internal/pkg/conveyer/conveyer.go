package conveyer

import (
	"context"
	"errors"
	"sync"
)

const (
	ErrMsgChanNotFound = "chan not found"
	ResUndefined       = "undefined"
)

type taskFunc func(ctx context.Context) error

type Conveyer struct {
	mu          sync.Mutex
	channels    map[string]chan string
	chanSize    int
	tasks       []taskFunc
	channelsKey []string
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels:    make(map[string]chan string),
		chanSize:    size,
		tasks:       make([]taskFunc, 0),
		channelsKey: make([]string, 0),
	}
}

func (c *Conveyer) getOrMakeChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	newCh := make(chan string, c.chanSize)
	c.channels[name] = newCh
	c.channelsKey = append(c.channelsKey, name)
	return newCh
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	in := c.getOrMakeChan(input)
	out := c.getOrMakeChan(output)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	var inChs []chan string
	for _, name := range inputs {
		inChs = append(inChs, c.getOrMakeChan(name))
	}
	out := c.getOrMakeChan(output)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inChs, out)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	in := c.getOrMakeChan(input)
	var outChs []chan string
	for _, name := range outputs {
		outChs = append(outChs, c.getOrMakeChan(name))
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, in, outChs)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, len(c.tasks))

	for _, task := range c.tasks {
		wg.Add(1)
		go func(t taskFunc) {
			defer wg.Done()
			if err := t(ctx); err != nil {
				select {
				case errCh <- err:
					cancel()
				default:
				}
			}
		}(task)
	}

	wg.Wait()

	c.mu.Lock()
	for _, name := range c.channelsKey {
		close(c.channels[name])
	}
	c.mu.Unlock()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	ch, ok := c.channels[input]
	c.mu.Unlock()

	if !ok {
		return errors.New(ErrMsgChanNotFound)
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	ch, ok := c.channels[output]
	c.mu.Unlock()

	if !ok {
		return "", errors.New(ErrMsgChanNotFound)
	}

	val, opened := <-ch
	if !opened {
		return ResUndefined, nil
	}
	return val, nil
}

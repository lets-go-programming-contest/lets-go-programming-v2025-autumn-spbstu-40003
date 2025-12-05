package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"

	"ivantsov.egor/task-5/pkg/handlers"
)

type DecoratorFunc func(context.Context, chan string, chan string) error

type pipeline struct {
	bufferSize int
	mu         sync.Mutex

	decorators []DecoratorFunc

	inCh  chan string
	outCh chan string
}

func New(bufferSize int) *pipeline {
	return &pipeline{
		bufferSize: bufferSize,
		mu:         sync.Mutex{},
		decorators: make([]DecoratorFunc, 0),
	}
}

func (p *pipeline) RegisterDecorator(fn DecoratorFunc) error {
	if fn == nil {
		return errors.New("nil decorator")
	}

	p.mu.Lock()
	p.decorators = append(p.decorators, fn)
	p.mu.Unlock()

	return nil
}

func (p *pipeline) init() {
	if p.inCh != nil {
		return
	}

	p.inCh = make(chan string, p.bufferSize)

	prev := p.inCh

	for _, d := range p.decorators {
		next := make(chan string, p.bufferSize)

		go func(input, output chan string, dec DecoratorFunc) {
			_ = dec(context.Background(), input, output)
			close(output)
		}(prev, next, d)

		prev = next
	}

	p.outCh = prev
}

func (p *pipeline) Send(ctx context.Context, v string) error {
	p.mu.Lock()
	p.init()
	in := p.inCh
	p.mu.Unlock()

	select {
	case in <- v:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *pipeline) Recv(ctx context.Context) (string, error) {
	p.mu.Lock()
	p.init()
	out := p.outCh
	p.mu.Unlock()

	select {
	case v, ok := <-out:
		if !ok {
			return "", errors.New("pipeline closed")
		}
		return v, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (p *pipeline) Run(ctx context.Context) error {
	p.mu.Lock()
	p.init()
	in, out := p.inCh, p.outCh
	p.mu.Unlock()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		defer close(in)
		<-ctx.Done()
		return nil
	})

	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _, ok := <-out:
				if !ok {
					return nil
				}
			}
		}
	})

	return group.Wait()
}

func RegisterDefaultHandlers(p *pipeline) error {
	return p.RegisterDecorator(handlers.PrefixDecoratorFunc)
}

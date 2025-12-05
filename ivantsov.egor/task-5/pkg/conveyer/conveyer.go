package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"ivantsov.egor/task-5/pkg/handlers"
)

type DecoratorFunc = func(context.Context, chan string, chan string) error

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
		inCh:       make(chan string),
		outCh:      make(chan string),
	}
}

func (p *pipeline) RegisterDecorator(fn DecoratorFunc, _, _ string) error {
	if fn == nil {
		return errors.New("nil decorator")
	}

	p.mu.Lock()
	p.decorators = append(p.decorators, fn)
	p.mu.Unlock()

	return nil
}

func (p *pipeline) Send(value string) error {
	select {
	case p.inCh <- value:
		return nil
	default:
		return errors.New("pipeline closed")
	}
}

func (p *pipeline) Recv(_ string) (string, error) {
	value, ok := <-p.outCh
	if !ok {
		return "", errors.New("pipeline closed")
	}

	return value, nil
}

func (p *pipeline) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	in := p.inCh

	for _, decorator := range p.decorators {
		out := make(chan string, p.bufferSize)

		currentDecorator := decorator

		group.Go(func() error {
			if err := currentDecorator(ctx, in, out); err != nil {
				return fmt.Errorf("decorator error: %w", err)
			}

			close(out)

			return nil
		})

		in = out
	}

	group.Go(func() error {
		return handlers.MultiplexerFunc(ctx, []chan string{in}, p.outCh)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("pipeline error: %w", err)
	}

	close(p.outCh)

	return nil
}

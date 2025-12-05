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

	multiplexers []multiplexer

	inputs  map[string]chan string
	outputs map[string]chan string
}

type multiplexer struct {
	inputNames []string
	outputName string
}

func New(bufferSize int) *pipeline {
	return &pipeline{
		bufferSize: bufferSize,
		inputs:     map[string]chan string{},
		outputs:    map[string]chan string{},
	}
}

func (p *pipeline) getInput(name string) chan string {
	ch, ok := p.inputs[name]
	if !ok {
		ch = make(chan string, p.bufferSize)
		p.inputs[name] = ch
	}
	return ch
}

func (p *pipeline) getOutput(name string) chan string {
	ch, ok := p.outputs[name]
	if !ok {
		ch = make(chan string, p.bufferSize)
		p.outputs[name] = ch
	}
	return ch
}

func (p *pipeline) RegisterDecorator(fn DecoratorFunc, inputName string, outputName string) error {
	if fn == nil {
		return errors.New("nil decorator")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	in := p.getInput(inputName)
	out := p.getOutput(outputName)

	p.decorators = append(p.decorators, func(ctx context.Context, _, _ chan string) error {
		return fn(ctx, in, out)
	})

	return nil
}

func (p *pipeline) RegisterMultiplexer(inputNames []string, outputName string) error {
	if len(inputNames) == 0 {
		return errors.New("no inputs")
	}

	p.mu.Lock()
	p.multiplexers = append(p.multiplexers, multiplexer{
		inputNames: inputNames,
		outputName: outputName,
	})
	p.mu.Unlock()

	return nil
}

func (p *pipeline) Send(inputName string, value string) error {
	ch, ok := p.inputs[inputName]
	if !ok {
		return errors.New("input not found")
	}

	select {
	case ch <- value:
		return nil
	default:
		return errors.New("pipeline closed")
	}
}

func (p *pipeline) Recv(outputName string) (string, error) {
	ch, ok := p.outputs[outputName]
	if !ok {
		return "", errors.New("output not found")
	}

	val, ok := <-ch
	if !ok {
		return "", errors.New("output closed")
	}

	return val, nil
}

func (p *pipeline) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, d := range p.decorators {
		decorator := d

		group.Go(func() error {
			if err := decorator(ctx, nil, nil); err != nil {
				return fmt.Errorf("decorator error: %w", err)
			}
			return nil
		})
	}

	for _, m := range p.multiplexers {
		mp := m

		group.Go(func() error {
			var inputs []chan string

			for _, name := range mp.inputNames {
				inputs = append(inputs, p.inputs[name])
			}

			output := p.getOutput(mp.outputName)

			if err := handlers.MultiplexerFunc(ctx, inputs, output); err != nil {
				return fmt.Errorf("multiplexer error: %w", err)
			}

			close(output)

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

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

var (
	ErrNilDecorator   = errors.New("nil decorator")
	ErrNoInputs       = errors.New("no inputs")
	ErrInputNotFound  = errors.New("input not found")
	ErrOutputNotFound = errors.New("output not found")
	ErrOutputClosed   = errors.New("output closed")
	ErrPipelineClosed = errors.New("pipeline closed")
)

type pipeline struct {
	bufferSize int
	mu         sync.Mutex

	decorators   []DecoratorFunc
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
		bufferSize:   bufferSize,
		mu:           sync.Mutex{},
		decorators:   nil,
		multiplexers: nil,
		inputs:       map[string]chan string{},
		outputs:      map[string]chan string{},
	}
}

func (p *pipeline) getInput(name string) chan string {
	channel, ok := p.inputs[name]
	if !ok {
		channel = make(chan string, p.bufferSize)
		p.inputs[name] = channel
	}
	return channel
}

func (p *pipeline) getOutput(name string) chan string {
	channel, ok := p.outputs[name]
	if !ok {
		channel = make(chan string, p.bufferSize)
		p.outputs[name] = channel
	}
	return channel
}

func (p *pipeline) RegisterDecorator(
	decoratorFunc DecoratorFunc,
	inputName string,
	outputName string,
) error {

	if decoratorFunc == nil {
		return ErrNilDecorator
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	inputChannel := p.getInput(inputName)
	outputChannel := p.getOutput(outputName)

	p.decorators = append(
		p.decorators,
		func(ctx context.Context, _, _ chan string) error {
			return decoratorFunc(ctx, inputChannel, outputChannel)
		},
	)

	return nil
}

func (p *pipeline) RegisterMultiplexer(
	inputNames []string,
	outputName string,
) error {

	if len(inputNames) == 0 {
		return ErrNoInputs
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.multiplexers = append(p.multiplexers, multiplexer{
		inputNames: inputNames,
		outputName: outputName,
	})

	return nil
}

func (p *pipeline) Send(inputName string, value string) error {
	channel, ok := p.inputs[inputName]
	if !ok {
		return ErrInputNotFound
	}

	select {
	case channel <- value:
		return nil
	default:
		return ErrPipelineClosed
	}
}

func (p *pipeline) Recv(outputName string) (string, error) {
	channel, ok := p.outputs[outputName]
	if !ok {
		return "", ErrOutputNotFound
	}

	value, ok := <-channel
	if !ok {
		return "", ErrOutputClosed
	}

	return value, nil
}

func (p *pipeline) Run(ctx context.Context) error {
	group, groupContext := errgroup.WithContext(ctx)

	for _, decorator := range p.decorators {
		decoratorFunc := decorator

		group.Go(func() error {
			if err := decoratorFunc(groupContext, nil, nil); err != nil {
				return fmt.Errorf("decorator error: %w", err)
			}
			return nil
		})
	}

	for _, multiplexerCfg := range p.multiplexers {
		cfg := multiplexerCfg

		group.Go(func() error {
			var inputChannels []chan string

			for _, name := range cfg.inputNames {
				inputChannels = append(inputChannels, p.inputs[name])
			}

			outputChannel := p.getOutput(cfg.outputName)

			if err := handlers.MultiplexerFunc(groupContext, inputChannels, outputChannel); err != nil {
				return fmt.Errorf("multiplexer error: %w", err)
			}

			close(outputChannel)

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("run error: %w", err)
	}

	return nil
}

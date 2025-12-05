package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type DecoratorFunc = func(context.Context, chan string, chan string) error
type SeparatorFunc = func(context.Context, chan string, []chan string) error
type MultiplexerFunc = func(context.Context, []chan string, chan string) error

var (
	ErrNilDecorator   = errors.New("nil decorator")
	ErrNilSeparator   = errors.New("nil separator")
	ErrNilMultiplexer = errors.New("nil multiplexer")
	ErrPipelineClosed = errors.New("pipeline closed")
	ErrInputNotFound  = errors.New("input not found")
	ErrOutputNotFound = errors.New("output not found")
)

type pipeline struct {
	bufferSize int

	mu sync.Mutex

	decorators   []DecoratorFunc
	separators   []separator
	multiplexers []multiplexer

	inputs  map[string]chan string
	outputs map[string]chan string
}

type separator struct {
	fn          SeparatorFunc
	inputName   string
	outputNames []string
}

type multiplexer struct {
	fn         MultiplexerFunc
	inputNames []string
	outputName string
}

func New(bufferSize int) *pipeline {
	return &pipeline{
		bufferSize:   bufferSize,
		decorators:   []DecoratorFunc{},
		separators:   []separator{},
		multiplexers: []multiplexer{},
		inputs:       map[string]chan string{},
		outputs:      map[string]chan string{},
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

func (p *pipeline) RegisterDecorator(
	fn DecoratorFunc,
	inputName string,
	outputName string,
) error {
	if fn == nil {
		return ErrNilDecorator
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	input := p.getInput(inputName)
	output := p.getOutput(outputName)

	p.decorators = append(p.decorators,
		func(ctx context.Context, _, _ chan string) error {
			return fn(ctx, input, output)
		},
	)

	return nil
}

func (p *pipeline) RegisterSeparator(
	fn SeparatorFunc,
	inputName string,
	outputNames []string,
) error {
	if fn == nil {
		return ErrNilSeparator
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.separators = append(p.separators, separator{
		fn:          fn,
		inputName:   inputName,
		outputNames: append([]string(nil), outputNames...),
	})

	return nil
}

func (p *pipeline) RegisterMultiplexer(
	fn MultiplexerFunc,
	inputNames []string,
	outputName string,
) error {
	if fn == nil {
		return ErrNilMultiplexer
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.multiplexers = append(p.multiplexers, multiplexer{
		fn:         fn,
		inputNames: append([]string(nil), inputNames...),
		outputName: outputName,
	})

	return nil
}

func (p *pipeline) Send(inputName string, value string) error {
	input, ok := p.inputs[inputName]
	if !ok {
		return ErrInputNotFound
	}

	select {
	case input <- value:
		return nil

	default:
		return ErrPipelineClosed
	}
}

func (p *pipeline) Recv(outputName string) (string, error) {
	output, ok := p.outputs[outputName]
	if !ok {
		return "", ErrOutputNotFound
	}

	value, ok := <-output
	if !ok {
		return "", ErrPipelineClosed
	}

	return value, nil
}

func (p *pipeline) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, decorator := range p.decorators {
		decoratorFunc := decorator

		group.Go(func() error {
			if err := decoratorFunc(groupCtx, nil, nil); err != nil {
				return fmt.Errorf("decorator: %w", err)
			}
			return nil
		})
	}

	for _, sep := range p.separators {
		separatorData := sep

		group.Go(func() error {
			input := p.getInput(separatorData.inputName)

			outputs := make([]chan string, 0, len(separatorData.outputNames))
			for _, name := range separatorData.outputNames {
				outputs = append(outputs, p.getOutput(name))
			}

			if err := separatorData.fn(groupCtx, input, outputs); err != nil {
				return fmt.Errorf("separator: %w", err)
			}

			for _, ch := range outputs {
				close(ch)
			}

			return nil
		})
	}

	for _, mux := range p.multiplexers {
		multiplexerData := mux

		group.Go(func() error {
			inputs := make([]chan string, 0, len(multiplexerData.inputNames))
			for _, name := range multiplexerData.inputNames {
				inputs = append(inputs, p.getInput(name))
			}

			output := p.getOutput(multiplexerData.outputName)

			if err := multiplexerData.fn(groupCtx, inputs, output); err != nil {
				return fmt.Errorf("multiplexer: %w", err)
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

package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type DecoratorFunc = func(context.Context, chan string, chan string) error
type MultiplexerFunc = func(context.Context, []chan string, chan string) error
type SeparatorFunc = func(context.Context, chan string, []chan string) error

var (
	ErrNilFunc        = errors.New("nil function")
	ErrPipelineClosed = errors.New("pipeline closed")
)

type pipeline struct {
	bufferSize int
	lock       sync.Mutex

	decorators   []DecoratorFunc
	multiplexers []multiplexer
	separators   []separator

	inputs  map[string]chan string
	outputs map[string]chan string
}

type multiplexer struct {
	fn         MultiplexerFunc
	inputNames []string
	outputName string
}

type separator struct {
	fn        SeparatorFunc
	inputName string
	outNames  []string
}

func New(bufferSize int) *pipeline {
	return &pipeline{
		bufferSize:   bufferSize,
		decorators:   make([]DecoratorFunc, 0),
		multiplexers: make([]multiplexer, 0),
		separators:   make([]separator, 0),
		inputs:       make(map[string]chan string),
		outputs:      make(map[string]chan string),
	}
}

func (p *pipeline) getInput(name string) chan string {
	p.lock.Lock()
	defer p.lock.Unlock()

	ch, ok := p.inputs[name]
	if ok {
		return ch
	}

	ch = make(chan string, p.bufferSize)
	p.inputs[name] = ch

	return ch
}

func (p *pipeline) getOutput(name string) chan string {
	p.lock.Lock()
	defer p.lock.Unlock()

	ch, ok := p.outputs[name]
	if ok {
		return ch
	}

	ch = make(chan string, p.bufferSize)
	p.outputs[name] = ch

	return ch
}

func (p *pipeline) RegisterDecorator(
	fn DecoratorFunc,
	input string,
	output string,
) error {
	if fn == nil {
		return ErrNilFunc
	}

	inChan := p.getInput(input)
	outChan := p.getOutput(output)

	p.lock.Lock()
	defer p.lock.Unlock()

	p.decorators = append(
		p.decorators,
		func(ctx context.Context, _, _ chan string) error {
			return fn(ctx, inChan, outChan)
		},
	)

	return nil
}

func (p *pipeline) RegisterMultiplexer(
	fn MultiplexerFunc,
	inputs []string,
	output string,
) error {
	if fn == nil {
		return ErrNilFunc
	}

	if len(inputs) == 0 {
		return errors.New("empty inputs")
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.multiplexers = append(
		p.multiplexers,
		multiplexer{
			fn:         fn,
			inputNames: inputs,
			outputName: output,
		},
	)

	return nil
}

func (p *pipeline) RegisterSeparator(
	fn SeparatorFunc,
	input string,
	outputs []string,
) error {
	if fn == nil {
		return ErrNilFunc
	}

	if len(outputs) == 0 {
		return errors.New("empty outputs")
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.separators = append(
		p.separators,
		separator{
			fn:        fn,
			inputName: input,
			outNames:  outputs,
		},
	)

	return nil
}

func (p *pipeline) Send(name string, value string) error {
	ch, ok := p.inputs[name]
	if !ok {
		return ErrPipelineClosed
	}

	select {
	case ch <- value:
		return nil
	default:
		return ErrPipelineClosed
	}
}

func (p *pipeline) Recv(name string) (string, error) {
	ch, ok := p.outputs[name]
	if !ok {
		return "", ErrPipelineClosed
	}

	res, ok := <-ch
	if !ok {
		return "", ErrPipelineClosed
	}

	return res, nil
}

func (p *pipeline) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, decorator := range p.decorators {
		fn := decorator

		group.Go(func() error {
			err := fn(ctx, nil, nil)
			if err != nil {
				return fmt.Errorf("decorator failed: %w", err)
			}

			return nil
		})
	}

	for _, mul := range p.multiplexers {
		item := mul

		group.Go(func() error {
			ins := make([]chan string, 0, len(item.inputNames))

			for _, n := range item.inputNames {
				ins = append(ins, p.inputs[n])
			}

			out := p.getOutput(item.outputName)

			err := item.fn(ctx, ins, out)
			if err != nil {
				return fmt.Errorf("multiplexer failed: %w", err)
			}

			close(out)

			return nil
		})
	}

	for _, sep := range p.separators {
		item := sep

		group.Go(func() error {
			in := p.getInput(item.inputName)

			outs := make([]chan string, 0, len(item.outNames))

			for _, n := range item.outNames {
				outs = append(outs, p.getOutput(n))
			}

			err := item.fn(ctx, in, outs)
			if err != nil {
				return fmt.Errorf("separator failed: %w", err)
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("pipeline error: %w", err)
	}

	return nil
}

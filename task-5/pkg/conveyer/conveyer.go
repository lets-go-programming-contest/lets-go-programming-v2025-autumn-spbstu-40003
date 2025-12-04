package conveyer

import (
	"context"
	"errors"
	"sync"
)

const Undefined = "undefined"

type Conveyer interface {
	RegisterDecorator(
		fn func(context.Context, chan string, chan string) error,
		input string,
		output string,
	)

	RegisterMultiplexer(
		fn func(context.Context, []chan string, chan string) error,
		inputs []string,
		output string,
	)

	RegisterSeparator(
		fn func(context.Context, chan string, []chan string) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type pipeline struct {
	size    int
	mu      sync.Mutex
	chans   map[string]chan string
	workers []func(context.Context) error
}

func New(size int) Conveyer {
	return &pipeline{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (p *pipeline) getOrCreate(name string) chan string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, ok := p.chans[name]; ok {
		return ch
	}

	ch := make(chan string, p.size)
	p.chans[name] = ch
	return ch
}

func (p *pipeline) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	in := p.getOrCreate(input)
	out := p.getOrCreate(output)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (p *pipeline) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	var ins []chan string
	for _, name := range inputs {
		ins = append(ins, p.getOrCreate(name))
	}
	out := p.getOrCreate(output)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (p *pipeline) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	in := p.getOrCreate(input)

	var outs []chan string
	for _, name := range outputs {
		outs = append(outs, p.getOrCreate(name))
	}

	p.workers = append(p.workers, func(ctx context.Context) error {
		return fn(ctx, in, outs)
	})
}

func (p *pipeline) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, len(p.workers))

	for _, w := range p.workers {
		wg.Add(1)
		go func(job func(context.Context) error) {
			defer wg.Done()
			if err := job(ctx); err != nil {
				errCh <- err
				cancel()
			}
		}(w)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var first error
	for err := range errCh {
		if err != nil && first == nil {
			first = err
		}
	}

	p.closeAll()
	return first
}

func (p *pipeline) closeAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, ch := range p.chans {
		close(ch)
	}
}

func (p *pipeline) Send(input string, data string) error {
	p.mu.Lock()
	ch, ok := p.chans[input]
	p.mu.Unlock()

	if !ok {
		return errors.New("chan not found")
	}

	ch <- data
	return nil
}

func (p *pipeline) Recv(output string) (string, error) {
	p.mu.Lock()
	ch, ok := p.chans[output]
	p.mu.Unlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return Undefined, nil
	}

	return val, nil
}

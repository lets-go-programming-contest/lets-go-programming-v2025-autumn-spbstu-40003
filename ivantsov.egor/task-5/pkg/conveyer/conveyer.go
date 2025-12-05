package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type (
	JobFunc      func(context.Context, chan string, chan string) error
	JobsFunc     func(context.Context, []chan string, chan string) error
	SplitterFunc func(context.Context, chan string, []chan string) error
)

type Conveyer interface {
	AddJob(input string, output string, job JobFunc) error
	AddJobs(inputs []string, output string, job JobsFunc) error
	AddSplitter(input string, outputs []string, splitter SplitterFunc) error
	Run(ctx context.Context) error
}

type pipeline struct {
	bufferSize int
	mu         sync.Mutex
	chans      map[string]chan string
	workers    []func(context.Context) error
}

func New(bufferSize int) *pipeline {
	return &pipeline{
		bufferSize: bufferSize,
		mu:         sync.Mutex{},
		chans:      make(map[string]chan string),
		workers:    make([]func(context.Context) error, 0),
	}
}

func (p *pipeline) getOrCreate(name string) chan string {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch, ok := p.chans[name]
	if ok {
		return ch
	}

	newCh := make(chan string, p.bufferSize)
	p.chans[name] = newCh

	return newCh
}

func (p *pipeline) AddJob(input string, output string, job JobFunc) error {
	inputChan := p.getOrCreate(input)
	outputChan := p.getOrCreate(output)

	p.workers = append(
		p.workers,
		func(ctx context.Context) error {
			return job(ctx, inputChan, outputChan)
		},
	)

	return nil
}

func (p *pipeline) AddJobs(inputs []string, output string, job JobsFunc) error {
	outputChan := p.getOrCreate(output)

	inputChans := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		inputChans = append(inputChans, p.getOrCreate(name))
	}

	p.workers = append(
		p.workers,
		func(ctx context.Context) error {
			return job(ctx, inputChans, outputChan)
		},
	)

	return nil
}

func (p *pipeline) AddSplitter(input string, outputs []string, splitter SplitterFunc) error {
	inputChan := p.getOrCreate(input)

	outputChans := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outputChans = append(outputChans, p.getOrCreate(name))
	}

	p.workers = append(
		p.workers,
		func(ctx context.Context) error {
			return splitter(ctx, inputChan, outputChans)
		},
	)

	return nil
}

func (p *pipeline) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	errChan := make(chan error, len(p.workers))

	for _, worker := range p.workers {
		wg.Add(1)

		current := worker

		go func() {
			defer wg.Done()
			if err := current(ctx); err != nil {
				errChan <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return errors.New(ctx.Err().Error())
	}

	return nil
}

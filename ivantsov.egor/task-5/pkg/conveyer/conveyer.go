package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type (
	JobSingle   func(ctx context.Context, input chan string, output chan string) error
	JobMultiIn  func(ctx context.Context, inputs []chan string, output chan string) error
	JobMultiOut func(ctx context.Context, input chan string, outputs []chan string) error
)

type Conveyer interface {
	AddJob(inputName string, outputName string, function JobSingle)
	AddJobs(inputNames []string, outputName string, function JobMultiIn)
	AddSplitter(inputName string, outputNames []string, function JobMultiOut)
	Run(ctx context.Context) error
}

type pipeline struct {
	size    int
	chans   map[string]chan string
	workers []func(ctx context.Context) error
	mu      sync.Mutex
}

func New(size int) *pipeline {
	return &pipeline{
		size:    size,
		chans:   make(map[string]chan string),
		workers: []func(ctx context.Context) error{},
	}
}

func (p *pipeline) getOrCreate(name string) chan string {
	p.mu.Lock()
	defer p.mu.Unlock()

	channel, ok := p.chans[name]
	if ok {
		return channel
	}

	channel = make(chan string, p.size)
	p.chans[name] = channel

	return channel
}

func (p *pipeline) AddJob(inputName string, outputName string, function JobSingle) {
	inputChan := p.getOrCreate(inputName)
	outputChan := p.getOrCreate(outputName)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return function(ctx, inputChan, outputChan)
	})
}

func (p *pipeline) AddJobs(inputNames []string, outputName string, function JobMultiIn) {
	ins := make([]chan string, 0, len(inputNames))

	for _, name := range inputNames {
		ins = append(ins, p.getOrCreate(name))
	}

	outputChan := p.getOrCreate(outputName)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return function(ctx, ins, outputChan)
	})
}

func (p *pipeline) AddSplitter(inputName string, outputNames []string, function JobMultiOut) {
	inputChan := p.getOrCreate(inputName)
	outputNamesCopy := append([]string(nil), outputNames...)
	outputs := make([]chan string, 0, len(outputNamesCopy))

	for _, name := range outputNamesCopy {
		outputs = append(outputs, p.getOrCreate(name))
	}

	p.workers = append(p.workers, func(ctx context.Context) error {
		return function(ctx, inputChan, outputs)
	})
}

func (p *pipeline) Run(ctx context.Context) error {
	waitGroup := &sync.WaitGroup{}
	errCh := make(chan error, len(p.workers))

	for _, worker := range p.workers {
		waitGroup.Add(1)

		go func(job func(context.Context) error) {
			defer waitGroup.Done()

			if err := job(ctx); err != nil {
				errCh <- err
			}
		}(worker)
	}

	waitGroup.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

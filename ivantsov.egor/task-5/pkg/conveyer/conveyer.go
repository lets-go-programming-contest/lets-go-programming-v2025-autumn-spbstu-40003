package conveyer

import (
	"context"
	"sync"
)

type Conveyer interface {
	AddJob(
		inputName string,
		outputName string,
		handler func(context.Context, chan string, chan string) error,
	)

	AddJobs(
		inputNames []string,
		outputName string,
		handler func(context.Context, []chan string, chan string) error,
	)

	AddSplitter(
		inputName string,
		outputNames []string,
		handler func(context.Context, chan string, []chan string) error,
	)

	Run(ctx context.Context) error
}

type pipeline struct {
	chans   map[string]chan string
	workers []func(context.Context) error
	mu      sync.Mutex
}

func New(bufferSize int) Conveyer {
	return &pipeline{
		chans:   make(map[string]chan string),
		workers: nil,
		mu:      sync.Mutex{},
	}
}

func (p *pipeline) getOrCreate(name string) chan string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, ok := p.chans[name]; ok {
		return ch
	}

	created := make(chan string)
	p.chans[name] = created

	return created
}

func (p *pipeline) AddJob(
	inputName string,
	outputName string,
	handler func(context.Context, chan string, chan string) error,
) {
	input := p.getOrCreate(inputName)
	output := p.getOrCreate(outputName)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, input, output)
	})
}

func (p *pipeline) AddJobs(
	inputNames []string,
	outputName string,
	handler func(context.Context, []chan string, chan string) error,
) {
	inputs := make([]chan string, len(inputNames))

	for i, name := range inputNames {
		inputs[i] = p.getOrCreate(name)
	}

	output := p.getOrCreate(outputName)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, inputs, output)
	})
}

func (p *pipeline) AddSplitter(
	inputName string,
	outputNames []string,
	handler func(context.Context, chan string, []chan string) error,
) {
	input := p.getOrCreate(inputName)
	outputNamesCopy := append([]string(nil), outputNames...)

	outputs := make([]chan string, len(outputNamesCopy))

	for i, name := range outputNamesCopy {
		outputs[i] = p.getOrCreate(name)
	}

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, input, outputs)
	})
}

func (p *pipeline) Run(ctx context.Context) error {
	errChan := make(chan error, len(p.workers))
	waitGroup := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, worker := range p.workers {
		waitGroup.Add(1)

		go func(job func(context.Context) error) {
			defer waitGroup.Done()

			if err := job(ctx); err != nil {
				errChan <- err
				cancel()
			}
		}(worker)
	}

	waitGroup.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

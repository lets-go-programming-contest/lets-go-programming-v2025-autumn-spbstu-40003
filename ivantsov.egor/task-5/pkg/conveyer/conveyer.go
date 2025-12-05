package conveyer

import (
	"context"
	"errors"
	"sync"
)

var errChanNotFound = errors.New("channel not found")

type Job func(context.Context) error

type Conveyer interface {
	AddJob(
		string,
		string,
		func(context.Context, chan string, chan string) error,
	)

	AddJobs(
		[]string,
		string,
		func(context.Context, []chan string, chan string) error,
	)

	AddSplitter(
		string,
		[]string,
		func(context.Context, chan string, []chan string) error,
	)

	Run(context.Context) error
}

type pipeline struct {
	size    int
	chans   map[string]chan string
	workers []Job
	mu      sync.Mutex
}

func New(size int) *pipeline {
	return &pipeline{
		size:    size,
		chans:   make(map[string]chan string),
		workers: make([]Job, 0),
		mu:      sync.Mutex{},
	}
}

func (pipelineInstance *pipeline) getOrCreate(name string) chan string {
	pipelineInstance.mu.Lock()
	defer pipelineInstance.mu.Unlock()

	if channel, exists := pipelineInstance.chans[name]; exists {
		return channel
	}

	channel := make(chan string, pipelineInstance.size)

	pipelineInstance.chans[name] = channel

	return channel
}

func (pipelineInstance *pipeline) AddJob(
	input string,
	output string,
	function func(context.Context, chan string, chan string) error,
) {

	inputChan := pipelineInstance.getOrCreate(input)
	outputChan := pipelineInstance.getOrCreate(output)

	pipelineInstance.workers = append(
		pipelineInstance.workers,
		func(ctx context.Context) error {
			defer close(outputChan)
			return function(ctx, inputChan, outputChan)
		},
	)
}

func (pipelineInstance *pipeline) AddJobs(
	inputNames []string,
	output string,
	function func(context.Context, []chan string, chan string) error,
) {

	inputs := make([]chan string, 0, len(inputNames))

	for _, name := range inputNames {
		inputs = append(inputs, pipelineInstance.getOrCreate(name))
	}

	outputChan := pipelineInstance.getOrCreate(output)

	pipelineInstance.workers = append(
		pipelineInstance.workers,
		func(ctx context.Context) error {
			defer close(outputChan)
			return function(ctx, inputs, outputChan)
		},
	)
}

func (pipelineInstance *pipeline) AddSplitter(
	input string,
	outputs []string,
	function func(context.Context, chan string, []chan string) error,
) {

	inputChan := pipelineInstance.getOrCreate(input)

	outputChannels := make([]chan string, 0, len(outputs))

	for _, name := range outputs {
		outputChannels = append(
			outputChannels,
			pipelineInstance.getOrCreate(name),
		)
	}

	pipelineInstance.workers = append(
		pipelineInstance.workers,
		func(ctx context.Context) error {
			for _, ch := range outputChannels {
				defer close(ch)
			}

			return function(ctx, inputChan, outputChannels)
		},
	)
}

func (pipelineInstance *pipeline) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChannel := make(chan error, len(pipelineInstance.workers))

	var waitGroup sync.WaitGroup

	for _, job := range pipelineInstance.workers {
		waitGroup.Add(1)

		go func(jobFunction Job) {
			defer waitGroup.Done()

			if err := jobFunction(ctx); err != nil {
				errChannel <- err
				cancel()
			}

		}(job)
	}

	waitGroup.Wait()
	close(errChannel)

	for err := range errChannel {
		return err
	}

	return nil
}

func (pipelineInstance *pipeline) getChannel(
	name string,
) (chan string, error) {

	channel, exists := pipelineInstance.chans[name]

	if !exists {
		return nil, errChanNotFound
	}

	return channel, nil
}

func (pipelineInstance *pipeline) Receive(
	name string,
) (<-chan string, error) {

	ch, err := pipelineInstance.getChannel(name)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

package conveyer

import (
	"context"
	"sync"
)

type Conveyer interface {
	AddJob(
		inputName string,
		outputName string,
		handler func(ctx context.Context, inputChannel chan string, outputChannel chan string) error,
	)
	AddJobs(
		inputNames []string,
		outputName string,
		handler func(ctx context.Context, inputChannels []chan string, outputChannel chan string) error,
	)
	AddSplitter(
		inputName string,
		outputNames []string,
		handler func(ctx context.Context, inputChannel chan string, outputChannels []chan string) error,
	)
	Run(ctx context.Context) error
}

type workerFunc func(context.Context) error

type pipeline struct {
	chans   map[string]chan string
	size    int
	mutex   *sync.Mutex
	workers []workerFunc
}

//nolint:ireturn
func New(bufferSize int) Conveyer {
	return &pipeline{
		chans:   make(map[string]chan string),
		size:    bufferSize,
		mutex:   &sync.Mutex{},
		workers: nil,
	}
}

func (p *pipeline) AddJob(
	inputName string,
	outputName string,
	handler func(ctx context.Context, inputChannel chan string, outputChannel chan string) error,
) {
	inputChannelName := inputName
	outputChannelName := outputName

	job := func(ctx context.Context) error {
		inputChannel := p.getOrCreateChannel(inputChannelName)
		outputChannel := p.getOrCreateChannel(outputChannelName)

		return handler(ctx, inputChannel, outputChannel)
	}

	p.workers = append(p.workers, job)
}

func (p *pipeline) AddJobs(
	inputNames []string,
	outputName string,
	handler func(ctx context.Context, inputChannels []chan string, outputChannel chan string) error,
) {
	namesCopy := append([]string(nil), inputNames...)
	outputChannelName := outputName

	job := func(ctx context.Context) error {
		inputChannels := make([]chan string, 0, len(namesCopy))
		for _, channelName := range namesCopy {
			inputChannels = append(inputChannels, p.getOrCreateChannel(channelName))
		}

		outputChannel := p.getOrCreateChannel(outputChannelName)

		return handler(ctx, inputChannels, outputChannel)
	}

	p.workers = append(p.workers, job)
}

func (p *pipeline) AddSplitter(
	inputName string,
	outputNames []string,
	handler func(ctx context.Context, inputChannel chan string, outputChannels []chan string) error,
) {
	inputChannelName := inputName
	outputNamesCopy := append([]string(nil), outputNames...)

	job := func(ctx context.Context) error {
		inputChannel := p.getOrCreateChannel(inputChannelName)

		outputChannels := make([]chan string, 0, len(outputNamesCopy))
		for _, channelName := range outputNamesCopy {
			outputChannels = append(outputChannels, p.getOrCreateChannel(channelName))
		}

		return handler(ctx, inputChannel, outputChannels)
	}

	p.workers = append(p.workers, job)
}

func (p *pipeline) Run(ctx context.Context) error {
	waitGroup := &sync.WaitGroup{}
	errorChannel := make(chan error, len(p.workers))

	for _, workerJob := range p.workers {
		jobCopy := workerJob

		waitGroup.Add(1)

		go func(job workerFunc) {
			defer waitGroup.Done()

			if err := job(ctx); err != nil {
				errorChannel <- err
			}
		}(jobCopy)
	}

	waitGroup.Wait()
	close(errorChannel)

	for err := range errorChannel {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *pipeline) getOrCreateChannel(name string) chan string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	channel, exists := p.chans[name]
	if !exists {
		channel = make(chan string, p.size)
		p.chans[name] = channel
	}

	return channel
}

package conveyor

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type WorkerFunc func(ctx context.Context) error

type Conveyor struct {
	channelBufferSize int
	channelsByID      map[string]chan string
	workerFunctions   []WorkerFunc

	inputIDs  map[string]struct{}
	outputIDs map[string]struct{}
}

var ErrChanNotFound = errors.New("chan not found")

func New(channelBufferSize int) *Conveyor {
	return &Conveyor{
		channelBufferSize: channelBufferSize,
		channelsByID:      make(map[string]chan string),
		workerFunctions:   make([]WorkerFunc, 0),
		inputIDs:          make(map[string]struct{}),
		outputIDs:         make(map[string]struct{}),
	}
}

func (conveyorInstance *Conveyor) markInput(channelID string) {
	conveyorInstance.inputIDs[channelID] = struct{}{}
}

func (conveyorInstance *Conveyor) markOutput(channelID string) {
	conveyorInstance.outputIDs[channelID] = struct{}{}
}

func (conveyorInstance *Conveyor) RegisterDecorator(
	handlerFunc func(
		ctx context.Context,
		inputChannel chan string,
		outputChannel chan string,
	) error,
	inputChannelID string,
	outputChannelID string,
) {
	inputChannel := conveyorInstance.getOrCreateChannel(inputChannelID)
	outputChannel := conveyorInstance.getOrCreateChannel(outputChannelID)

	conveyorInstance.markInput(inputChannelID)
	conveyorInstance.markOutput(outputChannelID)

	workerWrapper := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannel)
	}

	conveyorInstance.workerFunctions = append(conveyorInstance.workerFunctions, workerWrapper)
}

func (conveyorInstance *Conveyor) RegisterMultiplexer(
	handlerFunc func(
		ctx context.Context,
		inputChannels []chan string,
		outputChannel chan string,
	) error,
	inputChannelIDs []string,
	outputChannelID string,
) {
	inputChannels := make([]chan string, len(inputChannelIDs))

	for index, channelID := range inputChannelIDs {
		inputChannels[index] = conveyorInstance.getOrCreateChannel(channelID)
		conveyorInstance.markInput(channelID)
	}

	outputChannel := conveyorInstance.getOrCreateChannel(outputChannelID)
	conveyorInstance.markOutput(outputChannelID)

	workerWrapper := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannels, outputChannel)
	}

	conveyorInstance.workerFunctions = append(conveyorInstance.workerFunctions, workerWrapper)
}

func (conveyorInstance *Conveyor) RegisterSeparator(
	handlerFunc func(
		ctx context.Context,
		inputChannel chan string,
		outputChannels []chan string,
	) error,
	inputChannelID string,
	outputChannelIDs []string,
) {
	inputChannel := conveyorInstance.getOrCreateChannel(inputChannelID)
	conveyorInstance.markInput(inputChannelID)

	outputChannels := make([]chan string, len(outputChannelIDs))

	for index, channelID := range outputChannelIDs {
		outputChannels[index] = conveyorInstance.getOrCreateChannel(channelID)
		conveyorInstance.markOutput(channelID)
	}

	workerWrapper := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannels)
	}

	conveyorInstance.workerFunctions = append(conveyorInstance.workerFunctions, workerWrapper)
}

func (conveyorInstance *Conveyor) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup

	errorChannel := make(chan error, len(conveyorInstance.workerFunctions))

	for _, workerFunc := range conveyorInstance.workerFunctions {
		localWorkerFunc := workerFunc

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			if err := localWorkerFunc(ctx); err != nil {
				errorChannel <- err
			}
		}()
	}

	doneChannel := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(doneChannel)
	}()

	select {
	case err := <-errorChannel:
		conveyorInstance.closeRootInputChannels()
		waitGroup.Wait()

		if err == nil {
			return nil
		}

		return fmt.Errorf("conveyor worker failed: %w", err)

	case <-doneChannel:
		return nil

	case <-ctx.Done():
		conveyorInstance.closeRootInputChannels()
		waitGroup.Wait()

		return nil
	}
}

func (conveyorInstance *Conveyor) Send(inputChannelID string, data string) error {
	channel, exists := conveyorInstance.channelsByID[inputChannelID]
	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (conveyorInstance *Conveyor) Recv(outputChannelID string) (string, error) {
	channel, exists := conveyorInstance.channelsByID[outputChannelID]
	if !exists {
		return "", ErrChanNotFound
	}

	value, hasValue := <-channel
	if !hasValue {
		return "undefined", nil
	}

	return value, nil
}

func (conveyorInstance *Conveyor) getOrCreateChannel(channelID string) chan string {
	channel, exists := conveyorInstance.channelsByID[channelID]
	if exists {
		return channel
	}

	newChannel := make(chan string, conveyorInstance.channelBufferSize)
	conveyorInstance.channelsByID[channelID] = newChannel

	return newChannel
}

func (conveyorInstance *Conveyor) closeRootInputChannels() {
	for channelID, channel := range conveyorInstance.channelsByID {
		_, isInput := conveyorInstance.inputIDs[channelID]
		_, isOutput := conveyorInstance.outputIDs[channelID]

		if isInput && !isOutput {
			func(ch chan string) {
				defer func() {
					_ = recover()
				}()

				close(ch)
			}(channel)
		}
	}
}

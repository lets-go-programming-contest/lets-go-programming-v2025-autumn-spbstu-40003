package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type channelMap struct {
	sync.RWMutex
	m map[string]chan string
}

type ConveyerImpl struct {
	channels    *channelMap
	tasks       []func(context.Context) error
	channelSize int
}

func New(size int) Conveyer {
	return &ConveyerImpl{
		channels: &channelMap{
			RWMutex: sync.RWMutex{},
			m:       make(map[string]chan string),
		},
		channelSize: size,
		tasks:       make([]func(context.Context) error, 0),
	}
}

func (conv *ConveyerImpl) getOrCreateChannel(name string) chan string {
	conv.channels.Lock()
	defer conv.channels.Unlock()

	if channel, exists := conv.channels.m[name]; exists {
		return channel
	}

	channel := make(chan string, conv.channelSize)
	conv.channels.m[name] = channel

	return channel
}

func (conv *ConveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	inCh := conv.getOrCreateChannel(inputName)
	outCh := conv.getOrCreateChannel(outputName)

	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}

	conv.tasks = append(conv.tasks, task)
}

func (conv *ConveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	inputChannels := make([]chan string, 0, len(inputNames))

	for _, name := range inputNames {
		inputChannels = append(inputChannels, conv.getOrCreateChannel(name))
	}

	outCh := conv.getOrCreateChannel(outputName)

	task := func(ctx context.Context) error {
		return fn(ctx, inputChannels, outCh)
	}

	conv.tasks = append(conv.tasks, task)
}

func (conv *ConveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	inCh := conv.getOrCreateChannel(inputName)

	outputChannels := make([]chan string, 0, len(outputNames))

	for _, name := range outputNames {
		outputChannels = append(outputChannels, conv.getOrCreateChannel(name))
	}

	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outputChannels)
	}

	conv.tasks = append(conv.tasks, task)
}

func (conv *ConveyerImpl) Run(ctx context.Context) error {
	errGroup, groupCtx := errgroup.WithContext(ctx)

	for _, task := range conv.tasks {
		taskFunc := task

		errGroup.Go(func() error {
			return taskFunc(groupCtx)
		})
	}

	err := errGroup.Wait()

	conv.closeAllChannels()

	if err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}

func (conv *ConveyerImpl) closeAllChannels() {
	conv.channels.Lock()
	defer conv.channels.Unlock()

	for _, channel := range conv.channels.m {
		func() {
			defer func() {
				_ = recover()
			}()

			close(channel)
		}()
	}
}

func (conv *ConveyerImpl) Send(inputName string, data string) error {
	conv.channels.RLock()
	channel, exists := conv.channels.m[inputName]
	conv.channels.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	defer func() {
		_ = recover()
	}()

	channel <- data

	return nil
}

func (conv *ConveyerImpl) Recv(outputName string) (string, error) {
	conv.channels.RLock()
	channel, exists := conv.channels.m[outputName]
	conv.channels.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	val, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return val, nil
}

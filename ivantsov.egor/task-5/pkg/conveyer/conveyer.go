package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound   = errors.New("chan not found")
	ErrChanFull       = errors.New("chan is full")
	ErrNotInitialized = errors.New("conveyer not initialized")
)

const StrUndefined = "undefined"

type decoratorItem struct {
	handlerFunc func(context.Context, chan string, chan string) error
	inputName   string
	outputName  string
}

type multiplexerItem struct {
	handlerFunc func(context.Context, []chan string, chan string) error
	inputNames  []string
	outputName  string
}

type separatorItem struct {
	handlerFunc func(context.Context, chan string, []chan string) error
	inputName   string
	outputNames []string
}

type Conveyer struct {
	bufferSize     int
	channelsMap    map[string]chan string
	decoratorSet   []decoratorItem
	multiplexerSet []multiplexerItem
	separatorSet   []separatorItem
	channelLock    sync.RWMutex
	initialized    bool
}

func New(bufferSize int) *Conveyer {
	return &Conveyer{
		bufferSize:     bufferSize,
		channelsMap:    make(map[string]chan string),
		decoratorSet:   make([]decoratorItem, 0),
		multiplexerSet: make([]multiplexerItem, 0),
		separatorSet:   make([]separatorItem, 0),
		channelLock:    sync.RWMutex{},
		initialized:    true,
	}
}

func (conv *Conveyer) ensureChannel(channelName string) chan string {
	conv.channelLock.Lock()
	defer conv.channelLock.Unlock()

	existingChannel, found := conv.channelsMap[channelName]
	if found {
		return existingChannel
	}

	createdChannel := make(chan string, conv.bufferSize)
	conv.channelsMap[channelName] = createdChannel

	return createdChannel
}

func (conv *Conveyer) getChannel(channelName string) (chan string, bool) {
	conv.channelLock.RLock()
	defer conv.channelLock.RUnlock()

	existingChannel, found := conv.channelsMap[channelName]

	return existingChannel, found
}

func (conv *Conveyer) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	conv.decoratorSet = append(conv.decoratorSet, decoratorItem{
		handlerFunc: handlerFunc,
		inputName:   inputName,
		outputName:  outputName,
	})

	conv.ensureChannel(inputName)
	conv.ensureChannel(outputName)
}

func (conv *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	conv.multiplexerSet = append(conv.multiplexerSet, multiplexerItem{
		handlerFunc: handlerFunc,
		inputNames:  inputNames,
		outputName:  outputName,
	})

	for _, inputItem := range inputNames {
		conv.ensureChannel(inputItem)
	}

	conv.ensureChannel(outputName)
}

func (conv *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	conv.separatorSet = append(conv.separatorSet, separatorItem{
		handlerFunc: handlerFunc,
		inputName:   inputName,
		outputNames: outputNames,
	})

	conv.ensureChannel(inputName)

	for _, outputItem := range outputNames {
		conv.ensureChannel(outputItem)
	}
}

func (conv *Conveyer) Run(parentContext context.Context) error {
	if !conv.initialized {
		return ErrNotInitialized
	}

	internalContext, cancelFunction := context.WithCancel(parentContext)
	defer cancelFunction()

	var workersGroup sync.WaitGroup

	errorChannel := make(chan error, 1)

	startWorker := func(workerFunc func()) {
		workersGroup.Add(1)

		go func() {
			defer workersGroup.Done()
			workerFunc()
		}()
	}

	conv.launchDecorators(
		internalContext,
		startWorker,
		errorChannel,
	)

	conv.launchMultiplexers(
		internalContext,
		startWorker,
		errorChannel,
	)

	conv.launchSeparators(
		internalContext,
		startWorker,
		errorChannel,
	)

	go func() {
		workersGroup.Wait()

		conv.closeAllChannels()

		close(errorChannel)
	}()

	select {
	case receivedErr := <-errorChannel:
		cancelFunction()

		return receivedErr
	case <-internalContext.Done():

		return nil
	}
}

func (conv *Conveyer) launchDecorators(
	internalContext context.Context,
	startWorker func(func()),
	errorChannel chan error,
) {
	for _, decoratorEntry := range conv.decoratorSet {
		entry := decoratorEntry

		startWorker(func() {
			execErr := entry.handlerFunc(
				internalContext,
				conv.ensureChannel(entry.inputName),
				conv.ensureChannel(entry.outputName),
			)
			if execErr != nil {
				conv.sendError(errorChannel, execErr)
			}
		})
	}
}

func (conv *Conveyer) launchMultiplexers(
	internalContext context.Context,
	startWorker func(func()),
	errorChannel chan error,
) {
	for _, multiplexerEntry := range conv.multiplexerSet {
		entry := multiplexerEntry

		startWorker(func() {
			inputChannels := make([]chan string, len(entry.inputNames))
			for indexValue, channelName := range entry.inputNames {
				inputChannels[indexValue] = conv.ensureChannel(channelName)
			}

			execErr := entry.handlerFunc(
				internalContext,
				inputChannels,
				conv.ensureChannel(entry.outputName),
			)
			if execErr != nil {
				conv.sendError(errorChannel, execErr)
			}
		})
	}
}

func (conv *Conveyer) launchSeparators(
	internalContext context.Context,
	startWorker func(func()),
	errorChannel chan error,
) {
	for _, separatorEntry := range conv.separatorSet {
		entry := separatorEntry

		startWorker(func() {
			outputChannels := make([]chan string, len(entry.outputNames))
			for indexValue, channelName := range entry.outputNames {
				outputChannels[indexValue] = conv.ensureChannel(channelName)
			}

			execErr := entry.handlerFunc(
				internalContext,
				conv.ensureChannel(entry.inputName),
				outputChannels,
			)
			if execErr != nil {
				conv.sendError(errorChannel, execErr)
			}
		})
	}
}

func (conv *Conveyer) Send(channelName string, payload string) error {
	existingChannel, found := conv.getChannel(channelName)
	if !found {
		return ErrChanNotFound
	}

	select {
	case existingChannel <- payload:

		return nil
	default:

		return ErrChanFull
	}
}

func (conv *Conveyer) Recv(channelName string) (string, error) {
	existingChannel, found := conv.getChannel(channelName)
	if !found {
		return "", ErrChanNotFound
	}

	receivedValue, channelOpen := <-existingChannel
	if !channelOpen {
		return StrUndefined, nil
	}

	return receivedValue, nil
}

func (conv *Conveyer) sendError(errorChannel chan error, execErr error) {
	select {
	case errorChannel <- execErr:
	default:
	}
}

func (conv *Conveyer) closeAllChannels() {
	conv.channelLock.Lock()
	defer conv.channelLock.Unlock()

	closedChannelSet := make(map[chan string]struct{})

	for _, currentChannel := range conv.channelsMap {
		if _, exists := closedChannelSet[currentChannel]; exists {
			continue
		}

		close(currentChannel)

		closedChannelSet[currentChannel] = struct{}{}
	}
}

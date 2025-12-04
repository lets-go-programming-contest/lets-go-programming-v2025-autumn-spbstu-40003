package conveyer

import (
	"errors"
)

type Conv struct {
	chanSize int
	inputs   map[string]chan string
	outputs  map[string]chan string

	decorators   []DecoratorStorage
	separators   []SeparatorStorage
	multiplexers []MultiplexerStorage
}

func New(size int) *Conv {
	return &Conv{
		chanSize:     size,
		inputs:       make(map[string]chan string),
		outputs:      make(map[string]chan string),
		decorators:   make([]DecoratorStorage, 0),
		separators:   make([]SeparatorStorage, 0),
		multiplexers: make([]MultiplexerStorage, 0),
	}
}

var errChanNotFound = errors.New("chan not found")

func (conv *Conv) Send(input string, data string) error {
	channel, found := conv.inputs[input]
	if !found {
		return errChanNotFound
	}

	channel <- data

	return nil
}

const undefinedData = "undefined"

func (conv *Conv) Recv(output string) (string, error) {
	channel, found := conv.outputs[output]
	if !found {
		return "", errChanNotFound
	}

	data, open := <-channel
	if !open {
		return undefinedData, nil
	}

	return data, nil
}

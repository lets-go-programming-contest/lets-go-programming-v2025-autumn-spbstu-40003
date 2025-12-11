package conveyer

import "errors"

// Send(), Recv(), Close()

var errChanNotFound = errors.New("chan not found")

const UndefinedValue = "undefined"

func (conv *Conveyer) Send(name string, data string) error {
	channel, exists := conv.channels[name]
	if !exists {
		return errChanNotFound
	}
	channel <- data

	return nil
}

func (conv *Conveyer) Recv(name string) (string, error) {
	channel, exists := conv.channels[name]
	if !exists {
		return "", errChanNotFound
	}

	message, exists := <-channel
	if !exists {
		return UndefinedValue, nil
	}

	return message, nil
}

func (conv *Conveyer) Close(name string) error {
	channel, exists := conv.channels[name]
	if !exists {
		return errChanNotFound
	}
	close(channel)

	return nil
}

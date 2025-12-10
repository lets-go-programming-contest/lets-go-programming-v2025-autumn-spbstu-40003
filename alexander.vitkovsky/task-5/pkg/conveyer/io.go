package conveyer

import "errors"

// Send(), Recv(), Close()

var errChanNotFound = errors.New("chan not found")

const UndefinedValue = "undefined"

func (conv *Conveyer) Send(name string, data string) error {
	ch, ok := conv.channels[name]
	if !ok {
		return errChanNotFound
	}
	ch <- data
	return nil
}

func (conv *Conveyer) Recv(name string) (string, error) {
	ch, ok := conv.channels[name]
	if !ok {
		return "", errChanNotFound
	}

	message, ok := <-ch
	if !ok {
		return UndefinedValue, nil
	}

	return message, nil
}

func (conv *Conveyer) Close(name string) error {
	ch, ok := conv.channels[name]
	if !ok {
		return errChanNotFound
	}
	close(ch)
	return nil
}

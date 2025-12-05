package conveyer

import (
	"errors"
)

/*
	external access to data: Send() and Recv()
	+ Close(input string) - supportive method
*/

const ChanNotFoundMsg = "chan not found"
const UndefinedValue = "undefined"

func (conv *Conveyer) Send(input string, data string) error {
	<-conv.ready

	conv.mutex.Lock()
	ch, ok := conv.channels[input]
	conv.mutex.Unlock()
	if !ok || ch == nil {
		return errors.New(ChanNotFoundMsg)
	}

	select {
	case ch <- data:
		return nil
	case <-conv.runCtx.Done():
		return conv.runCtx.Err()
	}
}

func (conv *Conveyer) Recv(output string) (string, error) {
	<-conv.ready

	conv.mutex.Lock()
	ch, ok := conv.channels[output]
	conv.mutex.Unlock()
	if !ok || ch == nil {
		return "", errors.New(ChanNotFoundMsg)
	}

	select {
	case msg, ok := <-ch:
		if !ok {
			return UndefinedValue, nil
		}
		return msg, nil
	case <-conv.runCtx.Done():
		return "", conv.runCtx.Err()
	}
}

func (conv *Conveyer) Close(input string) error {
	<-conv.ready
	conv.mutex.Lock()
	ch, ok := conv.channels[input]
	conv.mutex.Unlock()
	if !ok || ch == nil {
		return errors.New(ChanNotFoundMsg)
	}
	defer func() { _ = recover() }()
	close(ch)
	return nil
}

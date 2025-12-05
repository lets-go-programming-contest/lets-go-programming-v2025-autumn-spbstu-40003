package conveyer

import (
	"errors"
)

const ChanNotFoundMsg = "chan not found"
const UndefinedValue = "undefined"

func (conv *Conveyer) Send(input string, data string) error {
	<-conv.ready

	conv.mutex.Lock()
	ch, ok := conv.channels[input]
	runCtx := conv.runCtx
	conv.mutex.Unlock()

	if !ok || ch == nil {
		return errors.New(ChanNotFoundMsg)
	}

	select {
	case ch <- data:
		return nil
	case <-runCtx.Done():
		return runCtx.Err()
	}
}

func (conv *Conveyer) Recv(output string) (string, error) {
	// ждём Run()
	<-conv.ready

	conv.mutex.Lock()
	ch, ok := conv.channels[output]
	runCtx := conv.runCtx
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
	case <-runCtx.Done():
		return UndefinedValue, runCtx.Err()
	}
}

func (conv *Conveyer) Close(input string) error {
	conv.mutex.Lock()
	ch, ok := conv.channels[input]
	if !ok || ch == nil {
		conv.mutex.Unlock()
		return errors.New(ChanNotFoundMsg)
	}
	delete(conv.channels, input)
	conv.mutex.Unlock()

	defer func() { _ = recover() }()
	close(ch)
	return nil
}

package conveyer

import (
	"context"
	"errors"
)

func (conv *Conveyer) Run(ctx context.Context) error {
	conv.mutex.Lock()
	if conv.started {
		conv.mutex.Unlock()
		return errors.New("conveyor already started")
	}
	conv.started = true
	conv.mutex.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	conv.runCtx = ctx
	conv.cancel = cancel

	select {
	case <-conv.ready:
		// already closed
	default:
		close(conv.ready)
	}

	for _, handler := range conv.handlers {
		inputChannels := conv.resolveChannels(handler.inputs)
		outputChannels := conv.resolveChannels(handler.outputs)

		conv.wg.Add(1)
		handlerCopy := handler
		go func(handlerCopy handlerConfig, inputChannels, outputChannels []chan string) {
			defer conv.wg.Done()

			var err error
			switch handlerCopy.kind {
			case "decorator":
				fn := handlerCopy.function.(decoratorFn)
				err = fn(ctx, inputChannels[0], outputChannels[0])
			case "separator":
				fn := handlerCopy.function.(separatorFn)
				err = fn(ctx, inputChannels[0], outputChannels)
			case "multiplexer":
				fn := handlerCopy.function.(multiplexerFn)
				err = fn(ctx, inputChannels, outputChannels[0])
			}

			if err != nil {
				select {
				case conv.errCh <- err:
				default:
				}
			}
		}(handlerCopy, inputChannels, outputChannels)
	}

	go func() {
		conv.wg.Wait()
		close(conv.errCh)
	}()

	errC := make(chan error, 1)
	go func() {
		if err, ok := <-conv.errCh; ok {
			select {
			case errC <- err:
			default:
			}
		}
		close(errC)
	}()

	conv.wg.Wait()

	if err, ok := <-errC; ok && err != nil {
		conv.closeAllChannels()
		return err
	}

	return nil
}

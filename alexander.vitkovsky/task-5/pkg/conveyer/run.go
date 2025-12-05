package conveyer

import (
	"context"
	"errors"
)

// pipe executing: Run()

func (conv *Conveyer) Run(ctx context.Context) error {
	if conv.started {
		return errors.New("conveyor already started")
	}
	conv.started = true

	ctx, cancel := context.WithCancel(ctx)
	conv.runCtx = ctx
	conv.cancel = cancel
	close(conv.ready)

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
					conv.cancel()
				case <-ctx.Done(): // ignore
				}
			}
		}(handlerCopy, inputChannels, outputChannels)
	}

	go func() {
		conv.wg.Wait()
		close(conv.errCh)
	}()

	if err, ok := <-conv.errCh; ok {
		conv.closeAllChannels()
		return err
	}

	conv.closeAllChannels()
	return nil
}

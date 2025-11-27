package conveyer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var (
	errChanNotFound = errors.New("chan not found")
	errNoData       = errors.New("no data")
)

const undefined = "undefined"

type Conveyer struct {
	bufSize  int
	channels map[string]chan string
	handlers []func(context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		bufSize:  size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.closeChannels()

	group, groupCtx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		handlerCopy := handler

		group.Go(func() error {
			if err := handlerCopy(groupCtx); err != nil {
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	ch, ok := c.channels[input]
	if !ok {
		return errChanNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, isAvailable := c.channels[output]
	if !isAvailable {
		return "", errChanNotFound
	}

	data, isAvailable := <-ch
	if !isAvailable {
		return undefined, errNoData
	}

	return data, nil
}

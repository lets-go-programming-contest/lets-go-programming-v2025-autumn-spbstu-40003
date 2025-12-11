package conveyer

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Conveyer struct, New(), Run()

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

func (conv *Conveyer) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, handler := range conv.handlers {
		h := handler

		group.Go(func() error {
			return h(groupCtx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer failed: %w", err)
	}

	return nil
}

package handlers

import (
	"context"
	"errors"
	"strings"
)

var errCantBeDecorated = errors.New("can't be decorated")

const (
	noDecorator = "no decorator"
	decorated   = "decorated: "
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecorator) {
				return errCantBeDecorated
			}

			if !strings.HasPrefix(data, decorated) {
				data = decorated + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

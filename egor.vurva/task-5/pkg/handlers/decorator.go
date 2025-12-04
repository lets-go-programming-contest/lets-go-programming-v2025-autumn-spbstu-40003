package handlers

import (
	"context"
	"errors"
	"strings"
)

const (
	decPrefix = "decorated: "
	decFlag   = "no decorator"
)

var errCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, decFlag) {
				return errCantBeDecorated
			}

			if !strings.HasPrefix(data, decPrefix) {
				data = decPrefix + data
			}

			select {
			case <-ctx.Done():
				return nil

			case output <- data:
			}
		}
	}
}

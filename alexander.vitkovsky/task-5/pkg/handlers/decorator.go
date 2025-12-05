package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "
	for {
		select {
		case <-ctx.Done():
			close(output)
			return ctx.Err()

		case msg, ok := <-input:
			if !ok {
				close(output)
				return nil
			}

			if strings.Contains(msg, "no decorator") {
				close(output)
				return ErrCantDecorated
			}

			if strings.HasPrefix(msg, prefix) {
				output <- msg
			} else {
				output <- prefix + msg
			}
		}
	}
}

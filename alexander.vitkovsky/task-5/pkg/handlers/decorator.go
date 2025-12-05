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
				select {
				case output <- msg:
				case <-ctx.Done():
					return ctx.Err()
				}
			} else {
				select {
				case output <- prefix + msg:
				case <-ctx.Done():
					return ctx.Err()
				}
			}

		case <-ctx.Done():
			close(output)
			return ctx.Err()
		}
	}
}

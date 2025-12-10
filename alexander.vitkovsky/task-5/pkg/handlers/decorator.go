package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantDecorate = errors.New("can't decorate")

const prefix = "decorated: "

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

		case message, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(message, "no decorator") {
				return ErrCantDecorate
			}
			if !strings.HasPrefix(message, prefix) {
				message = prefix + message
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- message:
			}
		}
	}
}

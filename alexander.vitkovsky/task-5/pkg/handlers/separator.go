package handlers

import "context"

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	index := 0
	n := len(outputs)
	if n == 0 {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg, ok := <-input:
			if !ok {
				for _, output := range outputs {
					close(output)
				}
				return nil
			}
			outputs[index] <- msg
			index = (index + 1) % n
		}
	}
}

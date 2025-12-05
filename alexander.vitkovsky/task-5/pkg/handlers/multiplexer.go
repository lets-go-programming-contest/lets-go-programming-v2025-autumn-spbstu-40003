package handlers

import (
	"context"
	"strings"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	type msgWrap struct {
		msg string
		ok  bool
	}

	n := len(inputs)
	if n == 0 {
		close(output)
		return nil
	}

	out := make(chan msgWrap)

	for _, ch := range inputs {
		go func(ch chan string) {
			for {
				select {
				case <-ctx.Done():
					out <- msgWrap{ok: false}
					return

				case m, ok := <-ch:
					if !ok {
						out <- msgWrap{ok: false}
						return
					}
					if strings.Contains(m, "no multiplexer") {
						continue
					}
					out <- msgWrap{msg: m, ok: true}
				}
			}
		}(ch)
	}

	closedCount := 0

	for closedCount < n {
		select {
		case <-ctx.Done():
			close(output)
			return ctx.Err()

		case w := <-out:
			if !w.ok {
				closedCount++
				continue
			}
			output <- w.msg
		}
	}

	close(output)
	return nil
}

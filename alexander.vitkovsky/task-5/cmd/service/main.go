package main

import (
	"context"
	"log"

	"github.com/alexpi3/task-5/pkg/conveyer"
	"github.com/alexpi3/task-5/pkg/handlers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const bufferSize = 10
	conv := conveyer.New(bufferSize)

	conv.RegisterSeparator(
		handlers.SeparatorFunc,
		"input",
		[]string{"branch1", "branch2"},
	)

	conv.RegisterDecorator(
		handlers.PrefixDecoratorFunc,
		"branch1",
		"decor1",
	)
	conv.RegisterDecorator(
		handlers.PrefixDecoratorFunc,
		"branch2",
		"decor2",
	)

	conv.RegisterMultiplexer(
		handlers.MultiplexerFunc,
		[]string{"decor1", "decor2"},
		"output",
	)

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			message, err := conv.Recv("output")
			if err != nil {
				log.Println("recv error:", err)
				return
			}

			if message == conveyer.UndefinedValue {
				return
			}
			log.Println("Final:", message)
		}
	}()

	go func() {
		_ = conv.Send("input", "AAA")
		_ = conv.Send("input", "BBB")
		_ = conv.Send("input", "CCC")
		_ = conv.Send("input", "DDD")
		_ = conv.Close("input")
	}()

	if err := conv.Run(ctx); err != nil {
		log.Println("run error:", err)
	}

	<-done
}

package main

import (
	"context"
	"fmt"

	"github.com/alexpi3/task-5/pkg/conveyer"
	"github.com/alexpi3/task-5/pkg/handlers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conv := conveyer.New(10)

	conv.RegisterSeparator(
		handlers.SeparatorFunc,
		"input",
		[]string{"branch1", "branch2"})

	conv.RegisterDecorator(handlers.PrefixDecoratorFunc, "branch1", "decor1")
	conv.RegisterDecorator(handlers.PrefixDecoratorFunc, "branch2", "decor2")

	conv.RegisterMultiplexer(
		handlers.MultiplexerFunc,
		[]string{"decor1", "decor2"},
		"output")

	results := make(chan struct{})
	consumerStarted := make(chan struct{})

	go func() {
		close(consumerStarted)
		defer close(results)

		for {
			msg, err := conv.Recv("output")
			if err != nil {
				return
			}
			if msg == conveyer.UndefinedValue {
				return
			}
			fmt.Println("FINAL:", msg)
		}
	}()

	<-consumerStarted

	go func() {
		_ = conv.Send("input", "AAA")
		_ = conv.Send("input", "BBB")
		_ = conv.Send("input", "CCC")
		_ = conv.Send("input", "DDD")
		_ = conv.Close("input")
	}()

	if err := conv.Run(ctx); err != nil {
		fmt.Println("run error:", err)
	}

	<-results

}

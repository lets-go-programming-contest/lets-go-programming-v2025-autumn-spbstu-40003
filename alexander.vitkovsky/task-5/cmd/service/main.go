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

	go func() {
		if err := conv.Run(ctx); err != nil {
			fmt.Println("run error:", err)
		}
	}()

	go func() {
		if err := conv.Send("input", "AAA"); err != nil {
			fmt.Println(err)
		}
		if err := conv.Send("input", "BBB"); err != nil {
			fmt.Println(err)
		}
		if err := conv.Send("input", "CCC"); err != nil {
			fmt.Println(err)
		}
		if err := conv.Send("input", "DDD"); err != nil {
			fmt.Println(err)
		}

		if err := conv.Close("input"); err != nil {
			fmt.Println(err)
		}
	}()

	for {
		msg, err := conv.Recv("output")
		if err != nil {
			return
		}
		if msg == "undefined" {
			return
		}
		fmt.Println("FINAL:", msg)
	}
}

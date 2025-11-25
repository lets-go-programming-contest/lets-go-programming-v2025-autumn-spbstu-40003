package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nxgmvw/task-5/pkg/conveyer"
	"github.com/nxgmvw/task-5/pkg/handlers"
)

func main() {
	c := conveyer.New(10)

	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "decorated")
	c.RegisterSeparator(handlers.SeparatorFunc, "decorated", []string{"branch1", "branch2"})
	c.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"branch1", "branch2"}, "final")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Println("Pipeline started")
		if err := c.Run(ctx); err != nil {
			log.Printf("Pipeline finished with error: %v\n", err)
		} else {
			fmt.Println("Pipeline finished successfully")
		}
	}()

	inputs := []string{
		"hello",
		"world",
		"no multiplexer this string",
		"data",
	}

	for _, v := range inputs {
		if err := c.Send("input", v); err != nil {
			log.Fatalf("Send error: %v", err)
		}
	}

	for i := 0; i < 3; i++ {
		res, err := c.Recv("final")
		if err != nil {
			log.Printf("Recv error: %v", err)
			break
		}
		fmt.Printf("Received: %s\n", res)
	}

	c.Send("input", "no decorator trigger error")

	time.Sleep(100 * time.Millisecond)

	cancel()
	time.Sleep(100 * time.Millisecond)
}

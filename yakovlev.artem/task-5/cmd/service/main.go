package main

import (
	"context"
	"log"
	"time"

	"github.com/nxgmvw/task-5/pkg/conveyer"
	"github.com/nxgmvw/task-5/pkg/handlers"
)

const (
	chanSize      = 10
	sleepDuration = 100 * time.Millisecond
	receiveCount  = 3
)

func main() {
	conv := conveyer.New(chanSize)

	conv.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "decorated")
	conv.RegisterSeparator(handlers.SeparatorFunc, "decorated", []string{"branch1", "branch2"})
	conv.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"branch1", "branch2"}, "final")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Println("Pipeline started")

		if err := conv.Run(ctx); err != nil {
			log.Printf("Pipeline finished with error: %v\n", err)
		} else {
			log.Println("Pipeline finished successfully")
		}
	}()

	inputs := []string{
		"hello",
		"world",
		"no multiplexer this string",
		"data",
	}

	for _, v := range inputs {
		if err := conv.Send("input", v); err != nil {
			log.Fatalf("Send error: %v", err)
		}
	}

	for range receiveCount {
		res, err := conv.Recv("final")
		if err != nil {
			log.Printf("Recv error: %v", err)

			break
		}

		log.Printf("Received: %s\n", res)
	}

	if err := conv.Send("input", "no decorator trigger error"); err != nil {
		log.Printf("Send trigger error: %v", err)
	}

	time.Sleep(sleepDuration)

	cancel()

	time.Sleep(sleepDuration)
}

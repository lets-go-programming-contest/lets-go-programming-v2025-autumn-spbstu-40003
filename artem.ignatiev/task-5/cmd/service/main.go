package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kryjkaqq/task-5/pkg/conveyer"
	"github.com/kryjkaqq/task-5/pkg/handlers"
)

func main() {
	conv := conveyer.New(5)

	conv.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "decorated_stream")

	conv.RegisterSeparator(handlers.SeparatorFunc, "decorated_stream", []string{"part1", "part2"})

	conv.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"part1", "part2"}, "final_output")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		if err := conv.Run(ctx); err != nil {
			log.Printf("Conveyer stopped with error: %v", err)
		} else {
			log.Println("Conveyer finished successfully")
		}
	}()

	data := []string{"hello", "world", "test", "no multiplexer", "go"}
	for _, d := range data {
		if err := conv.Send("input", d); err != nil {
			log.Printf("Error sending %s: %v", d, err)
		}
	}

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 4; i++ {
		res, err := conv.Recv("final_output")
		if err != nil {
			log.Printf("Error receiving: %v", err)
			continue
		}
		fmt.Printf("Received: %s\n", res)
	}

	fmt.Println("--- Testing Error ---")
	conv.Send("input", "no decorator")

	<-ctx.Done()
}

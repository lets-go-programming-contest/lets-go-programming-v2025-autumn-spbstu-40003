package main

import (
	"context"
	"log"
	"time"

	"github.com/kryjkaqq/task-5/pkg/conveyer"
	"github.com/kryjkaqq/task-5/pkg/handlers"
)

const (
	ChanSize      = 5
	TimeoutSec    = 2
	SleepDuration = 100
	MessagesCount = 4
)

func main() {
	conv := conveyer.New(ChanSize)

	conv.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "decorated_stream")
	conv.RegisterSeparator(handlers.SeparatorFunc, "decorated_stream", []string{"part1", "part2"})
	conv.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"part1", "part2"}, "final_output")

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutSec*time.Second)
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

	time.Sleep(SleepDuration * time.Millisecond)

	for range MessagesCount {
		res, err := conv.Recv("final_output")
		if err != nil {
			log.Printf("Error receiving: %v", err)
			continue
		}
		log.Printf("Received: %s\n", res)
	}

	log.Println("--- Testing Error ---")

	if err := conv.Send("input", "no decorator"); err != nil {
		log.Printf("Error sending trigger: %v", err)
	}

	<-ctx.Done()
}

package conveyer

import "context"

// supportive types for registers.go

type decoratorFn func(
	ctx context.Context,
	input chan string,
	output chan string,
) error

type multiplexerFn func(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error

type separatorFn func(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error

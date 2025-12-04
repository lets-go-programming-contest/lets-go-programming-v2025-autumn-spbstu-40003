package conveyer

import "context"

type DecoratorFunc func(
	ctx context.Context,
	input chan string,
	output chan string,
) error

type SeparatorFunc func(
	ctx context.Context,
	input chan string,
	output []chan string,
) error

type MultiplexerFunc func(
	ctx context.Context,
	input []chan string,
	output chan string,
) error

type DecoratorStorage struct {
	fn     DecoratorFunc
	input  string
	output string
}

type SeparatorStorage struct {
	fn     SeparatorFunc
	input  string
	output []string
}
type MultiplexerStorage struct {
	fn     MultiplexerFunc
	input  []string
	output string
}

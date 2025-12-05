package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type DecoratorFunc func(context.Context, chan string, chan string) error
type SeparatorFunc func(context.Context, chan string, []chan string) error
type MultiplexerFunc func(context.Context, []chan string, chan string) error

type pipeline struct {
	queueSize int
	lock      sync.Mutex

	decorators []decoratorTask
	separators []separatorTask
	mults      []multiplexerTask

	inputMap  map[string]chan string
	outputMap map[string]chan string
}

type decoratorTask struct {
	handler   DecoratorFunc
	inputTag  string
	outputTag string
}

type separatorTask struct {
	handler    SeparatorFunc
	inputTag   string
	outputTags []string
}

type multiplexerTask struct {
	handler   MultiplexerFunc
	inputTags []string
	outputTag string
}

func New(buffer int) *pipeline {
	return &pipeline{
		queueSize: buffer,
		inputMap:  make(map[string]chan string),
		outputMap: make(map[string]chan string),
	}
}

func (p *pipeline) channelIn(name string) chan string {
	val, ok := p.inputMap[name]
	if ok {
		return val
	}

	channel := make(chan string, p.queueSize)
	p.inputMap[name] = channel

	return channel
}

func (p *pipeline) channelOut(name string) chan string {
	val, ok := p.outputMap[name]
	if ok {
		return val
	}

	channel := make(chan string, p.queueSize)
	p.outputMap[name] = channel

	return channel
}

func (p *pipeline) RegisterDecorator(fn DecoratorFunc, inLabel string, outLabel string) error {
	if fn == nil {
		return errors.New("nil decorator")
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.decorators = append(p.decorators, decoratorTask{
		handler:   fn,
		inputTag:  inLabel,
		outputTag: outLabel,
	})

	return nil
}

func (p *pipeline) RegisterSeparator(fn SeparatorFunc, inLabel string, outputs []string) error {
	if fn == nil {
		return errors.New("nil separator")
	}

	if len(outputs) == 0 {
		return errors.New("empty outputs")
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.separators = append(p.separators, separatorTask{
		handler:    fn,
		inputTag:   inLabel,
		outputTags: outputs,
	})

	return nil
}

func (p *pipeline) RegisterMultiplexer(fn MultiplexerFunc, inputLabels []string, outputLabel string) error {
	if fn == nil {
		return errors.New("nil multiplexer")
	}

	if len(inputLabels) == 0 {
		return errors.New("empty inputs")
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.mults = append(p.mults, multiplexerTask{
		handler:   fn,
		inputTags: append([]string(nil), inputLabels...),
		outputTag: outputLabel,
	})

	return nil
}

func (p *pipeline) Send(input string, value string) error {
	ch, ok := p.inputMap[input]
	if !ok {
		return errors.New("chan not found")
	}

	select {
	case ch <- value:
		return nil
	default:
		return errors.New("pipeline closed")
	}
}

func (p *pipeline) Recv(out string) (string, error) {
	ch, ok := p.outputMap[out]
	if !ok {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return "", errors.New("pipeline closed")
	}

	return val, nil
}

func (p *pipeline) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, job := range p.decorators {
		item := job

		group.Go(func() error {
			in := p.channelIn(item.inputTag)
			out := p.channelOut(item.outputTag)

			err := item.handler(ctx, in, out)
			if err != nil {
				return fmt.Errorf("decorator failed: %w", err)
			}
			close(out)

			return nil
		})
	}

	for _, job := range p.separators {
		item := job

		group.Go(func() error {
			input := p.channelIn(item.inputTag)

			var outputs []chan string
			for _, name := range item.outputTags {
				outputs = append(outputs, p.channelOut(name))
			}

			err := item.handler(ctx, input, outputs)
			if err != nil {
				return fmt.Errorf("separator failed: %w", err)
			}

			for _, channel := range outputs {
				close(channel)
			}

			return nil
		})
	}

	for _, job := range p.mults {
		item := job

		group.Go(func() error {
			var inputs []chan string

			for _, name := range item.inputTags {
				inputs = append(inputs, p.channelIn(name))
			}

			output := p.channelOut(item.outputTag)

			err := item.handler(ctx, inputs, output)
			if err != nil {
				return fmt.Errorf("multiplexer failed: %w", err)
			}

			close(output)

			return nil
		})
	}

	return group.Wait()
}

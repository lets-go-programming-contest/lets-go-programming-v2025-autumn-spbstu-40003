package conveyer

import "context"

func (conv *Conveyer) EnsureChannel(name string) chan string {
	channel, exists := conv.channels[name]
	if !exists {
		channel = make(chan string, conv.bufSize)
		conv.channels[name] = channel
	}

	return channel
}

func (conv *Conveyer) CreateChannel(names ...string) {
	for _, name := range names {
		conv.EnsureChannel(name)
	}
}

func (conv *Conveyer) AddHandler(fn func(ctx context.Context) error) {
	conv.handlers = append(conv.handlers, fn)
}

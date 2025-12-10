package conveyer

import "context"

func (conv *Conveyer) EnsureChannel(name string) chan string {
	ch, ok := conv.channels[name]
	if !ok {
		ch = make(chan string, conv.bufSize)
		conv.channels[name] = ch
	}
	return ch
}

func (conv *Conveyer) CreateChannel(names ...string) {
	for _, name := range names {
		conv.EnsureChannel(name)
	}
}

func (conv *Conveyer) AddHandler(fn func(ctx context.Context) error) {
	conv.handlers = append(conv.handlers, fn)
}

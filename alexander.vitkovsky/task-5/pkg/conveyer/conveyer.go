package conveyer

import (
	"context"
	"sync"
)

// Conveyer struct + constructor + supportive structures and functions

type handlerConfig struct {
	kind     string
	function interface{}
	inputs   []string
	outputs  []string
}

type Conveyer struct {
	channels map[string]chan string
	handlers []handlerConfig
	size     int

	wg     sync.WaitGroup
	errCh  chan error
	cancel context.CancelFunc
	runCtx context.Context

	ready   chan struct{}
	started bool

	mutex sync.Mutex
}

func New(size int) *Conveyer {
	conv := new(Conveyer)

	conv.channels = make(map[string]chan string)
	conv.size = size
	conv.errCh = make(chan error, 1)
	conv.ready = make(chan struct{})

	return conv
}

func (conv *Conveyer) getOrCreateChannel(name string) chan string {
	conv.mutex.Lock()
	defer conv.mutex.Unlock()

	if ch, ok := conv.channels[name]; ok {
		return ch
	}
	ch := make(chan string, conv.size)
	conv.channels[name] = ch
	return ch
}

func (conv *Conveyer) resolveChannels(names []string) []chan string {
	result := make([]chan string, len(names))
	for index, name := range names {
		result[index] = conv.getOrCreateChannel(name)
	}
	return result
}

func (conv *Conveyer) closeAllChannels() {
	conv.mutex.Lock()
	defer conv.mutex.Unlock()

	for name, ch := range conv.channels {
		func() {
			defer func() { _ = recover() }()
			if ch != nil {
				close(ch)
			}
			conv.channels[name] = nil
		}()
	}
}

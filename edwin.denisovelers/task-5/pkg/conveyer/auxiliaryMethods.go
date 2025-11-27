package conveyer

import "context"

func (c *Conveyer) ensureChannel(name string) chan string {
	channel, ok := c.channels[name]
	if !ok {
		channel = make(chan string, c.bufSize)
		c.channels[name] = channel
	}

	return channel
}

func (c *Conveyer) createChannels(channelNames ...string) {
	for _, channelName := range channelNames {
		c.ensureChannel(channelName)
	}
}

func (c *Conveyer) closeChannels() {
	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) addHandler(fn func(context.Context) error) {
	c.handlers = append(c.handlers, fn)
}

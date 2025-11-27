package chansyncmap

import "sync"

type ChanSyncMap struct {
	size         int
	channelsByID sync.Map
}

func New(size int) *ChanSyncMap {
	return &ChanSyncMap{
		size:         size,
		channelsByID: sync.Map{},
	}
}

func (csm *ChanSyncMap) GetOrCreateChan(chanName string) chan string {
	var channel chan string
	if preChannel, ok := csm.channelsByID.Load(chanName); !ok {
		channel = make(chan string, csm.size)
		csm.channelsByID.Store(chanName, channel)
	} else {
		channel, ok = preChannel.(chan string)
		if !ok {
			return nil
		}
	}

	return channel
}

func (csm *ChanSyncMap) GetChan(chanName string) (chan string, bool) {
	var channel chan string

	if preChannel, ok := csm.channelsByID.Load(chanName); !ok {
		return nil, false
	} else {
		channel, ok = preChannel.(chan string)

		if !ok {
			return nil, false
		}
	}

	return channel, true
}

func (csm *ChanSyncMap) CloseAllChans() {
	csm.channelsByID.Range(func(key, value interface{}) bool {
		channel, ok := value.(chan string)
		if !ok {
			return false
		}

		close(channel)

		return true
	})
}

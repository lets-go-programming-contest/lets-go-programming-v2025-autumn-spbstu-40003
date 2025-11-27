package chansyncmap

import "sync"

type ChanSyncMap struct {
	size         int
	channelsById sync.Map
}

func New(size int) *ChanSyncMap {
	return &ChanSyncMap{
		size: size,
	}
}

func (csm *ChanSyncMap) GetOrCreateChan(chanName string) chan string {
	var channel chan string
	if preChannel, ok := csm.channelsById.Load(chanName); !ok {
		channel = make(chan string, csm.size)
		csm.channelsById.Store(chanName, channel)
	} else {
		channel = preChannel.(chan string)
	}

	return channel
}

func (csm *ChanSyncMap) GetChan(chanName string) (chan string, bool) {
	var channel chan string
	if preChannel, ok := csm.channelsById.Load(chanName); !ok {
		return nil, false
	} else {
		channel = preChannel.(chan string)
	}

	return channel, true
}

func (csm *ChanSyncMap) CloseAllChans() {
	csm.channelsById.Range(func(key, value interface{}) bool {
		close(value.(chan string))
		return true
	})
}

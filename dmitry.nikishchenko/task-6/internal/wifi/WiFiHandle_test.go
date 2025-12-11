package wifi_test

import (
	"fmt"

	wifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type WiFiHandle struct {
	mock.Mock
}

func (_m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()
	if len(ret) == 0 {
		panic("no return value specified for Interfaces")
	}

	var (
		r0 []*wifi.Interface
		r1 error
	)

	rf0, ok := ret.Get(0).(func() ([]*wifi.Interface, error))
	if ok {
		out0, out1 := rf0()
		return out0, fmt.Errorf("wrapped: %w", out1)
	}

	rf1, ok := ret.Get(0).(func() []*wifi.Interface)
	if ok {
		r0 = rf1()
	}

	if !ok && ret.Get(0) != nil {
		v, ok := ret.Get(0).([]*wifi.Interface)
		if !ok {
			return nil, fmt.Errorf("unexpected type for return value 0: %T", ret.Get(0))
		}
		r0 = v
	}

	rfErr, ok := ret.Get(1).(func() error)
	if ok {
		r1 = rfErr()
	} else {
		r1 = fmt.Errorf("wrapped: %w", ret.Error(1))
	}

	return r0, r1
}

type testInterface interface {
	mock.TestingT
	Cleanup(func())
}

func NewWiFiHandle(t testInterface) *WiFiHandle {
	h := &WiFiHandle{}
	h.Mock.Test(t)

	t.Cleanup(func() {
		h.AssertExpectations(t)
	})

	return h
}

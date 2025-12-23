package wifi_test

import (
	"errors"
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errTypeAssertionFailed = errors.New("type assertion failed")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface

	if v := args.Get(0); v != nil {
		cast, ok := v.([]*wifi.Interface)
		if !ok {
			return nil, errTypeAssertionFailed
		}

		ifaces = cast
	}

	if err := args.Error(1); err != nil {
		return ifaces, fmt.Errorf("mock interfaces error: %w", err)
	}

	return ifaces, nil
}

package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		err := args.Error(1)
		if err != nil {
			return nil, fmt.Errorf("wifi error: %w", err)
		}
		return nil, nil
	}
	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}
	return interfaces, args.Error(1)
}

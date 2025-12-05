package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		hwAddr, _ := net.ParseMAC("00:00:5e:00:53:01")
		interfaces := []*wifi.Interface{
			{ID: 1, Name: "wlan0", HardwareAddr: hwAddr},
		}

		mockHandle.On("Interfaces").Return(interfaces, nil)

		addrs, err := service.GetAddresses()

		assert.NoError(t, err)
		assert.Len(t, addrs, 1)
		assert.Equal(t, hwAddr, addrs[0])
		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		expectedErr := errors.New("network error")

		mockHandle.On("Interfaces").Return(nil, expectedErr)

		addrs, err := service.GetAddresses()

		assert.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		interfaces := []*wifi.Interface{
			{ID: 1, Name: "wlan0"},
			{ID: 2, Name: "wlan1"},
		}

		mockHandle.On("Interfaces").Return(interfaces, nil)

		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errors.New("fail"))

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		mockHandle.AssertExpectations(t)
	})
}

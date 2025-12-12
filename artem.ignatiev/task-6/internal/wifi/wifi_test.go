package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		hwAddr, _ := net.ParseMAC("00:11:22:33:44:55")
		mockInterfaces := []*wifi.Interface{
			{
				Index:        1,
				Name:         "wlan0",
				HardwareAddr: hwAddr,
			},
		}

		mockWiFi.On("Interfaces").Return(mockInterfaces, nil)

		addrs, err := service.GetAddresses()

		assert.NoError(t, err)
		assert.Len(t, addrs, 1)
		assert.Equal(t, hwAddr, addrs[0])

		mockWiFi.AssertExpectations(t)
	})

	t.Run("error getting interfaces", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		expectedErr := errors.New("hardware error")

		mockWiFi.On("Interfaces").Return(nil, expectedErr)

		addrs, err := service.GetAddresses()

		assert.Error(t, err)
		assert.Nil(t, addrs)
		assert.ErrorIs(t, err, expectedErr)
		assert.Contains(t, err.Error(), "getting interfaces")

		mockWiFi.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		mockInterfaces := []*wifi.Interface{
			{Name: "eth0"},
			{Name: "wlan0"},
		}

		mockWiFi.On("Interfaces").Return(mockInterfaces, nil)

		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Len(t, names, 2)
		assert.Equal(t, []string{"eth0", "wlan0"}, names)

		mockWiFi.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errors.New("fail"))

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)

		mockWiFi.AssertExpectations(t)
	})
}

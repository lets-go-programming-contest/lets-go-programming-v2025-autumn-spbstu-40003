package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifiservice "github.com/kryjkaqq/task-6/internal/wifi"
)

var errWifi = errors.New("wifi error")

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("successful get addresses", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := wifiservice.New(mockWiFi)

		mac1, _ := net.ParseMAC("00:11:22:33:44:55")
		mac2, _ := net.ParseMAC("AA:BB:CC:DD:EE:FF")

		interfaces := []*wifi.Interface{
			{HardwareAddr: mac1, Name: "wlan0"},
			{HardwareAddr: mac2, Name: "wlan1"},
		}

		mockWiFi.On("Interfaces").Return(interfaces, nil)

		addrs, err := service.GetAddresses()

		require.NoError(t, err)
		assert.Len(t, addrs, 2)
		assert.Equal(t, mac1, addrs[0])
		assert.Equal(t, mac2, addrs[1])
		mockWiFi.AssertExpectations(t)
	})

	t.Run("empty interfaces list", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := wifiservice.New(mockWiFi)

		mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

		addrs, err := service.GetAddresses()

		require.NoError(t, err)
		assert.Empty(t, addrs)
		mockWiFi.AssertExpectations(t)
	})

	t.Run("interfaces error", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := wifiservice.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errWifi)

		addrs, err := service.GetAddresses()

		require.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockWiFi.AssertExpectations(t)
	})
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	t.Run("successful get names", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := wifiservice.New(mockWiFi)

		mac1, _ := net.ParseMAC("00:11:22:33:44:55")
		mac2, _ := net.ParseMAC("AA:BB:CC:DD:EE:FF")

		interfaces := []*wifi.Interface{
			{HardwareAddr: mac1, Name: "wlan0"},
			{HardwareAddr: mac2, Name: "wlan1"},
		}

		mockWiFi.On("Interfaces").Return(interfaces, nil)

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Len(t, names, 2)
		assert.Equal(t, "wlan0", names[0])
		assert.Equal(t, "wlan1", names[1])
		mockWiFi.AssertExpectations(t)
	})

	t.Run("empty interfaces list", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := wifiservice.New(mockWiFi)

		mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Empty(t, names)
		mockWiFi.AssertExpectations(t)
	})

	t.Run("interfaces error", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := wifiservice.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errWifi)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockWiFi.AssertExpectations(t)
	})
}

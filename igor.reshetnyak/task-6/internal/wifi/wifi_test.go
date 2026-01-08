package wifi_test

import (
	"errors"
	"net"
	"testing"

	wifiService "github.com/ReshetnyakIgor/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

func createMACAddress(t *testing.T, addr string) net.HardwareAddr {
	t.Helper()
	mac, err := net.ParseMAC(addr)
	require.NoError(t, err)
	return mac
}

func TestNewWiFiService(t *testing.T) {
	t.Parallel()

	mock := NewWiFiHandle(t)
	service := wifiService.New(mock)
	assert.NotNil(t, service)
	assert.Equal(t, mock, service.WiFi)
}

func TestGetAddresses_Success(t *testing.T) {
	t.Parallel()

	mock := NewWiFiHandle(t)
	service := wifiService.New(mock)

	interfaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: createMACAddress(t, "11:22:33:44:55:66")},
		{Name: "wlan1", HardwareAddr: createMACAddress(t, "aa:bb:cc:dd:ee:ff")},
	}
	mock.On("Interfaces").Return(interfaces, nil)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	assert.Equal(t, []net.HardwareAddr{
		createMACAddress(t, "11:22:33:44:55:66"),
		createMACAddress(t, "aa:bb:cc:dd:ee:ff"),
	}, addrs)
}

func TestGetAddresses_Empty(t *testing.T) {
	t.Parallel()

	mock := NewWiFiHandle(t)
	service := wifiService.New(mock)

	mock.On("Interfaces").Return([]*wifi.Interface{}, nil)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	assert.Empty(t, addrs)
}

func TestGetAddresses_Error(t *testing.T) {
	t.Parallel()

	mock := NewWiFiHandle(t)
	service := wifiService.New(mock)

	mock.On("Interfaces").Return(nil, errTest)

	addrs, err := service.GetAddresses()
	require.ErrorContains(t, err, "getting interfaces")
	assert.Nil(t, addrs)
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mock := NewWiFiHandle(t)
	service := wifiService.New(mock)

	interfaces := []*wifi.Interface{
		{Name: "wifi0", HardwareAddr: createMACAddress(t, "00:11:22:33:44:55")},
		{Name: "wifi1", HardwareAddr: createMACAddress(t, "66:77:88:99:aa:bb")},
	}
	mock.On("Interfaces").Return(interfaces, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"wifi0", "wifi1"}, names)
}

func TestGetNames_NoInterfaces(t *testing.T) {
	t.Parallel()

	mock := NewWiFiHandle(t)
	service := wifiService.New(mock)

	mock.On("Interfaces").Return([]*wifi.Interface{}, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestGetNames_Failure(t *testing.T) {
	t.Parallel()

	mock := NewWiFiHandle(t)
	service := wifiService.New(mock)

	mock.On("Interfaces").Return(nil, errTest)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "getting interfaces")
	assert.Nil(t, names)
}

func TestGetNames_WithNilAddress(t *testing.T) {
	t.Parallel()

	mock := NewWiFiHandle(t)
	service := wifiService.New(mock)

	interfaces := []*wifi.Interface{
		{Name: "eth0", HardwareAddr: nil},
		{Name: "wlan0", HardwareAddr: createMACAddress(t, "11:22:33:44:55:66")},
	}
	mock.On("Interfaces").Return(interfaces, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"eth0", "wlan0"}, names)
}

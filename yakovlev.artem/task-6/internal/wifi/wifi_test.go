package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

func TestWiFiService_GetAddresses_OK(t *testing.T) {
	m := NewMockWiFiHandle(t)

	ifaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}},
		{Name: "wlan1", HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
	}

	m.On("Interfaces").Return(ifaces, nil).Once()

	service := New(m)
	got, err := service.GetAddresses()
	require.NoError(t, err)

	require.Equal(t, []net.HardwareAddr{
		net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
	}, got)
}

func TestWiFiService_GetAddresses_InterfacesError(t *testing.T) {
	m := NewMockWiFiHandle(t)

	iErr := errors.New("boom")
	m.On("Interfaces").Return(([]*wifi.Interface)(nil), iErr).Once()

	service := New(m)
	got, err := service.GetAddresses()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces:")
	require.ErrorIs(t, err, iErr)
}

func TestWiFiService_GetNames_OK(t *testing.T) {
	m := NewMockWiFiHandle(t)

	ifaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}},
		{Name: "wlan1", HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
	}

	m.On("Interfaces").Return(ifaces, nil).Once()

	service := New(m)
	got, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "wlan1"}, got)
}

func TestWiFiService_GetNames_InterfacesError(t *testing.T) {
	m := NewMockWiFiHandle(t)

	iErr := errors.New("boom")
	m.On("Interfaces").Return(([]*wifi.Interface)(nil), iErr).Once()

	service := New(m)
	got, err := service.GetNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces:")
	require.ErrorIs(t, err, iErr)
}

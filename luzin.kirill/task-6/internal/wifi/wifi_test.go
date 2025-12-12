package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	myWifi "github.com/KiRy6A/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var errExpected = errors.New("expected error")

type rowTestSysInfo struct {
	names       []string
	addrs       []string
	expectedErr error
}

var testTable = []rowTestSysInfo{
	{
		names: []string{"GuestSPBPU", "PepeRun"},
		addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
	},
	{
		names:       []string{"GuestSPBPU", "Public"},
		expectedErr: errExpected,
	},
	{
		addrs:       []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
		expectedErr: errExpected,
	},
}

func TestNew(t *testing.T) {
	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.New(mockWifi)

	require.NotNil(t, mockWifi, "WiFiHandle should not be nil")
	require.NotNil(t, wifiService, "WiFiService should not be nil")
	require.Equal(t, mockWifi, wifiService.WiFi, "WiFiHandle and WiFiService.WiFiHandle should be equal")
}

func TestGetNames(t *testing.T) {
	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.WiFiService{mockWifi}

	for i, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfacesNames(row), row.expectedErr)
		names, err := wifiService.GetNames()
		if row.expectedErr != nil {
			require.ErrorIs(t, err, row.expectedErr, "row: %d, expected error: %w, actual error: %w", i, row.expectedErr, err)
			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, row.names, names,
			"row: %d, expected addrs: %s, actual addrs: %s", i,
			parseMACs(row.addrs), names)
	}
}

func TestGetAddresses(t *testing.T) {
	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.WiFiService{mockWifi}

	for i, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfacesAddrs(row), row.expectedErr)
		actualAddrs, err := wifiService.GetAddresses()
		if row.expectedErr != nil {
			require.ErrorIs(t, err, row.expectedErr, "row: %d, expected error: %w, actual error: %w", i, row.expectedErr, err)
			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, parseMACs(row.addrs), actualAddrs,
			"row: %d, expected addrs: %s, actual addrs: %s", i,
			parseMACs(row.addrs), actualAddrs)
	}
}

func mockIfacesAddrs(row rowTestSysInfo) []*wifi.Interface {
	var interfaces []*wifi.Interface

	for i, addrStr := range row.addrs {
		hwAddr := parseMAC(addrStr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("eth%d", i+1),
			HardwareAddr: hwAddr,
			PHY:          1,
			Device:       1,
			Type:         wifi.InterfaceTypeAPVLAN,
			Frequency:    0,
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces
}

func mockIfacesNames(row rowTestSysInfo) []*wifi.Interface {
	var interfaces []*wifi.Interface

	for i, name := range row.names {
		hwAddr := parseMAC(fmt.Sprintf("00:11:22:33:44:%02x", i+1))
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         name,
			HardwareAddr: hwAddr,
			PHY:          1,
			Device:       1,
			Type:         wifi.InterfaceTypeAPVLAN,
			Frequency:    0,
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces
}

func parseMACs(macStr []string) []net.HardwareAddr {
	var addrs []net.HardwareAddr

	for _, addr := range macStr {
		addrs = append(addrs, parseMAC(addr))
	}

	return addrs
}

func parseMAC(macStr string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(macStr)
	if err != nil {
		return nil
	}
	return hwAddr
}

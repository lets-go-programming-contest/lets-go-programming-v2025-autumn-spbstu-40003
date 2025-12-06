package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	myWifi "task-6/internal/wifi"
)

var ErrExpected = errors.New("expected error")

type rowTestSysInfo struct {
	addrs       []string
	names       []string
	errExpected error
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	testTable := getTestCases()
	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.WiFiService{WiFi: mockWifi}

	for _, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(helperMockIfaces(t, row.addrs), row.errExpected).Once()

		actualAddrs, err := wifiService.GetAddresses()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected)
			continue
		}

		require.NoError(t, err)
		require.Equal(t, helperParseMACs(t, row.addrs), actualAddrs)
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	testTable := getTestCases()
	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.WiFiService{WiFi: mockWifi}

	for _, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(helperMockIfaces(t, row.addrs), row.errExpected).Once()

		actualNames, err := wifiService.GetNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected)
			continue
		}

		require.NoError(t, err)
		require.Equal(t, row.names, actualNames)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.New(mockWifi)

	require.Equal(t, mockWifi, wifiService.WiFi)
}

func helperMockIfaces(t *testing.T, addrs []string) []*wifi.Interface {
	t.Helper()

	interfaces := make([]*wifi.Interface, 0)
	validIndex := 1

	for _, addrStr := range addrs {
		hwAddr := parseMAC(addrStr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        validIndex,
			Name:         fmt.Sprintf("eth%d", validIndex),
			HardwareAddr: hwAddr,
			Type:         wifi.InterfaceTypeAPVLAN,
		}
		interfaces = append(interfaces, iface)
		validIndex++
	}

	return interfaces
}

func helperParseMACs(t *testing.T, macStr []string) []net.HardwareAddr {
	t.Helper()

	addrs := make([]net.HardwareAddr, 0)
	for _, addr := range macStr {
		if hwAddr := parseMAC(addr); hwAddr != nil {
			addrs = append(addrs, hwAddr)
		}
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

func getTestCases() []rowTestSysInfo {
	return []rowTestSysInfo{
		{
			addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			names: []string{"eth1", "eth2"},
		},
		{
			addrs:       nil,
			names:       nil,
			errExpected: ErrExpected,
		},
		{
			addrs: []string{},
			names: []string{},
		},
		{
			addrs: []string{"00:11:22:33:44:55", "invalid-mac-addr", "bb:bb:cc:dd:ee:ff"},
			names: []string{"eth1", "eth2"},
		},
	}
}

package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	myWifi "github.com/wedwincode/task-6/internal/wifi"
)

var errExpected = errors.New("ExpectedError")

type rowTest struct {
	addrs       []string
	names       []string
	errExpected error
}

func getTestTable() []rowTest {
	return []rowTest{
		{
			addrs:       []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			names:       []string{"wlan0", "wlan1"},
			errExpected: nil,
		},
		{
			addrs:       []string{},
			names:       []string{},
			errExpected: errExpected,
		},
		{
			addrs:       []string{},
			names:       []string{},
			errExpected: nil,
		},
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.New(mockWifi)

	require.NotNil(t, wifiService.WiFi, "WiFi handle should not be nil")
	require.Equal(t, mockWifi, wifiService.WiFi, "WiFi handle should be set correctly")
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.New(mockWifi)

	for idx, row := range getTestTable() {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfacesFromAddrs(row.addrs), row.errExpected)

		actualAddrs, err := wifiService.GetAddresses()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected,
				"row: %d, expected error: %w, actual error: %w", idx, row.errExpected, err)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", idx)
		require.Equal(t, parseMACs(row.addrs), actualAddrs,
			"row: %d, expected addrs: %s, actual addrs: %s", idx, parseMACs(row.addrs), actualAddrs)
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.New(mockWifi)

	for idx, row := range getTestTable() {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfacesFromNames(row.names), row.errExpected)

		actualNames, err := wifiService.GetNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected,
				"row: %d, expected error: %w, actual error: %w", idx, row.errExpected, err)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", idx)
		require.Equal(t, row.names, actualNames,
			"row: %d, expected addrs: %s, actual addrs: %s", idx, row.names, actualNames)
	}
}

func mockIfacesFromAddrs(addrs []string) []*wifi.Interface {
	interfaces := make([]*wifi.Interface, 0, len(addrs))

	for idx, addrStr := range addrs {
		hwAddr := parseMAC(addrStr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        idx + 1,
			Name:         fmt.Sprintf("eth%d", idx+1),
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

func mockIfacesFromNames(names []string) []*wifi.Interface {
	interfaces := make([]*wifi.Interface, 0, len(names))

	for idx, name := range names {
		hwAddr := parseMAC(fmt.Sprintf("00:11:22:33:44:%02x", idx))

		iface := &wifi.Interface{
			Index:        idx + 1,
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
	if len(macStr) == 0 {
		return []net.HardwareAddr{}
	}

	addrs := make([]net.HardwareAddr, 0, len(macStr))

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

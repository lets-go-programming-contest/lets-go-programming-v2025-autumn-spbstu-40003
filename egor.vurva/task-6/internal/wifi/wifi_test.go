package wifi_test

//go:generate mockery --name=WiFiHandle --output=. --outpkg=wifi_test --case=underscore

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	wifipkg "github.com/Vurvaa/task-6/internal/wifi"
)

var ErrInterfaces = errors.New("interfaces error")

type testCase struct {
	interfaces []*wifi.Interface
	err        error

	expectedNames     []string
	expectedAddresses []net.HardwareAddr
	wantError         bool
}

func TestNew(t *testing.T) {
	t.Parallel()

	wifiMock := NewWiFiHandle(t)

	service := wifipkg.New(wifiMock)
	require.Equal(t, wifiMock, service.WiFi)
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	testTable := []testCase{
		{
			interfaces: []*wifi.Interface{
				{HardwareAddr: net.HardwareAddr{0x00, 0x22, 0x33, 0x55}},
				{HardwareAddr: net.HardwareAddr{0xaa, 0xcc, 0xdd, 0xff}},
			},
			expectedAddresses: []net.HardwareAddr{
				{0x00, 0x22, 0x33, 0x55},
				{0xaa, 0xcc, 0xdd, 0xff},
			},
		},
		{
			err:       ErrInterfaces,
			wantError: true,
		},
		{
			interfaces:        []*wifi.Interface{},
			expectedAddresses: []net.HardwareAddr{},
		},
	}

	runTestTable(t, testTable, func(service wifipkg.WiFiService) ([]string, []net.HardwareAddr, error) {
		addresses, err := service.GetAddresses()
		if err != nil {
			return nil, nil, fmt.Errorf("GetAddresses: %w", err)
		}

		return nil, addresses, nil
	})
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	testTable := []testCase{
		{
			interfaces: []*wifi.Interface{
				{Name: "wlan0"},
				{Name: "wlan1"},
			},
			expectedNames: []string{"wlan0", "wlan1"},
		},
		{
			err:       ErrInterfaces,
			wantError: true,
		},
		{
			interfaces:    []*wifi.Interface{},
			expectedNames: []string{},
		},
	}

	runTestTable(t, testTable, func(service wifipkg.WiFiService) ([]string, []net.HardwareAddr, error) {
		names, err := service.GetNames()
		if err != nil {
			return nil, nil, fmt.Errorf("GetNames: %w", err)
		}

		return names, nil, nil
	})
}

type callFunc func(wifipkg.WiFiService) ([]string, []net.HardwareAddr, error)

func runTestTable(t *testing.T, testTable []testCase, call callFunc) {
	t.Helper()

	for i, test := range testTable {
		wifiMock := NewWiFiHandle(t)

		if test.err != nil {
			wifiMock.On("Interfaces").Return(nil, test.err)
		} else {
			wifiMock.On("Interfaces").Return(test.interfaces, nil)
		}

		service := wifipkg.WiFiService{WiFi: wifiMock}

		names, addresses, err := call(service)

		if test.wantError {
			require.Error(t, err, "case %d", i)
			require.ErrorContains(t, err, "getting interfaces", "case %d", i)
			require.Nil(t, names, "case %d", i)
			require.Nil(t, addresses, "case %d", i)

			continue
		}

		require.NoError(t, err, "case %d", i)

		if test.expectedNames != nil {
			require.Equal(t, test.expectedNames, names, "case %d", i)
		}

		if test.expectedAddresses != nil {
			require.Equal(t, test.expectedAddresses, addresses, "case %d", i)
		}
	}
}

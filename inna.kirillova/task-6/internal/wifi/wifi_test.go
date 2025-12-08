package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWiFi struct {
	ifaces []*wifi.Interface
	err    error
}

func (m *mockWiFi) Interfaces() ([]*wifi.Interface, error) {
	return m.ifaces, m.err
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFi{}
	service := New(mockHandle)

	assert.NotNil(t, service)
	assert.Equal(t, mockHandle, service.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	mac1 := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	mac2 := net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	mac3 := net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	mac4 := net.HardwareAddr{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}

	tests := []struct {
		name        string
		mock        *mockWiFi
		expectAddrs []net.HardwareAddr
		expectErr   bool
		errContains string
	}{
		{
			name: "success two interfaces",
			mock: &mockWiFi{
				ifaces: []*wifi.Interface{
					{HardwareAddr: mac1},
					{HardwareAddr: mac2},
				},
			},
			expectAddrs: []net.HardwareAddr{mac1, mac2},
		},
		{
			name: "success four interfaces",
			mock: &mockWiFi{
				ifaces: []*wifi.Interface{
					{HardwareAddr: mac1},
					{HardwareAddr: mac2},
					{HardwareAddr: mac3},
					{HardwareAddr: mac4},
				},
			},
			expectAddrs: []net.HardwareAddr{mac1, mac2, mac3, mac4},
		},
		{
			name: "empty interfaces",
			mock: &mockWiFi{
				ifaces: []*wifi.Interface{},
			},
			expectAddrs: []net.HardwareAddr{},
		},
		{
			name: "interfaces error",
			mock: &mockWiFi{
				err: errors.New("system error"),
			},
			expectErr:   true,
			errContains: "getting interfaces:",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := New(tc.mock)
			addrs, err := svc.GetAddresses()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
				assert.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectAddrs, addrs)
			}
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	mac1 := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	mac2 := net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

	tests := []struct {
		name        string
		mock        *mockWiFi
		expectNames []string
		expectErr   bool
		errContains string
	}{
		{
			name: "success two names",
			mock: &mockWiFi{
				ifaces: []*wifi.Interface{
					{Name: "wlan0", HardwareAddr: mac1},
					{Name: "eth0", HardwareAddr: mac2},
				},
			},
			expectNames: []string{"wlan0", "eth0"},
		},
		{
			name: "empty list",
			mock: &mockWiFi{
				ifaces: []*wifi.Interface{},
			},
			expectNames: []string{},
		},
		{
			name: "interface with empty name",
			mock: &mockWiFi{
				ifaces: []*wifi.Interface{
					{Name: ""},
					{Name: "wlan1"},
				},
			},
			expectNames: []string{"", "wlan1"},
		},
		{
			name: "error from Interfaces",
			mock: &mockWiFi{
				err: errors.New("permission error"),
			},
			expectErr:   true,
			errContains: "getting interfaces:",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := New(tc.mock)
			names, err := svc.GetNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectNames, names)
			}
		})
	}
}


package wifi_test

import (
	"errors"
	"net"
	"testing"

	wifipkg "github.com/KrrMaxim/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errInterfaces = errors.New("interfaces error")

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandle := &MockWiFiHandle{}
	service := wifipkg.New(mockHandle)

	require.NotNil(t, service)
	assert.Equal(t, mockHandle, service.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	mac1 := net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	mac2 := net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

	tests := []struct {
		name      string
		ifaces    []*wifi.Interface
		err       error
		want      []net.HardwareAddr
		wantError bool
	}{
		{
			name: "success multiple interfaces",
			ifaces: []*wifi.Interface{
				{HardwareAddr: mac1},
				{HardwareAddr: mac2},
			},
			want: []net.HardwareAddr{mac1, mac2},
		},
		{
			name:   "success empty",
			ifaces: []*wifi.Interface{},
			want:   []net.HardwareAddr{},
		},
		{
			name:      "error",
			err:       errInterfaces,
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &MockWiFiHandle{}
			mockHandle.
				On("Interfaces").
				Return(tc.ifaces, tc.err)

			service := wifipkg.New(mockHandle)
			addrs, err := service.GetAddresses()

			if tc.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
				assert.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, addrs)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		ifaces    []*wifi.Interface
		err       error
		want      []string
		wantError bool
	}{
		{
			name: "success multiple names",
			ifaces: []*wifi.Interface{
				{Name: "wlan0"},
				{Name: "eth0"},
			},
			want: []string{"wlan0", "eth0"},
		},
		{
			name:   "success empty",
			ifaces: []*wifi.Interface{},
			want:   []string{},
		},
		{
			name:      "error",
			err:       errInterfaces,
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &MockWiFiHandle{}
			mockHandle.
				On("Interfaces").
				Return(tc.ifaces, tc.err)

			service := wifipkg.New(mockHandle)
			names, err := service.GetNames()

			if tc.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, names)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

package wifi_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ArtttNik/task-6/internal/wifi"

	mdlayherwifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mockHandle := &wifi.MockWiFiHandle{}
	s := wifi.New(mockHandle)
	assert.Equal(t, mockHandle, s.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mock *wifi.MockWiFiHandle)
		want    []net.HardwareAddr
		wantErr string
	}{
		{
			name: "success multiple interfaces",
			setup: func(mock *wifi.MockWiFiHandle) {
				ifaces := []*mdlayherwifi.Interface{
					{HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}},
					{HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
				}
				mock.On("Interfaces").Return(ifaces, nil)
			},
			want: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			name: "success empty",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface{}, nil)
			},
			want: []net.HardwareAddr{},
		},
		{
			name: "error",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return(([]*mdlayherwifi.Interface)(nil), fmt.Errorf("interfaces error"))
			},
			wantErr: "getting interfaces: interfaces error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandle := &wifi.MockWiFiHandle{}
			tt.setup(mockHandle)

			s := wifi.New(mockHandle)
			got, err := s.GetAddresses()

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mock *wifi.MockWiFiHandle)
		want    []string
		wantErr string
	}{
		{
			name: "success multiple interfaces",
			setup: func(mock *wifi.MockWiFiHandle) {
				ifaces := []*mdlayherwifi.Interface{
					{Name: "eth0"},
					{Name: "wlan0"},
				}
				mock.On("Interfaces").Return(ifaces, nil)
			},
			want: []string{"eth0", "wlan0"},
		},
		{
			name: "success empty",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface{}, nil)
			},
			want: []string{},
		},
		{
			name: "error",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return(([]*mdlayherwifi.Interface)(nil), fmt.Errorf("interfaces error"))
			},
			wantErr: "getting interfaces: interfaces error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandle := &wifi.MockWiFiHandle{}
			tt.setup(mockHandle)

			s := wifi.New(mockHandle)
			got, err := s.GetNames()

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

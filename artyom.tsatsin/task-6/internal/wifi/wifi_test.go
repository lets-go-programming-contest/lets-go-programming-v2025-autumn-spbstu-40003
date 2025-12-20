package wifi

import (
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/mdlayher/wifi"
)

func TestGetAddresses_OK(t *testing.T) {
	mock := &WiFiHandleMock{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22}},
				{HardwareAddr: net.HardwareAddr{0xAA, 0xBB, 0xCC}},
			}, nil
		},
	}

	service := New(mock)

	got, err := service.GetAddresses()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []net.HardwareAddr{
		{0x00, 0x11, 0x22},
		{0xAA, 0xBB, 0xCC},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetAddresses_Error(t *testing.T) {
	mock := &WiFiHandleMock{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return nil, errors.New("fail")
		},
	}

	service := New(mock)

	_, err := service.GetAddresses()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetNames_OK(t *testing.T) {
	mock := &WiFiHandleMock{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{Name: "wlan0"},
				{Name: "wlan1"},
			}, nil
		},
	}

	service := New(mock)

	got, err := service.GetNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"wlan0", "wlan1"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetNames_Error(t *testing.T) {
	mock := &WiFiHandleMock{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return nil, errors.New("fail")
		},
	}

	service := New(mock)

	_, err := service.GetNames()
	if err == nil {
		t.Fatal("expected error")
	}
}

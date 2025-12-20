package wifi

import "github.com/mdlayher/wifi"

type WiFiHandleMock struct {
	InterfacesFunc func() ([]*wifi.Interface, error)
}

func (m *WiFiHandleMock) Interfaces() ([]*wifi.Interface, error) {
	return m.InterfacesFunc()
}

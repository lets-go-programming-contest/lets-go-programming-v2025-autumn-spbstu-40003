package wifi_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/gituser549/task-6/internal/util"
	MyWifi "github.com/gituser549/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
)

type EthIfaceInfo struct {
	MacAddr string
	Name    string
}

type TestEthIfaceUnit struct {
	ifaces      []EthIfaceInfo
	errExpected string
}

func TestGetAddresses(t *testing.T) {
	TestWifiServiceTable := []TestEthIfaceUnit{
		{
			ifaces:      nil,
			errExpected: "interfaces error",
		},
		{
			ifaces: []EthIfaceInfo{
				{MacAddr: "AA:BB:CC:DD:EE:FF", Name: "correct1"},
				{MacAddr: "00:1A:2B:3C:4D:5E", Name: "correct2"},
			},
		},
	}

	t.Parallel()

	for numTestCase, testCase := range TestWifiServiceTable {
		t.Run(fmt.Sprintf("%s #%d", t.Name(), numTestCase), func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewMockWifiHandle(t)

			wifiService := MyWifi.New(mockWiFi)

			mockWiFi.On("Interfaces").Unset()
			mockWiFi.On("Interfaces").
				Return(createIfacesFromTestData(testCase.ifaces),
					util.MakeError(testCase.errExpected)).
				Once()

			result, err := wifiService.GetAddresses()

			if !util.IsEmpty(testCase.errExpected) {
				util.AssertError(t, result, err, testCase.errExpected)

				return
			}

			util.AssertNoError(t, extractAddressesFromIfaces(testCase.ifaces), result, err)

			mockWiFi.AssertExpectations(t)
		})
	}
}

func TestGetNames(t *testing.T) {
	TestWifiServiceTable := []TestEthIfaceUnit{
		{
			ifaces:      nil,
			errExpected: "interfaces error",
		},
		{
			ifaces: []EthIfaceInfo{
				{MacAddr: "AA:BB:CC:DD:EE:FF", Name: "correct1"},
				{MacAddr: "00:1A:2B:3C:4D:5E", Name: "correct2"},
			},
		},
	}

	t.Parallel()

	for numTestCase, testCase := range TestWifiServiceTable {
		t.Run(fmt.Sprintf("%s #%d", t.Name(), numTestCase), func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewMockWifiHandle(t)

			wifiService := MyWifi.New(mockWiFi)

			mockWiFi.On("Interfaces").Unset()
			mockWiFi.On("Interfaces").
				Return(createIfacesFromTestData(testCase.ifaces),
					util.MakeError(testCase.errExpected)).
				Once()

			result, err := wifiService.GetNames()

			if !util.IsEmpty(testCase.errExpected) {
				util.AssertError(t, result, err, testCase.errExpected)

				return
			}

			util.AssertNoError(t, extractNamesFromIfaces(testCase.ifaces), result, err)

			mockWiFi.AssertExpectations(t)
		})
	}
}

func createIfacesFromTestData(ifaces []EthIfaceInfo) []*wifi.Interface {
	ethIfaces := make([]*wifi.Interface, 0, len(ifaces))

	for _, iface := range ifaces {
		hardwareAddr := returnAddressIfCorrect(iface.MacAddr)
		if hardwareAddr == nil {
			continue
		}

		ethIfaces = append(ethIfaces, &wifi.Interface{
			Name:         iface.Name,
			HardwareAddr: hardwareAddr,
		})
	}

	return ethIfaces
}

func extractAddressesFromIfaces(ifaces []EthIfaceInfo) []net.HardwareAddr {
	hardwareAddrs := make([]net.HardwareAddr, 0, len(ifaces))

	for _, iface := range ifaces {
		hardwareAddrs = append(hardwareAddrs, returnAddressIfCorrect(iface.MacAddr))
	}

	return hardwareAddrs
}

func returnAddressIfCorrect(address string) net.HardwareAddr {
	hardwareAddr, err := net.ParseMAC(address)
	if err != nil {
		return nil
	}

	return hardwareAddr
}

func extractNamesFromIfaces(ifaces []EthIfaceInfo) []string {
	names := make([]string, 0, len(ifaces))

	for _, iface := range ifaces {
		if returnAddressIfCorrect(iface.MacAddr) == nil {
			continue
		}

		names = append(names, iface.Name)
	}

	return names
}

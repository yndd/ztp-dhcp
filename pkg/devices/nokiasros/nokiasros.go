package nokiasros

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/devices"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

const DeviceName = "NokiaSROS"

type NokiaSros struct {
	backend backend.DhcpBackend
}

func (srl *NokiaSros) AdjustOffer(req *dhcpv4.DHCPv4, resp *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {

}

func (srl *NokiaSros) AdjustAck(req *dhcpv4.DHCPv4, resp *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {

}

func (srl *NokiaSros) SetBackend(backend backend.DhcpBackend) {
	srl.backend = backend
}

func init() {
	devices.DeviceManager.RegisterDevice(DeviceName, &NokiaSros{})
}

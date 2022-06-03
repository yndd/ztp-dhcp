package nokiasrl

import (
	"log"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/devices"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

var managed_models = []string{"NokiaSRL"}

type NokiaSrl struct {
	backend backend.DhcpBackend
}

func (srl *NokiaSrl) AdjustReply(req *dhcpv4.DHCPv4, reply *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {
	// set Option66
	reply.Options.Update(dhcpv4.OptTFTPServerName(devinfo.Option66))
	// set Option67
	reply.Options.Update(dhcpv4.OptBootFileName(devinfo.Option67))
}

func (srl *NokiaSrl) SetBackend(backend backend.DhcpBackend) {
	srl.backend = backend
}

func init() {
	err := devices.GetDeviceManagerRegistrator().RegisterDevice(managed_models, &NokiaSrl{})
	if err != nil {
		log.Fatal(err)
	}
}

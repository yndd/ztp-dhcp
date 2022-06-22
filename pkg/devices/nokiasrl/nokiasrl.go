package nokiasrl

import (
	"log"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/devices"
	"github.com/yndd/ztp-dhcp/pkg/structs"
	websstructs "github.com/yndd/ztp-webserver/pkg/structs"
)

var managed_models = []string{"SRLinux"}

type NokiaSrl struct {
	backend backend.DhcpBackend
}

func (srl *NokiaSrl) AdjustReply(req *dhcpv4.DHCPv4, reply *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {
	// set Option66
	reply.Options.Update(dhcpv4.OptTFTPServerName(devinfo.Option66))
	// set Option67

	reply.Options.Update(dhcpv4.OptBootFileName(websstructs.NewUrlParams(string(devinfo.VendorType), devinfo.Platform, websstructs.Script).GetUrlRelative()))
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

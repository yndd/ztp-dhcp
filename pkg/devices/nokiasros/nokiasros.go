package nokiasros

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/devices"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

var managed_models = []string{"NokiaSROS"}

type NokiaSros struct {
	backend backend.DhcpBackend
}

func (sros *NokiaSros) AdjustReply(req *dhcpv4.DHCPv4, reply *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {
	// set Option66
	//reply.Options.Update(dhcpv4.OptTFTPServerName(devinfo.Option66))
	// set Option67
	//reply.Options.Update(dhcpv4.OptBootFileName(devinfo.Option67))
}

func (sros *NokiaSros) SetBackend(backend backend.DhcpBackend) {
	sros.backend = backend
}

func init() {
	err := devices.GetDeviceManagerRegistrator().RegisterDevice(managed_models, &NokiaSros{})
	if err != nil {
		log.Fatal(err)
	}
}

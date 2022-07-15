package nokiasros

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/devices"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

var managed_models = []string{"SROS"}

type NokiaSros struct {
	backend backend.ZtpBackend
}

// AdjustReply takes the general reply and adjusts the device specific values
func (sros *NokiaSros) AdjustReply(req *dhcpv4.DHCPv4, reply *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {
	// set Option66
	//reply.Options.Update(dhcpv4.OptTFTPServerName(devinfo.Option66))
	// set Option67
	//reply.Options.Update(dhcpv4.OptBootFileName(devinfo.Option67))
}

// SetBackend used for late binding of the backend
func (sros *NokiaSros) SetBackend(backend backend.ZtpBackend) {
	sros.backend = backend
}

func init() {
	err := devices.GetDeviceManagerRegistrator().RegisterDevice(managed_models, &NokiaSros{})
	if err != nil {
		log.Fatal(err)
	}
}

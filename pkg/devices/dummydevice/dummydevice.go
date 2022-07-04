package dummydevice

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/devices"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

var managed_models = []string{"DummyDevice"}

type DummyDevice struct {
	backend backend.ZtpBackend
}

func (dd *DummyDevice) AdjustReply(req *dhcpv4.DHCPv4, reply *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {
	// set Option66
	//reply.Options.Update(dhcpv4.OptTFTPServerName(devinfo.Option66))
	// set Option67
	//reply.Options.Update(dhcpv4.OptBootFileName(devinfo.Option67))
}

func (dd *DummyDevice) SetBackend(backend backend.ZtpBackend) {
	dd.backend = backend
}

func init() {
	err := devices.GetDeviceManagerRegistrator().RegisterDevice(managed_models, &DummyDevice{})
	if err != nil {
		log.Fatal(err)
	}
}

package devices

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

// Device the interface that all vendor specific devices need to support
type Device interface {
	AdjustReply(req, resp *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation)
	// SetBackend is used to inject the Backend after self
	// registration of the Device with the DeviceManager
	SetBackend(backend backend.DhcpBackend)
}

package devices

import (
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
)

// DeviceManager is the dispatcher and registry for all the vendor specific devices ztp implementations
var DeviceManager DeviceManagerInterf

// DeviceManagerInterf the general DeviceManager Interface
type DeviceManagerInterf interface {
	RegisterDevice(string, Device)
}

// DeviceManagerImpl the actual implementation of the DeviceManager
type DeviceManagerImpl struct {
	devices map[string]Device
	backend backend.DhcpBackend
}

// AddDevice add Device specific implementations to the DeviceManager
func (dm *DeviceManagerImpl) RegisterDevice(devicetype string, device Device) {
	log.Infof("Adding %s to Devicemanager", devicetype)
	dm.devices[devicetype] = device
}

func (dm *DeviceManagerImpl) SetBackend(backend backend.DhcpBackend) {
	dm.backend = backend
}

// Initialize the DeviceManager
func init() {
	DeviceManager = &DeviceManagerImpl{
		devices: map[string]Device{},
	}
}

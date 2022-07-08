package devices

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
)

// DeviceManager is the dispatcher and registry for all the vendor specific devices ztp implementations
// Devicemanager is a singleton
var deviceManager *DeviceManagerImpl

// DeviceManagerImpl the actual implementation of the DeviceManager
type DeviceManagerImpl struct {
	devices map[string]Device
	backend backend.ZtpBackend
}

// GetDeviceManagerHandler retrieve a reference to
// the GetDeviceManagerHandler interface of this DeviceManagerImpl
func GetDeviceManagerHandler() DeviceManagerHandler {
	return newDeviceManager()
}

// GetDeviceManagerRegistrator retrieve a reference to
// the DeviceManagerRegistrator interface of this DeviceManagerImpl
func GetDeviceManagerRegistrator() DeviceManagerRegistrator {
	return newDeviceManager()
}

// newDeviceManager returns the singleton instance of the DeviceManager if
// uninitialized, it will construct the instance
func newDeviceManager() *DeviceManagerImpl {
	if deviceManager == nil {
		deviceManager = &DeviceManagerImpl{
			devices: map[string]Device{},
		}
	}
	return deviceManager
}

// AddDevice Device specific implementations to the DeviceManager
// the device thereby specifies which models it is responsible for
func (dm *DeviceManagerImpl) RegisterDevice(managed_models []string, device Device) error {
	for _, entry := range managed_models {
		// check if the specified model name was already registered
		if _, exists := dm.devices[entry]; exists {
			return fmt.Errorf("Device implementation for %s already exists", entry)
		}
		log.Infof("Adding %s to Devicemanager", entry)
		// register model to the given Device implementation
		dm.devices[entry] = device
		device.SetBackend(dm.backend)
	}
	return nil
}

// GetModelHandler retrieves the ModelHandler referenced by its name
func (dm *DeviceManagerImpl) GetModelHandler(model string) (Device, error) {
	deviceHandler, exists := dm.devices[model]
	if !exists {
		return nil, fmt.Errorf("model %s not registered in DeviceManager", model)
	}
	return deviceHandler, nil
}

// SetBackend pushes the given ZtpBackend instance into all of the device instances
func (dm *DeviceManagerImpl) SetBackend(backend backend.ZtpBackend) {
	dm.backend = backend
	for _, device := range dm.devices {
		// updating backend reference for all devices
		device.SetBackend(dm.backend)
	}
}

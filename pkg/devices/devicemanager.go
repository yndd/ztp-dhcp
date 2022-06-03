package devices

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
)

// DeviceManager is the dispatcher and registry for all the vendor specific devices ztp implementations
// Devicemanager is a singleton
var deviceManager *DeviceManagerImpl

// DeviceManagerRegistrator the DeviceManager Interface for registering device model handlers
type DeviceManagerRegistrator interface {
	RegisterDevice([]string, Device) error
}

// DeviceManagerHandler the DeviceManager Interface for the handler side of things
type DeviceManagerHandler interface {
	GetModelHandler(model string) (Device, error)
}

// DeviceManagerImpl the actual implementation of the DeviceManager
type DeviceManagerImpl struct {
	devices map[string]Device
	backend backend.DhcpBackend
}

func GetDeviceManagerHandler() DeviceManagerHandler {
	return newDeviceManager()
}

func GetDeviceManagerRegistrator() DeviceManagerRegistrator {
	return newDeviceManager()
}

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
	}
	return nil
}

func (dm *DeviceManagerImpl) GetModelHandler(model string) (Device, error) {
	deviceHandler, exists := dm.devices[model]
	if !exists {
		return nil, fmt.Errorf("model %s not registered in DeviceManager", model)
	}
	return deviceHandler, nil
}

func (dm *DeviceManagerImpl) SetBackend(backend backend.DhcpBackend) {
	dm.backend = backend
}

package devices

// DeviceManagerRegistrator the DeviceManager Interface for registering device model handlers
type DeviceManagerRegistrator interface {
	RegisterDevice([]string, Device) error
}

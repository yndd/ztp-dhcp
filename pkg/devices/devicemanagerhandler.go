package devices

import "github.com/yndd/ztp-dhcp/pkg/backend"

// DeviceManagerHandler the DeviceManager Interface for the handler side of things
type DeviceManagerHandler interface {
	GetModelHandler(model string) (Device, error)
	SetBackend(backend.ZtpBackend)
}

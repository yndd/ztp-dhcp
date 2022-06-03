package backend

import (
	"errors"

	"github.com/yndd/ztp-dhcp/pkg/structs"
)

type DhcpBackend interface {
	// GetDeviceInformation
	GetDeviceInformation(cir *structs.ClientIdentifier) (*structs.DeviceInformation, error)
}

var ErrDeviceNotFound = errors.New("device not found")

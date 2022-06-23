package backend

import (
	"errors"

	"github.com/yndd/ztp-dhcp/pkg/structs"
)

type DhcpBackend interface {
	// GetDeviceInformation
	GetDeviceInformation(cir *structs.ClientIdentifier) (*structs.DeviceInformation, error)
	// GetWebserverInformation is used to query the backend for information about
	// the webserver that further ztp requests should be forwarded to
	GetWebserverInformation() (*structs.WebserverInfo, error)
}

var ErrDeviceNotFound = errors.New("device not found")

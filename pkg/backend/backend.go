package backend

import (
	"errors"

	"github.com/yndd/ztp-dhcp/pkg/structs"
)

type ZtpBackend interface {
	// GetDeviceInformationByClientIdentifier
	GetDeviceInformationByClientIdentifier(cir *structs.ClientIdentifier) (*structs.DeviceInformation, error)
	GetDeviceInformationByName(deviceId string) (*structs.DeviceInformation, error)
	// GetWebserverInformation is used to query the backend for information about
	// the webserver that further ztp requests should be forwarded to
	GetWebserverInformation() (*structs.WebserverInfo, error)
	GetDhcpserverInformation() (*structs.DhcpServerInfo, error)
}

var ErrDeviceNotFound = errors.New("device not found")

package static

import (
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

type ZtpStaticBackend struct {
	datastore map[string]*structs.DeviceInformation
}

func NewZtpStaticBackend() *ZtpStaticBackend {
	log.Infof("Instantiating ZtpStaticBackend")
	return &ZtpStaticBackend{
		datastore: map[string]*structs.DeviceInformation{},
	}
}

func (f *ZtpStaticBackend) GetDeviceInformation(cir *structs.ClientIdentifierResult) (*structs.DeviceInformation, error) {
	val, exists := f.datastore[cir.Value]
	if !exists {
		return nil, backend.ErrDeviceNotFound
	}
	return val, nil
}

func (f *ZtpStaticBackend) AddEntry(cir *structs.ClientIdentifierResult, di *structs.DeviceInformation) {
	f.datastore[cir.Value] = di
}

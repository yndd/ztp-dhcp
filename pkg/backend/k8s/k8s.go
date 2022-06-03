package k8s

import (
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

type ZtpK8sBackend struct {
}

func NewZtpK8sBackend() *ZtpK8sBackend {
	log.Infof("Instantiating K8sBackend")
	return &ZtpK8sBackend{}
}

func (k *ZtpK8sBackend) GetDeviceInformation(cir *structs.ClientIdentifierResult) (*structs.DeviceInformation, error) {
	return nil, nil
}

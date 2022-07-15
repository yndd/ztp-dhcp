package static

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

type ZtpStaticBackend struct {
	datastore             map[string]*structs.DeviceInformation
	webserverInformation  *structs.WebserverInfo
	dhcpserverInformation *structs.DhcpServerInfo
}

// NewZtpStaticBackend constructs a new StaticZTPBackend.
// The initial content is taken from the file referenced under the
// 'YNDD_ZTP_STATIC_DATASTORE_SOURCE' environemtn variable
func NewZtpStaticBackend() *ZtpStaticBackend {
	log.Infof("Instantiating ZtpStaticBackend")
	backend := &ZtpStaticBackend{
		datastore: map[string]*structs.DeviceInformation{},
		webserverInformation: &structs.WebserverInfo{
			Port:     80,
			IpFqdn:   "127.0.0.1",
			Protocol: "http",
		},
		dhcpserverInformation: &structs.DhcpServerInfo{Ip: net.ParseIP("1.2.3.4")},
	}

	// using the ENV var to provide static config content
	if val, exists := os.LookupEnv("YNDD_ZTP_STATIC_DATASTORE_SOURCE"); exists {
		log.Info("ENVVAR ", val)
		err := backend.loadDataStoreFromFile(val)
		if err != nil {
			log.Errorf("error loading static backend data from %s - %v", val, err)
		}
	}
	return backend
}

func (f *ZtpStaticBackend) GetDhcpserverInformation() (*structs.DhcpServerInfo, error) {
	return f.dhcpserverInformation, nil
}

func (f *ZtpStaticBackend) GetDeviceInformationByClientIdentifier(cir *structs.ClientIdentifier) (*structs.DeviceInformation, error) {
	val, exists := f.datastore[cir.Value]
	if !exists {
		return nil, backend.ErrDeviceNotFound
	}
	return val, nil
}

func (f *ZtpStaticBackend) GetDeviceInformationByName(name string) (*structs.DeviceInformation, error) {
	for _, x := range f.datastore {
		if x.Name == name {
			return x, nil
		}
	}
	return nil, fmt.Errorf("no entry with name '%s' found", name)
}

func (f *ZtpStaticBackend) AddEntry(cir *structs.ClientIdentifier, di *structs.DeviceInformation) {
	f.datastore[cir.Value] = di
}

func (f *ZtpStaticBackend) GetWebserverInformation() (*structs.WebserverInfo, error) {
	return f.webserverInformation, nil
}

func (f *ZtpStaticBackend) loadDataStoreFromFile(path string) error {
	// construct absolute path from the provided path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// read the file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	// unmarshal the data
	backendDatastore := &structs.StaticBackendDatastore{}
	err = json.Unmarshal(data, backendDatastore)
	if err != nil {
		return err
	}

	// add the entries to the StaticBackend
	for _, x := range backendDatastore.Datastore {
		f.AddEntry(x.ClientIdentifier, x.DeviceInformation)
		log.Debugf("adding %s with ClientIdentifier %s to static Datastore", x.DeviceInformation.MacAddress, x.ClientIdentifier.String())
	}

	return nil
}

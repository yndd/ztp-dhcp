package nokiasrl

import (
	"fmt"
	"strconv"

	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/devices"
	"github.com/yndd/ztp-dhcp/pkg/structs"
	websstructs "github.com/yndd/ztp-webserver/pkg/structs"
)

var managed_models = []string{"SRLinux"}

type NokiaSrl struct {
	backend backend.ZtpBackend
}

// AdjustReply takes the general reply and adjusts the device specific values
func (srl *NokiaSrl) AdjustReply(req *dhcpv4.DHCPv4, reply *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {
	// generate the url for the config script
	up := websstructs.NewUrlParamsDeviceName(string(devinfo.VendorType), devinfo.Platform, devinfo.Name, websstructs.Script)
	// retrieve the webserver port host and scheme
	wsinfo, err := srl.backend.GetWebserverInformation()
	if err != nil {
		log.Error(err)
		return
	}
	// retrieve the resulting url.URL from the Parameters
	theUrl := up.GetUrlRelative()
	// set port, Host and scheme
	portAsString := strconv.FormatInt(int64(wsinfo.Port), 10)
	theUrl.Host = fmt.Sprintf("%s:%s", wsinfo.IpFqdn, portAsString)
	theUrl.Scheme = wsinfo.Protocol

	// update / set the DHCP option
	reply.Options.Update(dhcpv4.OptBootFileName(theUrl.String()))
}

// SetBackend used for late binding of the backend
func (srl *NokiaSrl) SetBackend(backend backend.ZtpBackend) {
	srl.backend = backend
}

func init() {
	// register NokiaSRL with the DeviceManager
	err := devices.GetDeviceManagerRegistrator().RegisterDevice(managed_models, &NokiaSrl{})
	if err != nil {
		log.Fatal(err)
	}
}

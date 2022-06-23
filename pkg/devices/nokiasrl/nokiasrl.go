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
	backend backend.DhcpBackend
}

func (srl *NokiaSrl) AdjustReply(req *dhcpv4.DHCPv4, reply *dhcpv4.DHCPv4, devinfo *structs.DeviceInformation) {
	// set Option66
	//reply.Options.Update(dhcpv4.OptTFTPServerName(devinfo.Option66))
	// set Option67

	up := websstructs.NewUrlParamsDeviceId(string(devinfo.VendorType), devinfo.Platform, devinfo.Name, websstructs.Script)
	wsinfo, err := srl.backend.GetWebserverInformation()
	if err != nil {
		log.Error(err)
		return
	}
	theUrl := up.GetUrlRelative()
	portAsString := strconv.FormatInt(int64(wsinfo.Port), 10)
	theUrl.Host = fmt.Sprintf("%s:%s", wsinfo.IpFqdn, portAsString)
	theUrl.Scheme = wsinfo.Protocol
	reply.Options.Update(dhcpv4.OptBootFileName(theUrl.String()))
}

func (srl *NokiaSrl) SetBackend(backend backend.DhcpBackend) {
	srl.backend = backend
}

func init() {
	err := devices.GetDeviceManagerRegistrator().RegisterDevice(managed_models, &NokiaSrl{})
	if err != nil {
		log.Fatal(err)
	}
}

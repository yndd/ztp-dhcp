package dhcp

import (
	"encoding/base64"
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/devices"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

// ZtpServer is the ZTP Server
type ZtpServer struct {
	backend       backend.DhcpBackend
	deviceManager devices.DeviceManagerInterf
	settings      *ZtpSettings
}

// ZtpSettings contains settings to configure the ZTP Server
type ZtpSettings struct {
	LeaseTime uint32 // the default lease time assigned to any IP assignments
}

// NewZtpServer konstructor for a new ZtpServer instance
func NewZtpServer(backend backend.DhcpBackend, ztpSettings *ZtpSettings) *ZtpServer {
	return &ZtpServer{
		backend:       backend,
		deviceManager: devices.DeviceManager,
		settings:      ztpSettings,
	}
}

// handler is called whenever a Packet is received on the socket that the DHCP server is bound to
// from here we branch into handling the different message types and generate the responses
func (z *ZtpServer) handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// this function will just print the received DHCPv4 message, without replying
	if m == nil {
		log.Warn("Packet is nil!")
		return
	}
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		log.Warn("Not a BootRequest!")
		return
	}

	log.Debug("Received Paket:")
	log.Debug(m.Summary())
	log.Debugf("<%s>", base64.StdEncoding.EncodeToString(m.ToBytes()))

	var reply *dhcpv4.DHCPv4
	var err error

	// figure out if DHCP Offer (handleDiscover) or
	// DHCP Ack (handleRequest) is to be send
	switch m.MessageType() {
	case dhcpv4.MessageTypeDiscover:
		reply, err = z.handleDiscover(conn, peer, m)
	case dhcpv4.MessageTypeRequest:
		reply, err = z.handleRequest(conn, peer, m)
	default:
		z.fallbackHandler(conn, peer, m)
	}
	if err != nil {
		log.Error(err)
		return
	}

	// push reply to wire
	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		log.Printf("Cannot reply to client: %v", err)
	}
	log.Debugf("REPLY: %s", reply.Summary())
}

// handleDiscover handles DHCPDISCOVER messages and crafts a DHCPACK in response
func (z *ZtpServer) handleDiscover(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, error) {
	log.Info("DiscoverHandler - Start")

	//log.Print(m.ClassIdentifier())
	//ci := m.Options.Get(dhcpv4.OptionClientIdentifier)

	//vendor_data, err := ztpv4.ParseVendorData(m)
	//if err != nil {
	//	log.Error(err)
	// return nil, err
	//}

	//println("ClassID: %s %s %s", vendor_data.VendorName, vendor_data.Model, vendor_data.Serial)

	log.Info("CLIENTIDENT")
	ciresult, err := GetClientIdentifier(m)
	if err != nil {
		// no ClientIdentifier present or error happened
		return nil, err
	}

	deviceInfo, err := z.backend.GetDeviceInformation(ciresult)
	if err != nil {
		return nil, err
	}

	// create a reply from the request
	reply, err := createReply(m, dhcpv4.MessageTypeOffer, z.settings.LeaseTime, deviceInfo)
	if err != nil {
		log.Error(err)
	}

	log.Info("DiscoverHandler - Done")
	return reply, nil
}

// handleRequest handles DHCPREQUEST messages and crafts a DHCPACK in response
func (z *ZtpServer) handleRequest(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, error) {
	log.Info("RequestHandler - Start")

	ciresult, err := GetClientIdentifier(m)
	if err != nil {
		// no ClientIdentifier present or error happened
		return nil, err
	}

	deviceInfo, err := z.backend.GetDeviceInformation(ciresult)
	if err != nil {
		return nil, err
	}

	// create a reply from the request
	reply, err := createReply(m, dhcpv4.MessageTypeAck, z.settings.LeaseTime, deviceInfo)
	if err != nil {
		log.Error(err)
	}
	log.Info("RequestHandler - Done")
	return reply, nil
}

// StartDHCPServer instantiates the socket, binds it to the given interface if ifname != "" and
// starts the DHCP Server process on that socket.
func (z *ZtpServer) StartDHCPServer(serverport int, ifname string) {
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: serverport,
	}
	// define the server
	server, err := server4.NewServer(ifname, laddr, z.handler, server4.WithDebugLogger())
	if err != nil {
		log.Fatal(err)
	}
	// instead of empty string for all interfaces print ANY
	ifprintnane := "ANY"
	if ifname != "" {
		ifprintnane = ifname
	}
	log.Infof("Starting DHCP-Server on interface %s", ifprintnane)

	// start server
	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

// createReply create a DHCP Reply from a Request.
func createReply(m *dhcpv4.DHCPv4, messageType dhcpv4.MessageType, leaseTime uint32, deviceInfo *structs.DeviceInformation) (*dhcpv4.DHCPv4, error) {
	// Create the reply from the request
	reply, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		log.Printf("NewReplyFromRequest failed: %v", err)
		return nil, err
	}
	// Change the Message Type to DHCP-ACK
	reply.UpdateOption(dhcpv4.OptMessageType(messageType))
	// Set the lease time
	dhcpv4.WithLeaseTime(leaseTime)(reply)

	// Parse CIDR information
	clientip, clientnetmask, err := net.ParseCIDR(deviceInfo.CIDR)
	if err != nil {
		return nil, fmt.Errorf("Error parsing CIDR (%s): %v", deviceInfo.CIDR, err)
	}
	// Fill in the to be offered client ip
	dhcpv4.WithYourIP(clientip)(reply)
	// Fill in the netmask
	dhcpv4.WithNetmask(clientnetmask.Mask)(reply)

	// add router information if specified in deviceinfo
	if deviceInfo.Gateway != "" {
		router := net.ParseIP(deviceInfo.Gateway)
		if router == nil {
			return nil, fmt.Errorf("router (%s) could not be parsed as an IPv4 address", deviceInfo.Gateway)
		}
		dhcpv4.WithRouter(router)(reply)
	}

	for _, dnsentry := range deviceInfo.DnsServersV4 {
		dnsserver := net.ParseIP(dnsentry)
		if dnsserver == nil {
			return nil, fmt.Errorf("DNS Server (%s) could not be parsed as an IPv4 address", dnsentry)
		}
		dhcpv4.WithDNS(dnsserver)(reply)
	}

	return reply, nil
}

// fallbackHandler is the fallback handler for non implemented DHCP Message Types
func (z *ZtpServer) fallbackHandler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Warn("NO HANDLER FOR THE RECEIVED PACKET")
}

package dhcp

import (
	"encoding/base64"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/yndd/ztp-dhcp/pkg/backend/static"
	_ "github.com/yndd/ztp-dhcp/pkg/devices/all"
	"github.com/yndd/ztp-dhcp/pkg/mocks"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

func TestDhcp(t *testing.T) {

	// prepare the static Backend
	staticBackend := static.NewZtpStaticBackend()
	// add an entry to the backend
	staticBackend.AddEntry(
		&structs.ClientIdentifier{
			CIType: 0,
			Value:  "this is the identifier",
		},
		&structs.DeviceInformation{
			Name:              "MyFunnyTestDevice",
			MacAddress:        "b6:8d:0b:94:62:8d",
			Platform:          "Dummy",
			ExpectedSWVersion: "5.4.2",
			CIDR:              "192.168.50.33/24",
			SerialNumber:      "666",
			NtpServersV4:      []string{"1.2.3.4"},
			DnsServersV4:      []string{"8.8.8.8"},
		},
	)
	staticBackend.AddEntry(
		&structs.ClientIdentifier{
			CIType: 1,
			Value:  "b6:8d:0b:94:62:8d",
		},
		&structs.DeviceInformation{
			Name:              "MyFunnyTestDevice",
			MacAddress:        "b6:8d:0b:94:62:8d",
			ExpectedSWVersion: "5.4.2",
			Platform:          "Dummy",
			CIDR:              "192.168.50.33/24",
			SerialNumber:      "666",
			NtpServersV4:      []string{"1.2.3.4"},
			DnsServersV4:      []string{"8.8.8.8"},
		},
	)
	staticBackend.AddEntry(
		&structs.ClientIdentifier{
			CIType: 0,
			Value:  "NS2113T0295",
		},
		&structs.DeviceInformation{
			Name:              "IXR",
			MacAddress:        "b6:8d:0b:94:62:aa",
			ExpectedSWVersion: "5.4.2",
			Platform:          "SRLinux",
			CIDR:              "192.168.50.33/24",
			SerialNumber:      "NS2113T0295",
			NtpServersV4:      []string{"1.2.3.4"},
			DnsServersV4:      []string{"8.8.8.8"},
		},
	)
	// instantiate the ZTP Server, which is basically the handler implementation
	ztpServer := NewZtpServer(staticBackend, &ZtpSettings{LeaseTime: 3600})

	// run through the testdata
	for _, entry := range testData {
		//t.Log(entry.Description)

		// decode the Base64 encoded packet
		foo, err := base64.StdEncoding.DecodeString(entry.B64data)
		if err != nil {
			t.Error()
			t.Fail()
		}
		// Parse the packet to become a dhcpv4 library struct
		paket, err := dhcpv4.FromBytes(foo)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		// instantiate a Fake PacketConn, required as a parameter to the handler function
		packetConn := mocks.NewMockPacketConn(mockCtrl)

		ipaddr := &net.IPAddr{}

		// if we expect the handler to succeede the Expectation on the Mock is set
		if entry.SucceedWriteTo {
			packetConn.EXPECT().WriteTo(gomock.Any(), ipaddr)
		}

		// finally call the handler to process the packet
		ztpServer.handler(packetConn, ipaddr, paket)
	}
}

type testDataEntry struct {
	Description    string // description of the b64data, basically a human readable dump
	B64data        string // the request packet encoded as Base64
	SucceedWriteTo bool   // used to check if the ZTP Server made it to the point where it wrote the response packet to the packetConn
	Option61Type   structs.CITypeEnum
}

var testData = []*testDataEntry{
	{
		Description: `
		opcode: BootRequest
		hwtype: Ethernet
		hopcount: 0
		transaction ID: 0x7431a324
		num seconds: 0
		flags: Unicast (0x00)
		client IP: 0.0.0.0
		your IP: 0.0.0.0
		server IP: 0.0.0.0
		gateway IP: 0.0.0.0
		client MAC: b6:8d:0b:94:62:8d
		server hostname:
		bootfile name:
		options:
		  Host Name: vbox
		  DHCP Message Type: DISCOVER
		  Parameter Request List: Subnet Mask, Time Offset, Router, Domain Name Server, Host Name, Domain Name, Interface MTU, Broadcast Address, NTP Servers, NetBIOS over TCP/IP Name Server, NetBIOS over TCP/IP Scope, DNS Domain Search List, Classless Static Route
		  Client identifier: [116 104 105 115 32 105 115 32 116 104 101 32 105 100 101 110 116 105 102 105 101 114]`,
		B64data:        "AQEGAHQxoyQAAAAAAAAAAAAAAAAAAAAAAAAAALaNC5RijQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABjglNjDAR2Ym94NQEBNw0BHAIDDwZ3DCwvGnkqPRcAdGhpcyBpcyB0aGUgaWRlbnRpZmllcv8AAAAAAAAAAAAA",
		SucceedWriteTo: true,
		Option61Type:   structs.String,
	},
	{
		Description: `
		opcode: BootRequest
		hwtype: Ethernet
		hopcount: 0
		transaction ID: 0x7431a324
		num seconds: 0
		flags: Unicast (0x00)
		client IP: 0.0.0.0
		your IP: 0.0.0.0
		server IP: 0.0.0.0
		gateway IP: 0.0.0.0
		client MAC: b6:8d:0b:94:62:8d
		server hostname:
		bootfile name:
		options:
		  Host Name: vbox
		  Requested IP Address: 192.168.50.33
		  DHCP Message Type: REQUEST
		  Parameter Request List: Subnet Mask, Time Offset, Router, Domain Name Server, Host Name, Domain Name, Interface MTU, Broadcast Address, NTP Servers, NetBIOS over TCP/IP Name Server, NetBIOS over TCP/IP Scope, DNS Domain Search List, Classless Static Route
		  Client identifier: [116 104 105 115 32 105 115 32 116 104 101 32 105 100 101 110 116 105 102 105 101 114]
		`,
		B64data:        "AQEGAHQxoyQAAAAAAAAAAAAAAAAAAAAAAAAAALaNC5RijQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABjglNjDAR2Ym94MgTAqDIhNQEDNw0BHAIDDwZ3DCwvGnkqPRcAdGhpcyBpcyB0aGUgaWRlbnRpZmllcv8AAAAA",
		SucceedWriteTo: true,
		Option61Type:   structs.String,
	},
	{
		Description: `
		opcode: BootRequest
		hwtype: Ethernet
		hopcount: 0
		transaction ID: 0x5c660802
		num seconds: 0
		flags: Unicast (0x00)
		client IP: 0.0.0.0
		your IP: 0.0.0.0
		server IP: 0.0.0.0
		gateway IP: 0.0.0.0
		client MAC: b6:8d:0b:94:62:8d
		server hostname:
		bootfile name:
		options:
		  Host Name: vbox
		  Requested IP Address: 192.168.50.33
		  DHCP Message Type: REQUEST
		  Parameter Request List: Subnet Mask, Time Offset, Router, Domain Name Server, Host Name, Domain Name, Interface MTU, Broadcast Address, NTP Servers, NetBIOS over TCP/IP Name Server, NetBIOS over TCP/IP Scope, DNS Domain Search List, Classless Static Route
		  Client identifier: [1 182 141 11 148 98 141]
		`,
		B64data:        "AQEGAFxmCAIAAAAAAAAAAAAAAAAAAAAAAAAAALaNC5RijQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABjglNjDAR2Ym94MgTAqDIhNQEDNw0BHAIDDwZ3DCwvGnkqPQcBto0LlGKN/wAAAAAAAAAAAAAAAAAAAAAAAAAA",
		SucceedWriteTo: true,
		Option61Type:   structs.MAC,
	},
	{
		Description: `
		opcode: BootRequest
		hwtype: Ethernet
		hopcount: 0
		transaction ID: 0x8461f41d
		num seconds: 0
		flags: Unicast (0x00)
		client IP: 0.0.0.0
		your IP: 0.0.0.0
		server IP: 0.0.0.0
		gateway IP: 0.0.0.0
		client MAC: b6:8d:0b:94:62:8d
		server hostname:
		bootfile name:
		options:
		  Host Name: vbox
		  DHCP Message Type: DISCOVER
		  Parameter Request List: Subnet Mask, Time Offset, Router, Domain Name Server, Host Name, Domain Name, Interface MTU, Broadcast Address, NTP Servers, NetBIOS over TCP/IP Name Server, NetBIOS over TCP/IP Scope, DNS Domain Search List, Classless Static Route
		  Client identifier: [1 182 141 11 148 98 141]
		`,
		B64data:        "AQEGAIRh9B0AAAAAAAAAAAAAAAAAAAAAAAAAALaNC5RijQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABjglNjDAR2Ym94NQEBNw0BHAIDDwZ3DCwvGnkqPQcBto0LlGKN/wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		SucceedWriteTo: true,
		Option61Type:   structs.MAC,
	},
	{
		Description: `
		SRLINUX Discovery via relay agent
		`,
		SucceedWriteTo: true,
		Option61Type:   structs.String,
		B64data:        "AQEGAYxe73AADoAAAAAAAAAAAAAAAAAAwKgWAVDg75lgUQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABjglNjMwT/////NQEBNwsBAwYHDA8qK0JDdzwRTk9LSUE6NzIyMCBJWFItRDE9DABOUzIxMTNUMDI5NVIIAQZlbnA2czD/",
	},
}

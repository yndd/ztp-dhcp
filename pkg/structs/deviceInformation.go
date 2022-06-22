package structs

import targetv1 "github.com/yndd/target/apis/target/v1"

type DeviceInformation struct {
	// name of the device
	Name string `json:"name"`
	// the MacAddress of the device
	MacAddress string `json:"macAddress"`
	// The serialnumber of the device
	SerialNumber string `json:"serialNo"`
	// The IP address and netmask in CIDR notation (e.g. 10.10.10.10/24)
	CIDR string `json:"cidr"`
	// The IP address of the Gateway
	Gateway string `json:"gateway,omitempty"`

	// NtpServer addresses NO DNS NAMES only IPv4 addresses
	NtpServersV4 []string `json:"ntpServersV4,omitempty"`
	// DnsServer addresses
	DnsServersV4 []string `json:"dnsServersV4,omitempty"`

	ExpectedSWVersion string `json:"expectedSwVersion,omitempty"`

	// reference to the Platform of the device
	Platform   string              `json:"plattform"`
	VendorType targetv1.VendorType `json:"vendor"`
}

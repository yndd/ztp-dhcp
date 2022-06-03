package structs

type DeviceInformation struct {
	// name of the device
	Name string
	// the MacAddress of the device
	MacAddress string
	// The serialnumber of the device
	SerialNumber string
	// The IP address and netmask in CIDR notation (e.g. 10.10.10.10/24)
	CIDR string
	// The IP address of the Gateway
	Gateway string

	// NtpServer addresses NO DNS NAMES only IPv4 addresses
	NtpServersV4 []string
	// DnsServer addresses
	DnsServersV4 []string

	// either
	ExpectedSWVersion string
	Config            string // how will this be referenced
	// or
	Option66 string
	Option67 string
	Option43 string
}

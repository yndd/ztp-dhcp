package structs

import "net"

type DhcpServerInfo struct {
	Ip net.IP `json:"ip"`
}

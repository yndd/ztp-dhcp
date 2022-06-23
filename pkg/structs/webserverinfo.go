package structs

type WebserverInfo struct {
	Port     int32  `json:"port"`
	IpFqdn   string `json:"ipFqdn"`
	Protocol string `json:"protocol"`
}

package dhcp

import (
	"net"
	"reflect"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

func TestGetClientIdentifier(t *testing.T) {
	type args struct {
		m *dhcpv4.DHCPv4
	}
	tests := []struct {
		name    string
		args    args
		want    *structs.ClientIdentifier
		wantErr bool
	}{
		{
			name: "Basic MAC client identifier",
			args: args{
				m: addClientID(createDiscover(t), byte(0x1), []byte{0x1, 0x2, 0x3, 0x4, 0x56, 0xad}),
			},
			want: &structs.ClientIdentifier{
				CIType: structs.MAC,
				Value:  "01:02:03:04:56:ad",
			},
			wantErr: false,
		},
		{
			name: "Basic String client identifier",
			args: args{
				m: addClientID(createDiscover(t), byte(0x0), []byte("Hello this is the ID")),
			},
			want: &structs.ClientIdentifier{
				CIType: structs.String,
				Value:  "Hello this is the ID",
			},
			wantErr: false,
		},
		{
			name: "Fail, no ClientID present",
			args: args{
				m: createDiscover(t),
			},
			wantErr: true,
		},
		{
			name: "Basic MAC client identifier",
			args: args{
				m: addClientID(createDiscover(t), byte(0x1), []byte{0x3, 0x4, 0x56, 0xad}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClientIdentifier(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClientIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClientIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

// createDiscover helper to create a DhcpDiscover paket with client identifier set
func createDiscover(t *testing.T) *dhcpv4.DHCPv4 {
	hwaddr := net.HardwareAddr([]byte{0x1, 0x2, 0x3, 0x4, 0x56, 0xad})
	dhcpPak, err := dhcpv4.NewDiscovery(hwaddr)
	if err != nil {
		t.Error("error creating dhcp paket.")
	}
	return dhcpPak
}

func addClientID(dhcpPak *dhcpv4.DHCPv4, citype byte, civalue []byte) *dhcpv4.DHCPv4 {
	clientid := append([]byte{citype}, civalue...)
	dhcpv4.WithOption(dhcpv4.OptClientIdentifier(clientid))(dhcpPak)
	return dhcpPak
}

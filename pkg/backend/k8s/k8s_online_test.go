//go:build k8sonline

// Use `go test ./... --tags k8sonline` to run these tests.
// It requires having a working ~/.kube/config present and the k8s cluster reachable
package k8s

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yndd/ztp-dhcp/pkg/structs"
)

func TestK8s(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error("Unable to obtain home directory")
	}
	zb := NewZtpK8sBackend(filepath.Join(home, ".kube", "config"))
	zb.GetDeviceInformationByClientIdentifier(&structs.ClientIdentifier{CIType: structs.MAC, Value: "b6:8d:0b:94:62:8d"})
}

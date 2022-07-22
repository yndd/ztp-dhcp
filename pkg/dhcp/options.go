package dhcp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/structs"
)

// ErrNoClientIdentifier used to indicate to the caller that no Option 61 was present
var ErrNoClientIdentifier = errors.New("no ClientIdentifier present")

// GetClientIdentifier is used to retrieve the
// Option 61 (ClientIdentifier) from the given dhcp paket.
// NoClientIdentifierError is thrown when no Option 61 is present
func GetClientIdentifier(m *dhcpv4.DHCPv4) (*structs.ClientIdentifier, error) {
	log.Debug("Calling getClientIdentifier()")
	// retrieve the option parameter as
	ci := m.Options.Get(dhcpv4.OptionClientIdentifier)

	// check if Option 61 is present
	if len(ci) == 0 {
		return nil, ErrNoClientIdentifier
	}

	// create the ClientIdentifierResult struct
	cir := &structs.ClientIdentifier{}

	switch ci[0] {
	case uint8(1): // Type 1 indicates that a MAC address is present in Option 61
		log.Debug("Processing MAC ADDRESS based ClientIdentifier")
		cir.CIType = structs.MAC

		tmpmac := hex.EncodeToString(ci[1:])
		if len(tmpmac) != 12 {
			return nil, fmt.Errorf("not a valid MAC address: %s", tmpmac)
		}

		macsl := []string{}
		for i := 0; i < len(tmpmac); i = i + 2 {
			macsl = append(macsl, tmpmac[i:i+2])
		}
		cir.Value = strings.Join(macsl, ":")
	default: // Type 0 usually is string, we define this as default for now
		log.Debug("Processing String based ClientIdentifier")
		cir.CIType = structs.String
		cir.Value = string(ci[1:])
	}

	log.Debugf("ClientIdentifier: %s", cir.String())
	return cir, nil
}

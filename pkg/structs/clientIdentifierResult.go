package structs

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// CITypeEnum used to indicate the type of
// data encoded in Option 61
type CITypeEnum int

const (
	String CITypeEnum = iota
	MAC
)

// ClientIdentifier is the struct to deliver the result of the Option 61 parsing
type ClientIdentifier struct {
	CIType CITypeEnum `json:"type"`
	Value  string     `json:"value"`
}

// String provides a string representation of the ClientIdentifierResult
// actually that is mapped to its json encoding
func (cir *ClientIdentifier) String() string {
	result, err := json.Marshal(cir)
	if err != nil {
		log.Errorf("Error marshalling ClientIdentifierResult: %v", err)
		return ""
	}
	return string(result)
}

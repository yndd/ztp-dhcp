package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// DeduceNamespace evaluates the "POD_NAMESPACE" Environment variable and
// returns its value if set, otherwise it will return the given default value
func DeduceNamespace(defaultNS string) string {
	var isset bool
	var namespace string
	// deduce, which namespace is to be queried for the nodes from the environment variable POD_NAMESPACE
	if namespace, isset = os.LookupEnv("POD_NAMESPACE"); !isset {
		log.Warnf("environment variable 'POD_NAMESPACE' not set. falling back to 'ndd-system'")
		namespace = defaultNS
	}
	return namespace
}

package logsetup

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// This is to setup logging as early as possible.
// doing it within some cmd/<package> init() method, will cause the
// init of all the imports to be called already. As an import it is run early
func init() {
	log.SetFormatter(&log.TextFormatter{DisableColors: false, FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	// log.SetReportCaller(true)
	log.Info("Logger setup")
}

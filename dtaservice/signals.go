package dtaservice

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// CaptureSignals spans a signal handler for SIGINT and SIGTERM
func CaptureSignals(server *DocTransServer, registerURL string, wg *sync.WaitGroup) {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go HandleSignals(server, signalCh, registerURL, wg)
}

// HandleSignals reacts on Signals by managing registration at registry.
// On SIGINT (CTRL-C) Unregisters
// On SIGTERM (CTRL-D) Toogling Registration
func HandleSignals(server *DocTransServer, signalCh chan os.Signal, registerURL string, wg *sync.WaitGroup) {
	defer wg.Done()
	for sigs := range signalCh {
		switch sigs {
		case syscall.SIGTERM: // CTRL-D
			log.Debugln("Received SIGTERM")
			if server.InstanceInfo() != nil {
				server.UnregisterAtRegistry()
			} else {
				server.RegisterAtRegistry(registerURL)
			}
		case syscall.SIGINT: // CTRL-C
			log.Debugln("Received SIGINT")
			if server.InstanceInfo() != nil {
				server.UnregisterAtRegistry()
			}
			return
		}
	}
}

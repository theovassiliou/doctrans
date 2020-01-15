package qdsservices

import (
	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
	"os"
	"os/signal"
	"syscall"
)

// CaptureSignals spans a signal handler for SIGINT and SIGTERM
func CaptureSignals(server *pb.DocTransServer) {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go HandleSignals(server, signalCh)
}

// HandleSignals reacts on Signals by managing registration at registry.
// On SIGINT (CTRL-C) Unregisters
// On SIGTERM (CTRL-D) Toogling Registration
func HandleSignals(server *pb.DocTransServer, signalCh chan os.Signal) {
	for sigs := range signalCh {
		switch sigs {
		case syscall.SIGTERM: // CTRL-D
			log.Debugln("Received SIGTERM")
			if server.InstanceInfo() != nil {
				server.UnregisterAtRegistry()
			} else {
				server.RegisterAtRegistry(server.HostName, server.AppName, aux.GetIPAdress(), server.PortToListen, "Service", server.TTL, server.IsSSL)
			}
		case syscall.SIGINT: // CTRL-C
			log.Debugln("Received SIGINT")
			if server.InstanceInfo() != nil {
				server.UnregisterAtRegistry()
			}
			os.Exit(1)
		}
	}
}

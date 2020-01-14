package qdsservices

import (
	log "github.com/sirupsen/logrus"
	aux "github.com/theovassiliou/doctrans/ipaux"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	"os"
	"os/signal"
	"syscall"
)

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

func CaptureSignals(server *pb.DocTransServer) {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go HandleSignals(server, signalCh)
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	"github.com/theovassiliou/doctrans/qdsservices"
	count "github.com/theovassiliou/doctrans/services/qds_count_generic/serviceimplementation"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	serviceName = "COUNT"
)

// DtaService holds the infrastructure for performing the service.

func main() {
	workingHomeDir, _ := homedir.Dir()

	dts := &pb.DocTransServer{
		AppName:  serviceName,
		CfgFile:  workingHomeDir + "/.dta/" + serviceName + "/config.json",
		LogLevel: log.WarnLevel,
	}

	// (1) SetUp Configuration
	pb.SetupConfiguration(dts, workingHomeDir, VERSION)

	// init the resolver so that we have access to the list of apps
	service := &count.Manager{
		SrvHandler: dts,
		Registry: eureka.NewClient([]string{
			dts.RegistrarURL,
			// add others servers here
		}),
	}

	// (2) Init and register GRPC Service
	lis := pb.GrpcLisInitAndReg(service.SrvHandler)

	go pb.StartGrpcServer(lis, service)

	// Start dta service by using the listener

	if dts.REST {
		//(3) Let's instanciate the the HTTP Server
		qdsservices.CaptureSignals(service.SrvHandler)
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		pb.MuxHTTPGrpc(ctx, dts.HTTPPort, service.SrvHandler)
	} else {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		qdsservices.HandleSignals(service.SrvHandler, signalCh)
	}
	return
}

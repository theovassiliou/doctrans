package main

import (
	"context"

	qdstemplate "github.com/theovassiliou/doctrans/qdsservices"

	count "github.com/theovassiliou/doctrans/services/qds_count_generic/serviceimplementation"

	aux "github.com/theovassiliou/doctrans/ipaux"
	grpc "google.golang.org/grpc"

	"github.com/carlescere/scheduler"
	"github.com/jpillora/opts"
	"github.com/theovassiliou/go-eureka-client/eureka"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	serviceName = "COUNT"
)

// DtaService holds the infrastructure for performing the service.

func main() {

	countWorker := &qdstemplate.AQdsService{
		GRPCPortToListen: "50051",
		HostName:         aux.GetHostname(),
		LogLevel:         log.WarnLevel,
		RegisterURL:      qdstemplate.EurekaURL,
		TTL:              qdstemplate.EurekaTTL,
	}

	opts.New(countWorker).
		Repo("github.com/theovassiliou/doctrans").
		Version(version).
		Parse()

	if countWorker.LogLevel != 0 {
		log.SetLevel(countWorker.LogLevel)
	}

	if countWorker.IsSSL {
		log.Warnln("SSL currently not support. Ignoring,")
		countWorker.IsSSL = false
	}

	// 1. Load config
	// a. load service name
	// service name build in
	// b. get galaxy name
	// read from command line

	// c. get grpc ports / get REST ports
	// default ports
	// GRPC/GRPCS: 50051/60051
	// HTTP/HTTPS : 80 / 8080

	// enable optional SSL

	// 2. Start GRPC and/or REST

	// 3. Register @ Registry
	// "http://127.0.0.1:8761/eureka"

	// 4. Enable signals

	// (1) SetUp Configuration
	//	pb.SetupConfiguration(dts, workingHomeDir, VERSION)

	// init the resolver so that we have access to the list of apps

	// (2) Init and register GRPC Service
	// We first create the listener to know the dynamically allocated port we listen on
	_configuredPort := countWorker.GRPCPortToListen

	lis, assignedPort := countWorker.CreateListener()
	if _configuredPort != assignedPort {
		log.Warnf("Listing on port %v instead on configured, but port-inuse %v\n", assignedPort, _configuredPort)
	}

	countWorker.GRPCPortToListen = assignedPort

	if countWorker.Register {

		countWorker.SetRegistry(eureka.NewClient([]string{
			countWorker.RegisterURL,
			// add others servers here
		}))

		countWorker.RegisterGRPCService(serviceName, aux.GetIPAdress())

		job := countWorker.GetHearbeatFunc()
		// Run every 25 seconds but not now.
		hbj, _ := scheduler.Every(int(countWorker.TTL)).Seconds().NotImmediately().Run(job)
		countWorker.SetHeartBeatJob(hbj)
	}

	var impl = count.Implementation{}

	go func() {
		s := grpc.NewServer()
		pb.RegisterDTAServerServer(s, impl)
		if err := s.Serve(lis); err != nil {
			log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
		}
	}()

	if countWorker.REST {
		//(3) Let's instanciate the the HTTP Server
		// FIXME: Reimplement signals
		// qdsservices.CaptureSignals(dts)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		_configuredHTTPPort := countWorker.HTTPPort

		assignedHTTPPort := countWorker.RegisterHTTPService(ctx)

		if _configuredHTTPPort != assignedHTTPPort {
			log.Warnf("Listing on port %v instead on configured, but port-inuse %v\n", assignedPort, _configuredPort)
		}
		countWorker.HTTPPort = assignedHTTPPort

		// FIXME: Continue here
	}
	// else {
	// 	signalCh := make(chan os.Signal)
	// 	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// 	qdsservices.HandleSignals(dts, signalCh)
	// }
	return
}

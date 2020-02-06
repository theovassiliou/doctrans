package main

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	si "github.com/theovassiliou/doctrans/services/qds_count_generic/serviceimplementation"

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

type Manager struct {
	pb.UnimplementedDTAServerServer

	GalaxyName       string `opts:"group=Service" help:"ID of the service"`
	GRPCPortToListen string `opts:"group=Service" help:"On which port to listen for this service."`
	IsSSL            bool   `opts:"group=Service" help:"Service reached via SSL, if set."`
	REST             bool   `opts:"group=Service" help:"REST-API enabled on port 80, if set"`
	HTTPPort         string `opts:"group=Service" help:"On which httpPort to listen for REST, if enableREST is set. Ignored otherwise."`

	Register      bool   `opts:"group=Registrar" help:"Register service with EUREKA, if set"`
	RegistrarURL  string `opts:"group=Registrar" help:"Registry URL"`
	RegistrarUser string `opts:"group=Registrar" help:"Registry User, no user used if not provided"`
	RegistrarPWD  string `opts:"group=Registrar" help:"Registry User Password, no password used if not provided"`
	TTL           uint   `opts:"group=Registrar" help:"Time in seconds to reregister at Registrar."`
	HostName      string `opts:"group=Registrar" help:"If provided will be used as hostname, else automatically derived."`

	LogLevel log.Level `opts:"group=Generic" help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`

	registry     *eureka.Client
	instanceInfo *eureka.InstanceInfo
	heartBeatJob *scheduler.Job
}

// DtaService holds the infrastructure for performing the service.

func main() {

	dts := &Manager{
		GRPCPortToListen: "50051",
		HostName:         aux.GetHostname(),
		LogLevel:         log.WarnLevel,
	}

	opts.New(dts).
		Repo("github.com/theovassiliou/doctrans").
		Version(version).
		Parse()

	if dts.LogLevel != 0 {
		log.SetLevel(dts.LogLevel)
	}

	dts.registry = eureka.NewClient([]string{
		dts.RegistrarURL,
		// add others servers here
	})

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
	lis := func(srvHandler *Manager) net.Listener {
		lis :=
			func(srvHandler *Manager) net.Listener {
				// We first create the listener to know the dynamically allocated port we listen on
				const maxPortSeek int = 20
				_configuredPort := srvHandler.GRPCPortToListen
				lis := func(maxPortSeek int) net.Listener {
					var lis net.Listener
					var err error

					port, err := strconv.Atoi(dts.GRPCPortToListen)

					for i := 0; i < maxPortSeek; i++ {
						log.WithFields(log.Fields{"Service": "Server", "Status": "Trying"}).Infof("Trying to listen on port %d", (port + i))
						lis, err = net.Listen("tcp", ":"+strconv.Itoa(port+i))
						if err == nil {
							port = port + i
							log.WithFields(log.Fields{"Service": "Server", "Status": "Listening"}).Infof("Using port %d to listen for dta", port)
							i = maxPortSeek
						}
					}

					if err != nil {
						log.WithFields(log.Fields{"Service": "Server", "Status": "Abort"}).Infof("Failed to finally open ports between %d and %d", port, port+maxPortSeek)
						log.Fatalf("failed to listen: %v", err)
					}

					log.WithFields(log.Fields{"Service": "Server", "Status": "main"}).Debugln("Opend successfull a port")
					dts.GRPCPortToListen = strconv.Itoa(port)
					return lis
				}(10)

				if _configuredPort != srvHandler.GRPCPortToListen {
					log.Warnf("Listing on port %v instead on configured, but used port %v\n", srvHandler.GRPCPortToListen, _configuredPort)
				}
				return lis
			}(dts)

		// We register ourselfs by using the dyn.port
		if srvHandler.Register {
			func(hostname, app, ipAddress, port string, ttl uint, isSsl bool) {

				dts.registry.CheckRetry = eureka.ExpBackOffCheckRetry
				// Create the app instance
				dts.instanceInfo = eureka.NewInstanceInfo(hostname, app, ipAddress, port, ttl, isSsl) //Create a new instance to register
				// Add some meta data. Currently no meaning
				// TODO: Remove this playground if not further required
				dts.instanceInfo.Metadata = &eureka.MetaData{
					Map: make(map[string]string),
				}
				dts.instanceInfo.Metadata.Map["DTA-Type"] = "Service" //one of Gateway, Service
				// Register instance and heartbeat for Eureka
				dts.registry.RegisterInstance(app, dts.instanceInfo) // Register new instance in your eureka(s)
				log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Init"}).Infof("Registering service %s\n", app)

				job := func() {
					log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Up"}).Trace("sending heartbeat : %v\n", time.Now().UTC())
					dts.registry.SendHeartbeat(dts.instanceInfo.App, dts.instanceInfo.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
				}

				// Run every 25 seconds but not now.
				// FIXME:0 We have somehow be able to deregister the heartbeat
				dts.heartBeatJob, _ = scheduler.Every(25).Seconds().NotImmediately().Run(job)
			}(srvHandler.HostName, serviceName, aux.GetIPAdress(), srvHandler.GRPCPortToListen, srvHandler.TTL, srvHandler.IsSSL)

		}

		return lis
	}(dts)

	var theWorker = si.WorkerImplementation{}
	go func(lis net.Listener, dtaServer Manager) {
		s := grpc.NewServer()
		pb.RegisterDTAServerServer(s, theWorker)
		if err := s.Serve(lis); err != nil {
			log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
		}
	}(lis, *dts)

	// Start dta service by using the listener

	if dts.REST {
		//(3) Let's instanciate the the HTTP Server
		// FIXME: Reimplement signals
		// qdsservices.CaptureSignals(dts)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		func(ctx context.Context, HTTPPort string, srvHandler *Manager) {
			mux := runtime.NewServeMux()
			opts := []grpc.DialOption{grpc.WithInsecure()}
			grpcPort := srvHandler.GRPCPortToListen
			log.Debugf("GRPC Endpoint localhost:%s\n", grpcPort)
			err := pb.RegisterDTAServerHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
			if err != nil {
				log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to register: %v", err)
			}

			// (4) Start HTTP Server
			// Start HTTP server (and proxy calls to gRPC server endpoint)
			log.WithFields(log.Fields{"Service": "HTTP", "Status": "Running"}).Debugf("Starting HTTP server on: %v", HTTPPort)
			// FIXME implement a fall back if port is in use.
			if err := http.ListenAndServe(":"+HTTPPort, mux); err != nil {
				log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
			}
		}(ctx, dts.HTTPPort, dts)

	}
	// else {
	// 	signalCh := make(chan os.Signal)
	// 	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// 	qdsservices.HandleSignals(dts, signalCh)
	// }
	return
}

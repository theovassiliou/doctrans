package dtaservice

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpc "google.golang.org/grpc"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
	aux "github.com/theovassiliou/doctrans/ipaux"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

const (
	appName = "DE.TU-BERLIN.QDS.ABSTRACT-SERVER"
)

// DocTransServer is a generic server
type DocTransServer struct {
	UnimplementedDTAServerServer
	//
	Register      bool   `opts:"group=Registrar" help:"Register service with EUREKA, if set"`
	RegistrarURL  string `opts:"group=Registrar" help:"Registry URL"`
	RegistrarUser string `opts:"group=Registrar" help:"Registry User, no user used if not provided"`
	RegistrarPWD  string `opts:"group=Registrar" help:"Registry User Password, no password used if not provided"`
	TTL           uint   `opts:"group=Registrar" help:"Time in seconds to reregister at Registrar."`

	ResolverURL          string `opts:"group=Resolver" help:"Resolver URL"`
	ResolverUser         string `opts:"group=Resolver" help:"Resolver User, no user used if not provided"`
	ResolverPWD          string `opts:"group=Resolver" help:"Resolver User Password, no password used if not provided"`
	ResolverTTL          uint   `opts:"group=Resolver" help:"Time in seconds to reregister at Resolver."`
	ResolverRegistration bool   `opts:"group=Resolver" help:"Register in addition also to the resolver"`

	HostName     string `opts:"group=Service" help:"If provided will be used as hostname, else automatically derived."`
	AppName      string `opts:"group=Service" help:"ID of the service"`
	PortToListen string `opts:"group=Service" help:"On which port to listen for this service."`
	DtaType      string `opts:"group=Service" help:"One of Gateway or Service. Service is assumed if not provided."`
	IsSSL        bool   `opts:"group=Service" help:"Service reached via SSL, if set."`
	REST         bool   `opts:"group=Service" help:"REST-API enabled on port 80, if set"`
	HTTPPort     string `opts:"group=Service" help:"On which httpPort to listen for REST, if enableREST is set. Ignored otherwise."`

	LogLevel log.Level `opts:"group=Generic" help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	CfgFile  string    `opts:"group=Generic" help:"The config file to use" json:"-"`
	Init     bool      `opts:"group=Generic" help:"Create a default config file as defined by cfg-file, if set. If not set ~/.dta/{AppName}/config.json will be created." json:"-"`

	registrar    *eureka.Client
	instanceInfo *eureka.InstanceInfo
	heartBeatJob *scheduler.Job
}

//CreateListener creates the grpc listener and returns it
func (dtas *DocTransServer) CreateListener(maxPortSeek int) net.Listener {
	var lis net.Listener
	var err error

	port, err := strconv.Atoi(dtas.PortToListen)

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
	dtas.PortToListen = strconv.Itoa(port)
	return lis
}

// ListServices implements dta.
func (dtas *DocTransServer) ListServices(ctx context.Context, req *ListServiceRequest) (*ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", dtas.ApplicationName())
	services := (&ListServicesResponse{}).Services
	services = append(services, dtas.ApplicationName())
	return &ListServicesResponse{Services: services}, nil
}

// ApplicationName returns the application name of the service
func (dtas *DocTransServer) ApplicationName() string {
	return appName
}

// GrpcLisInitAndReg initialises a listener for the GRPC server and registers the grps services
// Returns the listeners and the port on which the GRPC server listens
func GrpcLisInitAndReg(srvHandler *DocTransServer) net.Listener {
	lis := GrpcLisInit(srvHandler)
	// We register ourselfs by using the dyn.port
	if srvHandler.Register {
		srvHandler.RegisterAtRegistry(srvHandler.HostName, srvHandler.AppName, aux.GetIPAdress(), srvHandler.PortToListen, srvHandler.DtaType, srvHandler.TTL, srvHandler.IsSSL)
	}

	return lis
}

func GrpcLisInit(srvHandler *DocTransServer) net.Listener {
	// We first create the listener to know the dynamically allocated port we listen on
	const maxPortSeek int = 20
	_configuredPort := srvHandler.PortToListen
	lis := srvHandler.CreateListener(maxPortSeek) // for the service

	if _configuredPort != srvHandler.PortToListen {
		log.Warnf("Listing on port %v instead on configured, but used port %v\n", srvHandler.PortToListen, _configuredPort)
	}
	return lis
}

// StartGrpcServer starts the server for a given listener
func StartGrpcServer(lis net.Listener, dtaServer DTAServerServer) {
	s := grpc.NewServer()
	RegisterDTAServerServer(s, dtaServer)
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}
}

// MuxHTTPGrpc starts the HTTP server in a given context
func MuxHTTPGrpc(ctx context.Context, HTTPPort string, srvHandler *DocTransServer) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	grpcPort := srvHandler.PortToListen
	log.Debugf("GRPC Endpoint localhost:%s\n", grpcPort)
	err := RegisterDTAServerHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
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
}

package dtaservice

import (
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpc "google.golang.org/grpc"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

const (
	appName = "DE.TU-BERLIN.QDS.ABSTRACT-SERVER"
)

type DocTransServerOptions struct {
	GRPC         bool   `opts:"group=Protocols" help:"Start service only with GRPC protocol support, if set"`
	HTTP         bool   `opts:"group=Protocols" help:"Start service only with HTTP protocol support, if set"`
	Port         int    `opts:"group=Protocols" help:"On which port (starting point) to listen for the supported protocol(s)."`
	HostName     string `opts:"group=Service" help:"If provided will be used as hostname, else automatically derived."`
	RegistrarURL string `opts:"group=Registrar" help:"Registry URL (ex http://eureka:8761/eureka)"`
}
type DocTransServerGenericOptions struct {
	LogLevel log.Level `opts:"group=Generic" help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	CfgFile  string    `opts:"group=Generic" help:"The config file to use" json:"-"`
	Init     bool      `opts:"group=Generic" help:"Create a default config file as defined by cfg-file, if set. If not set ~/.dta/{AppName}/config.json will be created." json:"-"`
}

// DocTransServer is a generic server
type DocTransServer struct {
	UnimplementedDTAServerServer

	AppName string `opts:"-"`
	DtaType string `opts:"-"`
	Proto   string `opts:"-"`

	registrar    *eureka.Client
	instanceInfo *eureka.InstanceInfo
	heartBeatJob *scheduler.Job
}

//CreateListener creates the grpc listener and returns it
func CreateListener(port int, maxPortSeek int) (net.Listener, int) {
	var lis net.Listener
	var err error

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
	return lis, port
}

func InitListener(initialPort int) (net.Listener, int) {
	// We first create the listener to know the dynamically allocated port we listen on
	const maxPortSeek int = 20
	_configuredPort := initialPort
	lis, _configuredPort := CreateListener(initialPort, maxPortSeek) // for the service

	if _configuredPort != initialPort {
		log.Warnf("Listing on port %v instead on configured, but used port %v\n", initialPort, _configuredPort)
	}
	return lis, _configuredPort
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
func MuxHTTPGrpc(ctx context.Context, httpListener net.Listener, grpcPort int) {
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	log.Debugf("GRPC Endpoint localhost:%d\n", grpcPort)
	err := RegisterDTAServerHandlerFromEndpoint(ctx, gwmux, "localhost:"+strconv.Itoa(grpcPort), opts)
	if err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to register: %v", err)
	}

	// FIXME Continue here and pull the handler out. Remember to change this in all services
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, strings.NewReader("This is a test"))
	})

	mux.Handle("/", gwmux)

	// (4) Start HTTP Server
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.WithFields(log.Fields{"Service": "HTTP", "Status": "Running"}).Debugf("Starting HTTP server on: %v", httpListener.Addr().String())

	if err := http.Serve(httpListener, mux); err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}
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

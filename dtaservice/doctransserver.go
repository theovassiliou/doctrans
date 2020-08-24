package dtaservice

import (
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	sync "sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
	aux "github.com/theovassiliou/doctrans/ipaux"
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

type IDocTransServer interface {
	GetDocTransServer() GenDocTransServer
	DTAServerServer
}

// GenDocTransServer is a generic server
type GenDocTransServer struct {
	AppName string `opts:"-"`
	DtaType string `opts:"-"`
	Proto   string `opts:"-"`

	registrar    *eureka.Client
	instanceInfo *eureka.InstanceInfo
	heartBeatJob *scheduler.Job
	// UnimplementedDTAServerServer
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
func (dtas *GenDocTransServer) ListServices(ctx context.Context, req *ListServiceRequest) (*ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", dtas.ApplicationName())
	services := (&ListServicesResponse{}).Services
	services = append(services, dtas.ApplicationName())
	return &ListServicesResponse{Services: services}, nil
}

func (*GenDocTransServer) Options(context.Context, *OptionsRequest) (*OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}
func (*GenDocTransServer) TransformDocument(context.Context, *DocumentRequest) (*TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformDocument not implemented")
}
func (*GenDocTransServer) TransformPipe(context.Context, *TransformPipeRequest) (*TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}

// ApplicationName returns the application name of the service
func (dtas *GenDocTransServer) ApplicationName() string {
	return dtas.AppName
}

func LaunchServices(grpcGateway, httpGateway IDocTransServer, appName, dtaType, homepageURL string, d DocTransServerOptions) {
	var gDTS, hDTS GenDocTransServer

	var _httpListener net.Listener
	var _httpPort int

	// create GRPC Listener
	// -- take initial port
	_initialPort := d.Port
	// -- start listener and save used grpc port
	_grpcListener, _grpcPort := InitListener(_initialPort)
	_ipAddressUsed, _ := aux.ExternalIP()

	var registerGRPC, registerHTTP bool
	if grpcGateway != nil {
		registerGRPC = true
		gDTS = grpcGateway.GetDocTransServer()
		gDTS.NewInstanceInfo("grpc@"+d.HostName, appName, _ipAddressUsed, _grpcPort,
			0, false, dtaType, "grpc",
			homepageURL,
			calcStatusURL(d.RegistrarURL, appName, "grpc@"+d.HostName),
			"")
	}

	if httpGateway != nil {
		registerHTTP = true
		// create HTTP Listener (optional)
		// -- take GRPC port + 1
		// -- start listener and save used http port
		_httpListener, _httpPort = InitListener(_grpcPort + 1)
		hDTS = httpGateway.GetDocTransServer()
		hDTS.NewInstanceInfo("http@"+d.HostName, appName, _ipAddressUsed, _httpPort,
			0, false, dtaType, "http",
			homepageURL,
			calcStatusURL(d.RegistrarURL, appName, "http@"+d.HostName),
			"")

	}

	var wg sync.WaitGroup

	// Register at registrar
	// -- Register service with GRPC protocol
	log.Tracef("RegistrarURL: %s\n", d.RegistrarURL)
	if registerGRPC && d.RegistrarURL != "" {
		hDTS.RegisterAtRegistry(d.RegistrarURL)
	}
	if registerGRPC {
		go StartGrpcServer(_grpcListener, grpcGateway)
		CaptureSignals(grpcGateway, d.RegistrarURL, &wg)
		wg.Add(1)
	}

	// -- Register service with HTTP protocol (optional)
	if registerHTTP && d.RegistrarURL != "" {
		hDTS.RegisterAtRegistry(d.RegistrarURL)
	}

	if registerHTTP {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		go MuxHTTPGrpc(ctx, _httpListener, _grpcPort)
		CaptureSignals(httpGateway, d.RegistrarURL, &wg)
		wg.Add(1)
	}

	wg.Wait()
}

func calcStatusURL(url, appName, instanceID string) string {
	return url + "/apps/" + appName + "/" + instanceID
}

package main

import (
	"context"
	"net"
	"sync"

	"github.com/jpillora/opts"
	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

// TODO Splitt APP Name into appName (ECHO) and scope (DE.TU-BERLIN.QDS) so that it
// can be configured individual
const (
	appName = "DE.TU-BERLIN.QDS.ECHO"
	dtaType = "Service"
)

// Work just retuns the document (ECHO)
func Work(input []byte) (string, []string, error) {
	return string(input), []string{}, nil
}

type ServiceOptions struct {
	pb.DocTransServerOptions
	pb.DocTransServerGenericOptions
}

func calcStatusURL(url, appName, instanceId string) string {
	return url + "/apps/" + appName + "/" + instanceId
}

func NewDtaService(options ServiceOptions, appName, proto string) pb.IDocTransServer {
	var gw = DtaService{
		GenDocTransServer: pb.GenDocTransServer{
			AppName: appName,
			DtaType: dtaType,
			Proto:   proto,
		},
	}
	return &gw
}

type DtaService struct {
	pb.UnimplementedDTAServerServer
	pb.GenDocTransServer
	resolver *eureka.Client
	listener net.Listener
	pb.IDocTransServer
}

func main() {
	workingHomeDir, _ := homedir.Dir()
	homepageURL := "https://github.com/theovassiliou/doctrans/blob/master/gateway/README.md"

	serviceOptions := ServiceOptions{}
	serviceOptions.CfgFile = workingHomeDir + "/.dta/" + appName + "/config.json"
	serviceOptions.LogLevel = log.WarnLevel
	serviceOptions.HostName = aux.GetHostname()
	serviceOptions.RegistrarURL = "http://eureka:8761/eureka"

	opts.New(&serviceOptions).
		Repo("github.com/theovassiliou/doctrans").
		ConfigPath(serviceOptions.CfgFile).
		Version(VERSION).
		Parse()

	if serviceOptions.LogLevel != 0 {
		log.SetLevel(serviceOptions.LogLevel)
	}

	// Calc Configuration
	registerGRPC, registerHTTP := determineServerConfig(serviceOptions)

	var _httpListener net.Listener
	var _httpPort int

	// create GRPC Listener
	// -- take initial port
	_initialPort := serviceOptions.Port
	// -- start listener and save used grpc port
	_grpcListener, _grpcPort := pb.InitListener(_initialPort)

	// create HTTP Listener (optional)
	if registerHTTP {
		// -- take GRPC port + 1
		// -- start listener and save used http port
		_httpListener, _httpPort = pb.InitListener(_grpcPort + 1)
	}

	ipAddressUsed, _ := aux.ExternalIP()

	grpcGateway := NewDtaService(serviceOptions, appName, "grpc")
	gDTS := grpcGateway.GetDocTransServer()
	gDTS.NewInstanceInfo("grpc@"+serviceOptions.HostName, appName, ipAddressUsed, _grpcPort,
		0, false, dtaType, "grpc",
		homepageURL,
		calcStatusURL(serviceOptions.RegistrarURL, appName, "grpc@"+serviceOptions.HostName),
		"")

	httpGateway := NewDtaService(serviceOptions, appName, "http")
	hDTS := httpGateway.GetDocTransServer()
	hDTS.NewInstanceInfo("http@"+serviceOptions.HostName, appName, ipAddressUsed, _httpPort,
		0, false, dtaType, "http",
		homepageURL,
		calcStatusURL(serviceOptions.RegistrarURL, appName, "http@"+serviceOptions.HostName),
		"")

	var wg sync.WaitGroup

	// Register at registrar
	// -- Register service with GRPC protocol
	log.Tracef("RegistrarURL: %s\n", serviceOptions.RegistrarURL)
	if registerGRPC && serviceOptions.RegistrarURL != "" {
		gDTS.RegisterAtRegistry(serviceOptions.RegistrarURL)
	}
	if registerGRPC {
		go pb.StartGrpcServer(_grpcListener, grpcGateway)
		pb.CaptureSignals(grpcGateway, serviceOptions.RegistrarURL, &wg)
		wg.Add(1)
	}

	// -- Register service with HTTP protocol (optional)
	if registerHTTP && serviceOptions.RegistrarURL != "" {
		hDTS.RegisterAtRegistry(serviceOptions.RegistrarURL)
	}

	if registerHTTP {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		go pb.MuxHTTPGrpc(ctx, _httpListener, _grpcPort)
		pb.CaptureSignals(httpGateway, serviceOptions.RegistrarURL, &wg)
		wg.Add(1)
	}

	wg.Wait()
	return
}

func determineServerConfig(gwOptions ServiceOptions) (registerGRPC, registerHTTP bool) {
	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.GRPC {
		registerGRPC = true
	}

	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.HTTP {
		registerHTTP = true
	}
	return
}

// TransformDocument
func (s *DtaService) TransformDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.TransformDocumentResponse, error) {

	l, sOut, sErr := Work(req.GetDocument())
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	} else {
		errorS = []string{}
	}
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "TransformDocument"}).Tracef("Received document: %s and echoing", string(req.GetDocument()))

	return &pb.TransformDocumentResponse{
		TransDocument: []byte(l),
		TransOutput:   sOut,
		Error:         errorS,
	}, nil

}

func (s *DtaService) ListServices(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", s.ApplicationName())
	services := (&pb.ListServicesResponse{}).Services
	services = append(services, s.ApplicationName())
	return &pb.ListServicesResponse{Services: services}, nil

}
func (*DtaService) TransformPipe(context.Context, *pb.TransformPipeRequest) (*pb.TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}
func (*DtaService) Options(context.Context, *pb.OptionsRequest) (*pb.OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}

func (s *DtaService) ApplicationName() string {
	return s.AppName
}

package main

import (
	"context"
	"net"
	"sync"

	"github.com/jpillora/opts"
	"github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
	service "github.com/theovassiliou/doctrans/services/qds_transcribe/serviceimplementation"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

// TODO Splitt APP Name into appName (ECHO) and scope (DE.TU-BERLIN.QDS) so that it
// can be configured individual
const (
	appName = "BERLIN.VASSILIOUTHEO.TRANSCRIBE"
	dtaType = "Service"
)

type ServiceOptions struct {
	pb.DocTransServerOptions
	pb.DocTransServerGenericOptions
	LocalExecution      bool   `opts:"group=Local Execution, short=x" help:"If set, execute the service locally once."`
	LocalAdditionalInfo bool   `opts:"group=Local Execution, short=a" help:"Additional information on local execution. Otherwise ignored."`
	LocalFileName       string `opts:"group=Local Execution, short=f" help:"media file name if executed locally, Otherwise ignored."`
}

func calcStatusURL(url, appName, instanceId string) string {
	return url + "/apps/" + appName + "/" + instanceId
}

func NewDtaService(options ServiceOptions, appName, proto string) service.DtaService {
	var gw = service.DtaService{
		DocTransServer: pb.DocTransServer{
			AppName: appName,
			DtaType: dtaType,
			Proto:   proto,
		},
	}
	return gw
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

	if serviceOptions.LocalExecution {
		s := service.DtaService{}
		s.AppName = appName

		service.ExecuteWorkerLocally(s, serviceOptions.LocalFileName, serviceOptions.LocalAdditionalInfo)
		return
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
	grpcGateway.NewInstanceInfo("grpc@"+serviceOptions.HostName, appName, ipAddressUsed, _grpcPort,
		0, false, dtaType, "grpc",
		homepageURL,
		calcStatusURL(serviceOptions.RegistrarURL, appName, "grpc@"+serviceOptions.HostName),
		"")

	httpGateway := NewDtaService(serviceOptions, appName, "http")
	httpGateway.NewInstanceInfo("http@"+serviceOptions.HostName, appName, ipAddressUsed, _httpPort,
		0, false, dtaType, "http",
		homepageURL,
		calcStatusURL(serviceOptions.RegistrarURL, appName, "http@"+serviceOptions.HostName),
		"")

	var wg sync.WaitGroup

	// Register at registrar
	// -- Register service with GRPC protocol
	log.Tracef("RegistrarURL: %s\n", serviceOptions.RegistrarURL)
	if registerGRPC && serviceOptions.RegistrarURL != "" {
		grpcGateway.RegisterAtRegistry(serviceOptions.RegistrarURL)
	}
	if registerGRPC {
		go pb.StartGrpcServer(_grpcListener, &grpcGateway)
		pb.CaptureSignals(&grpcGateway.DocTransServer, serviceOptions.RegistrarURL, &wg)
		wg.Add(1)
	}

	// -- Register service with HTTP protocol (optional)
	if registerHTTP && serviceOptions.RegistrarURL != "" {
		httpGateway.RegisterAtRegistry(serviceOptions.RegistrarURL)
	}

	if registerHTTP {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		go pb.MuxHTTPGrpc(ctx, _httpListener, _grpcPort)
		pb.CaptureSignals(&httpGateway.DocTransServer, serviceOptions.RegistrarURL, &wg)
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

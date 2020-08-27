package main

import (
	"github.com/jpillora/opts"
	"github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
	dta "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
	service "github.com/theovassiliou/doctrans/services/qds_transcribe/serviceimplementation"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "BERLIN.VASSILIOUTHEO.TRANSCRIBE"
	dtaType = "Service"
)

type serviceCmdLineOptions struct {
	dta.DocTransServerOptions
	dta.DocTransServerGenericOptions
	LocalExecution      bool   `opts:"group=Local Execution, short=x" help:"If set, execute the service locally once."`
	LocalAdditionalInfo bool   `opts:"group=Local Execution, short=a" help:"Additional information on local execution. Otherwise ignored."`
	LocalFileName       string `opts:"group=Local Execution, short=f" help:"media file name if executed locally, Otherwise ignored."`
}

func main() {
	workingHomeDir, _ := homedir.Dir()
	homepageURL := "https://github.com/theovassiliou/doctrans/blob/master/gateway/README.md"

	serviceOptions := serviceCmdLineOptions{}
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

	var _grpcGateway, _httpGateway service.DtaService
	// Calc Configuration
	registerGRPC, registerHTTP := determineServerConfig(serviceOptions)
	if registerGRPC {
		_grpcGateway = newDtaService(serviceOptions, appName, "grpc")
	}
	if registerHTTP {
		_httpGateway = newDtaService(serviceOptions, appName, "http")
	}

	dta.LaunchServices(&_grpcGateway, &_httpGateway, appName, dtaType, homepageURL, serviceOptions.DocTransServerOptions)
}

func newDtaService(options serviceCmdLineOptions, appName, proto string) service.DtaService {
	var gw = service.DtaService{
		GenDocTransServer: dta.GenDocTransServer{
			AppName: appName,
			DtaType: dtaType,
			Proto:   proto,
		},
	}
	return gw
}

func determineServerConfig(gwOptions serviceCmdLineOptions) (registerGRPC, registerHTTP bool) {
	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.GRPC {
		registerGRPC = true
	}

	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.HTTP {
		registerHTTP = true
	}
	return
}

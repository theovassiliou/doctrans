package main

// A simple implemenation of using the Golang DocTrans Framework
import (
	"github.com/jpillora/opts"
	"github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
	dta "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
	service "github.com/theovassiliou/doctrans/services/qds_html2text/serviceimplementation"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "DE.TU-BERLIN.QDS.HTML2TEXT"
	dtaType = "Service"
)

type serviceCmdLineOptions struct {
	dta.DocTransServerOptions
	dta.DocTransServerGenericOptions
	LocalExecution string `opts:"group=Local Execution, short=x" help:"If set, execute the service locally once and read from this file"`
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

	if serviceOptions.LocalExecution != "" {
		s := service.DtaService{}
		s.AppName = appName

		service.ExecuteWorkerLocally(s, serviceOptions.LocalExecution)
		return
	}

	var _grpcGateway, _httpGateway dta.IDocTransServer
	// Calc Configuration
	registerGRPC, registerHTTP := determineServerConfig(serviceOptions)
	if registerGRPC {
		_grpcGateway = newDtaService(serviceOptions, appName, "grpc")
	}
	if registerHTTP {
		_httpGateway = newDtaService(serviceOptions, appName, "http")
	}

	dta.LaunchServices(_grpcGateway, _httpGateway, appName, dtaType, homepageURL, serviceOptions.DocTransServerOptions)
}

func newDtaService(options serviceCmdLineOptions, appName, proto string) dta.IDocTransServer {
	gw := service.DtaService{
		GenDocTransServer: dta.GenDocTransServer{
			AppName: appName,
			DtaType: dtaType,
			Proto:   proto,
		},
	}
	gw.AppName = appName

	return &gw
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

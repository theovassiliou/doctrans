package qdsservices

import (
	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
	dtaf "github.com/theovassiliou/doctrans/dtaservice"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

type QdsServer struct {
	dtaf.UnimplementedDTAServerServer

	// -- Galaxies Registrar -- The service has be registered there
	Register      bool   `opts:"group=Registrar" help:"Register service with EUREKA, if set"`
	RegistrarURL  string `opts:"group=Registrar" help:"Registry URL"`
	RegistrarUser string `opts:"group=Registrar" help:"Registry User, no user used if not provided"`
	RegistrarPWD  string `opts:"group=Registrar" help:"Registry User Password, no password used if not provided"`
	TTL           uint   `opts:"group=Registrar" help:"Time in seconds to reregister at Registrar."`

	// -- Define the service properties
	ServiceName string `opts:"group=Service" help:"ID of the service"`
	DtaType     string `opts:"group=Service" help:"One of Gateway or Service. Service is assumed if not provided."`

	// -- Describes the actual instance of the service
	HostName     string `opts:"group=Instance" help:"If provided will be used as hostname, else automatically derived."`
	GalaxyName   string `opts:"group=Instance" help:"Fully qualified galaxy name the service lives in."`
	PortToListen string `opts:"group=Instance" help:"On which port to listen for this service."`
	IsSSL        bool   `opts:"group=Instance" help:"Service reached via SSL, if set."`
	REST         bool   `opts:"group=Instance" help:"REST-API enabled on port 80, if set"`
	HTTPPort     string `opts:"group=Instance" help:"On which httpPort to listen for REST, if enableREST is set. Ignored otherwise."`

	LogLevel log.Level `opts:"group=Generic" help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	CfgFile  string    `opts:"group=Generic" help:"The config file to use" json:"-"`
	Init     bool      `opts:"group=Generic" help:"Create a default config file as defined by cfg-file, if set. If not set ~/.dta/{AppName}/config.json will be created." json:"-"`

	registrar    *eureka.Client
	instanceInfo *eureka.InstanceInfo
	heartBeatJob *scheduler.Job
}

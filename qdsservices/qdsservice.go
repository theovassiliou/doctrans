package qdsservices

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/go-eureka-client/eureka"
	"google.golang.org/grpc"

	pb "github.com/theovassiliou/doctrans/dtaservice"
)

const maxPortSeek int = 20
const EurekaURL string = "http://127.0.0.1:8761/eureka"
const EurekaTTL uint = 25

type AQdsService struct {

	// -- Galaxies Registrar -- The service has be registered there
	Register     bool   `opts:"group=Registrar" help:"Register service with EUREKA, if set"`
	RegisterURL  string `opts:"group=Registrar" help:"Registry URL"`
	RegisterUser string `opts:"group=Registrar" help:"Registry User, no user used if not provided"`
	RegisterPWD  string `opts:"group=Registrar" help:"Registry User Password, no password used if not provided"`
	TTL          uint   `opts:"group=Registrar" help:"Time in seconds to reregister at Registrar."`

	// -- Describes the actual instance of the service
	HostName         string `opts:"group=Instance" help:"If provided will be used as hostname, else automatically derived."`
	GalaxyName       string `opts:"group=Instance" help:"Fully qualified galaxy name the service lives in."`
	GRPCPortToListen string `opts:"group=Instance" help:"On which port to listen for this service."`
	IsSSL            bool   `opts:"group=Instance" help:"Service reached via SSL, if set."`
	REST             bool   `opts:"group=Instance" help:"REST-API enabled on port 80, if set"`
	HTTPPort         string `opts:"group=Instance" help:"On which httpPort to listen for REST, if enableREST is set. Ignored otherwise."`

	LogLevel log.Level `opts:"group=Generic" help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	CfgFile  string    `opts:"group=Generic" help:"The config file to use" json:"-"`
	Init     bool      `opts:"group=Generic" help:"Create a default config file as defined by cfg-file, if set. If not set ~/.dta/{AppName}/config.json will be created." json:"-"`

	registry     *eureka.Client
	instanceInfo *eureka.InstanceInfo
	heartBeatJob *scheduler.Job
}

func (serviceWorker *AQdsService) HeartBeatJob() *scheduler.Job {
	return serviceWorker.heartBeatJob
}

func (serviceWorker *AQdsService) SetHeartBeatJob(hbj *scheduler.Job) {
	serviceWorker.heartBeatJob = hbj
}

func (serviceWorker *AQdsService) InstanceInfo() *eureka.InstanceInfo {
	return serviceWorker.instanceInfo
}

func (serviceWorker *AQdsService) SetInstanceInfo(ii *eureka.InstanceInfo) {
	serviceWorker.instanceInfo = ii
}

func (serviceWorker *AQdsService) Registry() *eureka.Client {
	return serviceWorker.registry
}

func (serviceWorker *AQdsService) SetRegistry(reg *eureka.Client) {
	serviceWorker.registry = reg
}

func (serviceWorker *AQdsService) CreateListener() (net.Listener, string) {
	var lis net.Listener
	var err error

	port, err := strconv.Atoi(serviceWorker.GRPCPortToListen)

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
	return lis, strconv.Itoa(port)
}

func (serviceWorker *AQdsService) RegisterGRPCService(serviceName, ipAddress string) {
	serviceWorker.registry.CheckRetry = eureka.ExpBackOffCheckRetry
	// Create the app instance
	serviceWorker.instanceInfo = eureka.NewInstanceInfo(serviceWorker.HostName, serviceName,
		ipAddress, serviceWorker.GRPCPortToListen,
		serviceWorker.TTL, serviceWorker.IsSSL) //Create a new instance to register

	// Add some meta data. Currently no meaning
	// TODO: Remove this playground if not further required
	serviceWorker.instanceInfo.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	serviceWorker.instanceInfo.Metadata.Map["DTA-Type"] = "Service" //one of Gateway, Service
	// Register instance and heartbeat for Eureka
	serviceWorker.registry.RegisterInstance(serviceName, serviceWorker.instanceInfo) // Register new instance in your eureka(s)
	log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Init"}).Infof("Registering service %s\n", serviceName)

}

func (serviceWorker *AQdsService) GetHearbeatFunc() func() {
	return func() {
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Up"}).Trace("sending heartbeat : %v\n", time.Now().UTC())
		serviceWorker.registry.SendHeartbeat(serviceWorker.instanceInfo.App, serviceWorker.instanceInfo.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
	}
}

func (serviceWork *AQdsService) RegisterHTTPService(ctx context.Context) string {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	grpcPort := serviceWork.GRPCPortToListen
	log.Debugf("GRPC Endpoint localhost:%s\n", grpcPort)
	err := pb.RegisterDTAServerHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
	if err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to register: %v", err)
	}

	port, err := strconv.Atoi(serviceWork.HTTPPort)

	for i := 0; i < maxPortSeek; i++ {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Trying"}).Infof("Trying to start HTTP server on port %d", (port + i))

		// (4) Start HTTP Server
		// Start HTTP server (and proxy calls to gRPC server endpoint)
		err := http.ListenAndServe(":"+strconv.Itoa(port+i), mux)
		if err == nil {
			port = port + i
			log.WithFields(log.Fields{"Service": "HTTP", "Status": "Running"}).Infof("Using port %d to listen for HTTP/DTA", port)
			i = maxPortSeek
		}
	}

	if err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Infof("Failed to finally open ports between %d and %d", port, port+maxPortSeek)
		log.Fatalf("failed to listen HTTP: %v", err)
	}

	return strconv.Itoa(port)
}

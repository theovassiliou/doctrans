package dtaservice

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jpillora/opts"
	aux "github.com/theovassiliou/dta-server/ipaux"
	grpc "google.golang.org/grpc"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
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

func (dtas *DocTransServer) InstanceInfo() *eureka.InstanceInfo {
	return dtas.instanceInfo
}
func (dtas *DocTransServer) UnregisterAtRegistry() {
	if dtas.instanceInfo != nil {
		dtas.registrar.UnregisterInstance(dtas.instanceInfo.App, dtas.instanceInfo.HostName) // unregister the instance in your eureka(s)
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Unregister"}).Infof("Unregister service %s with id %s", dtas.instanceInfo.App, dtas.instanceInfo.InstanceID)
	} else {
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Unregister"}).Infof("service %s allready unregistered", dtas.instanceInfo.App)
	}
	dtas.heartBeatJob.Quit <- true
	dtas.instanceInfo = nil
}

// RegisterAtRegistry registers the DocTransServer at the Area Registry
func (dtas *DocTransServer) RegisterAtRegistry(hostname, app, ipAddress, port, dtaType string, ttl uint, isSsl bool) {

	// Build Eureka Configuration
	dtas.registrar = eureka.NewClient([]string{
		dtas.RegistrarURL,
		// add others servers here
	})

	// Create the app instance
	dtas.instanceInfo = eureka.NewInstanceInfo(hostname, app, ipAddress, port, ttl, isSsl) //Create a new instance to register
	// Add some meta data. Currently no meaning
	// TODO: Remove this playground if not further required
	dtas.instanceInfo.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	dtas.instanceInfo.Metadata.Map["DTA-Type"] = dtaType //one of Gateway, Service
	// Register instance and heartbeat for Eureka
	dtas.registrar.RegisterInstance(app, dtas.instanceInfo) // Register new instance in your eureka(s)
	log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Init"}).Infof("Registering service %s", app)

	job := func() {
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Up"}).Trace("sending heartbeat : ", time.Now().UTC())
		dtas.registrar.SendHeartbeat(dtas.instanceInfo.App, dtas.instanceInfo.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
	}

	// Run every 25 seconds but not now.
	// FIXME:0 We have somehow be able to deregister the heartbeat
	dtas.heartBeatJob, _ = scheduler.Every(25).Seconds().NotImmediately().Run(job)
}

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

// NewConfigFile creates a new example config file and terminates.
func (dtas *DocTransServer) NewConfigFile() error {

	dir := path.Dir(dtas.CfgFile)

	_, err := os.Open(dir)
	if err != nil {
		os.MkdirAll(dir, os.ModePerm)
		_, err := os.Open(dir)
		if err != nil {
			return err
		}
	}

	configJSON, _ := json.MarshalIndent(dtas, "", "  ")
	err = ioutil.WriteFile(dtas.CfgFile, configJSON, 0644)
	log.Infof("Wrote example configuration file to %s. Exiting.", dtas.CfgFile)
	return nil
}

// ListServices implements dta.
func (dtas *DocTransServer) ListServices(ctx context.Context, req *ListServiceRequest) (*ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", dtas.ApplicationName())
	services := (&ListServicesResponse{}).Services
	services = append(services, dtas.ApplicationName())
	return &ListServicesResponse{Services: services}, nil
}

// NewDocTransFromFile creates a DocTransServer from a given file path.
// The given file is expected to use the JSON format.
func NewDocTransFromFile(fpath string) (*DocTransServer, error) {
	fi, err := os.Open(fpath)
	if err != nil {
		return newDefaultDTS(), err
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	return NewDocTransFromReader(fi)
}

func newDefaultDTS() *DocTransServer {
	return &DocTransServer{
		Register:     false,
		REST:         true,
		HTTPPort:     defaultOrNot("80", os.Getenv("DTS_HTTPPort")),
		HostName:     defaultOrNot(getHostname(), os.Getenv("DTS_HostName")),
		AppName:      defaultOrNot("", os.Getenv("DTS_AppName")),
		PortToListen: defaultOrNot("50051", os.Getenv("DTS_PortToListen")),
		RegistrarURL: defaultOrNot("http://127.0.0.1:8761/eureka", os.Getenv("DTS_RegistrarURL")),
		DtaType:      "Service",
		LogLevel:     log.WarnLevel,
	}
}

// NewDocTransFromReader creates a Client configured from a given reader.
// The configuration is expected to use the JSON format.
func NewDocTransFromReader(reader io.Reader) (*DocTransServer, error) {
	d := newDefaultDTS()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
func defaultOrNot(d, v string) string {
	if v == "" {
		return d
	}
	return v
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Info("Unable to find hostname from OS")
		return ""
	}
	return hostname
}

func GetIPAdress() string {
	ipAddress, err := aux.ExternalIP()
	if err != nil {
		log.Info("Unable to find IP address from OS")
	}
	return ipAddress
}

func (dtas *DocTransServer) ApplicationName() string {
	return appName
}

func SetupConfiguration(config *DocTransServer, workingHomeDir, version string) {
	opts.New(config).
		Repo("github.com/theovassiliou/doctrans").
		Version(version).
		Parse()

	if config.LogLevel != 0 {
		log.SetLevel(config.LogLevel)
	}

	if config.AppName != "" && config.CfgFile != "" {
		config.CfgFile = workingHomeDir + "/.dta/" + config.AppName + "/config.json"
	}

	if config.Init {
		config.CfgFile = config.CfgFile + ".example"
		err := config.NewConfigFile()
		if err != nil {
			log.Fatalln(err)
		}
		log.Exit(0)
	}

	// Parse config file
	config, err := NewDocTransFromFile(config.CfgFile)
	if err != nil {
		log.Infoln("No config file found. Consider creating one using --init option.")
	}

	// Parse command line parameters again to insist on config parameters
	opts.New(config).Parse()
	if config.LogLevel != 0 {
		log.SetLevel(config.LogLevel)
	}

}

// GrpcLisInitAndReg initialises a listener for the GRPC server and registers the grps services
// Returns the listeners and the port on which the GRPC server listens
func GrpcLisInitAndReg(srvHandler *DocTransServer) net.Listener {
	// We first create the listener to know the dynamically allocated port we listen on
	const maxPortSeek int = 20
	_configuredPort := srvHandler.PortToListen

	lis := srvHandler.CreateListener(maxPortSeek) // for the service

	if _configuredPort != srvHandler.PortToListen {
		log.Warnf("Listing on port %v instead on configured, but used port %v\n", srvHandler.PortToListen, _configuredPort)
	}

	// We register ourselfs by using the dyn.port
	if srvHandler.Register {
		srvHandler.RegisterAtRegistry(srvHandler.HostName, srvHandler.AppName, GetIPAdress(), srvHandler.PortToListen, "Gateway", srvHandler.TTL, srvHandler.IsSSL)
	}
	return lis
}

func StartGrpcServer(lis net.Listener, dtaServer DTAServerServer) {
	s := grpc.NewServer()
	RegisterDTAServerServer(s, dtaServer)
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}
}

func MuxHttpGrpc(ctx context.Context, HTTPPort string, srvHandler *DocTransServer) {
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
	if err := http.ListenAndServe(":"+HTTPPort, mux); err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}
}

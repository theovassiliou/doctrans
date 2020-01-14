package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"
	"jaytaylor.com/html2text"

	"google.golang.org/grpc"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "DE.TU-BERLIN.QDS.HTML2TEXT"
)

// Work returns a nicely formatted text from a HTML input
func Work(input []byte, options []string) (string, []string, error) {
	text, err := html2text.FromString(string(input), html2text.Options{PrettyTables: true})
	return string(text), []string{}, err
}

type DtaService struct {
	pb.UnimplementedDTAServerServer
	srvHandler *pb.DocTransServer
	resolver   *eureka.Client
}

func main() {
	workingHomeDir, _ := homedir.Dir()

	dts := &pb.DocTransServer{
		RegistrarURL: "http://127.0.0.1:8761/eureka",
		AppName:      appName,
		PortToListen: "50051",
		HTTPPort:     "80",
		CfgFile:      workingHomeDir + "/.dta/" + appName + "/config.json",
		LogLevel:     log.WarnLevel,
	}

	// (1) SetUp Configuration
	setupConfiguration(dts, workingHomeDir)

	// init the resolver so that we have access to the list of apps
	// (2) Start GRPC Service
	gateway := &DtaService{
		srvHandler: dts,
		resolver: eureka.NewClient([]string{
			dts.RegistrarURL,
			// add others servers here
		}),
	}
	if !dts.REST {
		startGrpcServer(gateway, nil) // blocking
	}

	a := make(chan string)
	go startGrpcServer(gateway, a)
	grpcPort := <-a // receive the port it has registered at

	//(3) Let's instanciate the the HTTP Server
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	log.Debugf("GRPC Endpoint localhost:%s\n", grpcPort)
	err := pb.RegisterDTAServerHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
	if err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to register: %v", err)
	}

	// (4) Start HTTP Server
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.WithFields(log.Fields{"Service": "HTTP", "Status": "Running"}).Debugf("Starting HTTP server on: %v", dts.HTTPPort)
	if err := http.ListenAndServe(":"+dts.HTTPPort, mux); err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}

	return

}

func setupConfiguration(config *pb.DocTransServer, workingHomeDir string) {
	opts.New(config).
		Repo("github.com/theovassiliou/doctrans").
		Version(VERSION).
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
	config, err := pb.NewDocTransFromFile(config.CfgFile)
	if err != nil {
		log.Infoln("No config file found. Consider creating one using --init option.")
	}

	// Parse command line parameters again to insist on config parameters
	opts.New(config).Parse()
	if config.LogLevel != 0 {
		log.SetLevel(config.LogLevel)
	}

}

func handleSignals(dtaService *DtaService, signalCh chan os.Signal) {
	for sigs := range signalCh {
		switch sigs {
		case syscall.SIGTERM: // CTRL-D
			log.Debugln("Received SIGTERM")
			if dtaService.srvHandler.InstanceInfo() != nil {
				dtaService.srvHandler.UnregisterAtRegistry()
			} else {
				dtaService.srvHandler.RegisterAtRegistry(dtaService.srvHandler.HostName, dtaService.srvHandler.AppName, pb.GetIPAdress(), dtaService.srvHandler.PortToListen, "Service", dtaService.srvHandler.TTL, dtaService.srvHandler.IsSSL)
			}
		case syscall.SIGINT: // CTRL-C
			log.Debugln("Received SIGINT")
			if dtaService.srvHandler.InstanceInfo() != nil {
				dtaService.srvHandler.UnregisterAtRegistry()
			}
			os.Exit(1)
		}
	}
}

func (s *DtaService) captureSignals() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go handleSignals(s, signalCh)
}

func startGrpcServer(dtaService *DtaService, portCh chan string) {
	// We first create the listener to know the dynamically allocated port we listen on
	const maxPortSeek int = 20
	_configuredPort := dtaService.srvHandler.PortToListen

	lis := dtaService.srvHandler.CreateListener(maxPortSeek) // for the service

	if _configuredPort != dtaService.srvHandler.PortToListen {
		log.Warnf("Listing on port %v instead on configured, but used port %v\n", dtaService.srvHandler.PortToListen, _configuredPort)
	}

	if portCh != nil {
		portCh <- dtaService.srvHandler.PortToListen
	}

	s := grpc.NewServer()

	// We register ourselfs by using the dyn.port
	if dtaService.srvHandler.Register {
		dtaService.srvHandler.RegisterAtRegistry(dtaService.srvHandler.HostName, dtaService.srvHandler.AppName, pb.GetIPAdress(), dtaService.srvHandler.PortToListen, "Gateway", dtaService.srvHandler.TTL, dtaService.srvHandler.IsSSL)
	}

	dtaService.captureSignals()

	pb.RegisterDTAServerServer(s, dtaService)
	// Start dta service by using the listener
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}

}

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, docReq *pb.DocumentRequest) (*pb.TransformDocumentReply, error) {
	transResult, stdOut, stdErr := Work(docReq.GetDocument(), docReq.GetOptions())
	var errorS []string = []string{}
	if stdErr != nil {
		errorS = []string{stdErr.Error()}
	}
	log.WithFields(log.Fields{"Service": "count", "Status": "TransformDocument"}).Tracef("Received document: %s", string(docReq.GetDocument()))
	log.WithFields(log.Fields{"Service": "count", "Status": "TransformDocument"}).Tracef("Transformation Result: %s", transResult)

	return &pb.TransformDocumentReply{
		TransDocument: []byte(transResult),
		TransOutput:   stdOut,
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

func (s *DtaService) ApplicationName() string {
	return appName
}

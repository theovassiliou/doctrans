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
func Work(in *pb.DocumentRequest) (string, []string, error) {

	input := in.GetDocument()
	text, err := html2text.FromString(string(input), html2text.Options{PrettyTables: true})
	log.WithFields(log.Fields{"Service": "Work"}).Debugf("The transformation result %s\n", text)

	return string(text), []string{}, err
}

type DtaService struct {
	pb.UnimplementedDTAServerServer
	dts      *pb.DocTransServer
	resolver *eureka.Client
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
		REST:         true,
	}

	// (1) SetUp Configuration
	setupConfiguration(dts, workingHomeDir)

	// init the resolver so that we have access to the list of apps
	// (2) Start GRPC Service
	gateway := &DtaService{
		dts: dts,
		resolver: eureka.NewClient([]string{
			dts.RegistrarURL,
			// add others servers here
		}),
	}
	mux := gateway.startGRPCService()

	// (3) Start HTTP Server

	if dts.REST {
		// Start HTTP server (and proxy calls to gRPC server endpoint)
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Running"}).Debugf("Starting HTTP server on: %v", dts.HTTPPort)
		if err := http.ListenAndServe(":"+dts.HTTPPort, mux); err != nil {
			log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
		}
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

func (s *DtaService) startGRPCService() *runtime.ServeMux {

	a := make(chan string)
	go startGrpcServer(s, a)
	grpcPort := <-a // receive the port it has registered at

	// Setup context and mux
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	log.Debugf("GRPC Endpoint localhost:%s\n", grpcPort)
	if err := pb.RegisterDTAServerHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to register: %v", err)
	}
	return mux
}

func startGrpcServer(dtaService *DtaService, a chan string) {
	// We first create the listener to know the dynamically allocated port we listen on
	const maxPortSeek int = 20
	_configuredPort := dtaService.dts.PortToListen

	lis := dtaService.dts.CreateListener(maxPortSeek) // for the service

	if _configuredPort != dtaService.dts.PortToListen {
		log.Warnf("Listing on port %v instead on configured, but used port %v\n", dtaService.dts.PortToListen, _configuredPort)
	}

	a <- dtaService.dts.PortToListen

	s := grpc.NewServer()

	// We register ourselfs by using the dyn.port
	if dtaService.dts.Register {
		dtaService.dts.RegisterAtRegistry(dtaService.dts.HostName, dtaService.dts.AppName, pb.GetIPAdress(), dtaService.dts.PortToListen, "Gateway", dtaService.dts.TTL, dtaService.dts.IsSSL)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for sigs := range c {
			switch sigs {
			case syscall.SIGTERM: // CTRL-D
				if dtaService.dts.InstanceInfo() != nil {
					dtaService.dts.UnregisterAtRegistry()
				} else {
					dtaService.dts.RegisterAtRegistry(dtaService.dts.HostName, dtaService.dts.AppName, pb.GetIPAdress(), dtaService.dts.PortToListen, "Service", dtaService.dts.TTL, dtaService.dts.IsSSL)
				}
			case syscall.SIGINT: // CTRL-C
				if dtaService.dts.InstanceInfo() != nil {
					dtaService.dts.UnregisterAtRegistry()
				}
				os.Exit(1)
			}
		}
	}()
	pb.RegisterDTAServerServer(s, dtaService)
	// Start dta service by using the listener
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}

}

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, in *pb.DocumentRequest) (*pb.TransformDocumentReply, error) {
	l, sOut, sErr := Work(in)
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	} else {
		errorS = []string{}
	}
	log.WithFields(log.Fields{"Service": "count", "Status": "TransformDocument"}).Tracef("Received document: %s and has lines %s", string(in.GetDocument()), l)

	return &pb.TransformDocumentReply{
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

func (s *DtaService) ApplicationName() string {
	return appName
}

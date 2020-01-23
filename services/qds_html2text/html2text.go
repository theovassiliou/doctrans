package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"
	"jaytaylor.com/html2text"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	"github.com/theovassiliou/doctrans/qdsservices"
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

// DtaService holds the infrastructure for performing the service.
type DtaService struct {
	pb.UnimplementedDTAServerServer
	srvHandler *pb.DocTransServer
	resolver   *eureka.Client
}

func main() {
	workingHomeDir, _ := homedir.Dir()

	dts := &pb.DocTransServer{
		AppName:  appName,
		CfgFile:  workingHomeDir + "/.dta/" + appName + "/config.json",
		LogLevel: log.WarnLevel,
	}

	// (1) SetUp Configuration

	dts = pb.SetupConfiguration(dts, workingHomeDir, VERSION)
	if dts.AppName == "" {
		dts.AppName = appName
	}
	// init the resolver so that we have access to the list of apps
	service := &DtaService{
		srvHandler: dts,
	}

	// (2) Init and register GRPC Service
	lis := pb.GrpcLisInitAndReg(service.srvHandler)

	// In case of a service the resolver is the registry. Might be nil in case no registration was configured
	service.resolver = service.srvHandler.Registrar()

	go pb.StartGrpcServer(lis, service)

	// Start dta service by using the listener

	if dts.REST {
		//(3) Let's instanciate the the HTTP Server
		qdsservices.CaptureSignals(service.srvHandler)
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		pb.MuxHTTPGrpc(ctx, dts.HTTPPort, service.srvHandler)
	} else {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		qdsservices.HandleSignals(service.srvHandler, signalCh)
	}
	return
}

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, docReq *pb.DocumentRequest) (*pb.TransformDocumentResponse, error) {
	transResult, stdOut, stdErr := Work(docReq.GetDocument(), docReq.GetOptions())
	var errorS []string = []string{}
	if stdErr != nil {
		errorS = []string{stdErr.Error()}
	}
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "TransformDocument"}).Tracef("Received document: %s", string(docReq.GetDocument()))
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "TransformDocument"}).Tracef("Transformation Result: %s", transResult)

	return &pb.TransformDocumentResponse{
		TransDocument: []byte(transResult),
		TransOutput:   stdOut,
		Error:         errorS,
	}, nil
}

// ListServices list available services provided by this implementation
func (s *DtaService) ListServices(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", s.ApplicationName())
	services := (&pb.ListServicesResponse{}).Services
	services = append(services, s.ApplicationName())
	return &pb.ListServicesResponse{Services: services}, nil

}

// ApplicationName returns the name of the service application
func (s *DtaService) ApplicationName() string {
	return s.srvHandler.AppName
}

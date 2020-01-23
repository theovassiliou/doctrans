package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	"github.com/theovassiliou/doctrans/qdsservices"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

// TODO Splitt APP Name into appName (ECHO) and scope (DE.TU-BERLIN.QDS) so that it
// can be configured individual
const (
	appName = "DE.TU-BERLIN.QDS.ECHO"
)

// Work just retuns the document (ECHO)
func Work(input []byte) (string, []string, error) {
	return string(input), []string{}, nil
}

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
	gateway := &DtaService{
		srvHandler: dts,
		resolver: eureka.NewClient([]string{
			dts.RegistrarURL,
			// add others servers here
		}),
	}

	// (2) Init and register GRPC Service
	lis := pb.GrpcLisInitAndReg(gateway.srvHandler)

	go pb.StartGrpcServer(lis, gateway)

	// Start dta service by using the listener

	if dts.REST {
		//(3) Let's instanciate the the HTTP Server
		qdsservices.CaptureSignals(gateway.srvHandler)
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		pb.MuxHTTPGrpc(ctx, dts.HTTPPort, gateway.srvHandler)
	} else {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		qdsservices.HandleSignals(gateway.srvHandler, signalCh)
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

func (s *DtaService) ApplicationName() string {
	return s.srvHandler.AppName
}

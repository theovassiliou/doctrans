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
	pb.SetupConfiguration(dts, workingHomeDir, VERSION)

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
		pb.MuxHttpGrpc(ctx, dts.HTTPPort, gateway.srvHandler)
	} else {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		qdsservices.HandleSignals(gateway.srvHandler, signalCh)
	}
	return
}

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, docReq *pb.DocumentRequest) (*pb.TransformDocumentReply, error) {
	transResult, stdOut, stdErr := Work(docReq.GetDocument(), docReq.GetOptions())
	var errorS []string = []string{}
	if stdErr != nil {
		errorS = []string{stdErr.Error()}
	}
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "TransformDocument"}).Tracef("Received document: %s", string(docReq.GetDocument()))
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "TransformDocument"}).Tracef("Transformation Result: %s", transResult)

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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"regexp"
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

const (
	appName = "DE.TU-BERLIN.QDS.COUNT"
)

// CountResults describes the results of the transformation
type CountResults struct {
	Bytes int
	Lines int
	Words int
}

var re *regexp.Regexp = regexp.MustCompile(`[\S]+`)

// Work returns an encoded JSON object containing the
// bytes 	count the number of bytes
// lines	count the numnber of lines
// words		count the number of words
// The Service returns  the number of lines, words, and bytes contained in the input document
func Work(input []byte, options []string) (string, []string, error) {

	b := len(input)
	l, err := counter(bytes.NewReader(input), []byte{'\n'})
	w := len(re.FindAllString(string(input), -1))

	res := &CountResults{
		Bytes: b,
		Lines: l,
		Words: w,
	}
	resB, _ := json.MarshalIndent(res, "", "  ")
	log.WithFields(log.Fields{"Service": "Work"}).Infof("The result %s\n", resB)

	return string(resB), []string{}, err
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
	pb.SetupConfiguration(dts, workingHomeDir, VERSION)

	// init the resolver so that we have access to the list of apps
	service := &DtaService{
		srvHandler: dts,
		resolver: eureka.NewClient([]string{
			dts.RegistrarURL,
			// add others servers here
		}),
	}

	// (2) Init and register GRPC Service
	lis := pb.GrpcLisInitAndReg(service.srvHandler)

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

func counter(r io.Reader, sep []byte) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], sep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, in *pb.DocumentRequest) (*pb.TransformDocumentResponse, error) {
	l, sOut, sErr := Work(in.GetDocument(), in.GetOptions())
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	} else {
		errorS = []string{}
	}
	log.WithFields(log.Fields{"Service": "count", "Status": "TransformDocument"}).Tracef("Received document: %s and has lines %s", string(in.GetDocument()), l)

	return &pb.TransformDocumentResponse{
		TransDocument: []byte(l),
		TransOutput:   sOut,
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
	return appName
}

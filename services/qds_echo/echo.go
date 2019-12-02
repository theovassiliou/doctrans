package main

import (
	"context"

	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"

	"google.golang.org/grpc"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "DE.TU-BERLIN.QDS.ECHO"
)

// Work just retuns the document (ECHO)
func Work(input []byte) (string, []string, error) {
	return string(input), []string{}, nil
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
		PortToListen: "50051",
		CfgFile:      workingHomeDir + "/.dta/" + appName + "/config.json",
	}

	// Parse to fill the defaults
	opts.New(dts).
		Repo("github.com/theovassiliou/doctrans").
		Version(VERSION).
		Parse()

	if dts.AppName != "" {
		dts.CfgFile = workingHomeDir + "/.dta/" + dts.AppName + "/config.json"
	}

	if dts.Init {
		dts.CfgFile = dts.CfgFile + ".example"
		err := dts.NewConfigFile()
		if err != nil {
			log.Fatalln(err)
		}
		log.Exit(0)
	}

	// Parse config file
	dts, err := pb.NewDocTransFromFile(dts.CfgFile)
	if err != nil {
		log.Fatal(err)
	}

	// Parse command line parameters again to insist on config parameters
	opts.New(dts).Parse()
	log.SetLevel(dts.LogLevel)

	// init the resolver so that we have access to the list of apps
	gateway := &DtaService{
		dts: dts,
		resolver: eureka.NewClient([]string{
			dts.RegistrarURL,
			// add others servers here
		}),
	}

	// We first create the listener to know the dynamically allocated port we listen on
	const maxPortSeek int = 20
	lis := gateway.dts.CreateListener(maxPortSeek) // for the service
	s := grpc.NewServer()

	// We register ourselfs by using the dyn.port
	dts.RegisterAtRegistry(dts.HostName, dts.AppName, pb.GetIPAdress(), dts.PortToListen, "Service", dts.TTL, dts.IsSSL)

	pb.RegisterDTAServerServer(s, gateway)
	// Start dta service by using the listener
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}

}

// Transform Document
func (dtas *DtaService) TransformDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.TransformDocumentReply, error) {

	l, sOut, sErr := Work(req.GetDocument())
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	} else {
		errorS = []string{}
	}
	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocument"}).Tracef("Received document: %s and echoing", string(req.GetDocument()))

	return &pb.TransformDocumentReply{
		TransDocument: []byte(l),
		TransOutput:   sOut,
		Error:         errorS,
	}, nil

}

func (dtas *DtaService) ListServices(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", dtas.ApplicationName())
	services := (&pb.ListServicesResponse{}).Services
	services = append(services, dtas.ApplicationName())
	return &pb.ListServicesResponse{Services: services}, nil

}

func (dtas *DtaService) ApplicationName() string {
	return appName
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	"github.com/theovassiliou/doctrans/qdsservices"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "DE.TU-BERLIN.QDS.GW"
)

// Gateway holds the infrastructure for performing the service.
type Gateway struct {
	pb.UnimplementedDTAServerServer
	dts      *pb.DocTransServer
	resolver *eureka.Client
}

func main() {
	workingHomeDir, _ := homedir.Dir()

	// Take the build in defaults as configuration.
	// Basically this is required to print out something in case of --help
	dts := &pb.DocTransServer{
		AppName:  appName,
		CfgFile:  workingHomeDir + "/.dta/" + appName + "/config.json",
		LogLevel: log.WarnLevel,
	}

	// (1) SetUp Configuration
	pb.SetupConfiguration(dts, workingHomeDir, VERSION)

	// init the resolver so that we have access to the list of apps
	gateway := &Gateway{
		dts: dts,
	}

	// (2) Init and register GRPC Service
	lis := pb.GrpcLisInitAndReg(gateway.dts)

	// In case of a service the resolver is the registry. Might be nil in case no registration was configured
	gateway.resolver = gateway.dts.Registrar()

	go pb.StartGrpcServer(lis, gateway)

	// Start dta service by using the listener

	if dts.REST {
		//(3) Let's instanciate the the HTTP Server
		qdsservices.CaptureSignals(gateway.dts)
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		pb.MuxHTTPGrpc(ctx, dts.HTTPPort, gateway.dts)
	} else {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		qdsservices.HandleSignals(gateway.dts, signalCh)
	}
	return
}

// TransformDocument looks up the requested services via the resolver and forwards the request to the resolved service.
func (dtas *Gateway) TransformDocument(ctx context.Context, in *pb.DocumentRequest) (*pb.TransformDocumentResponse, error) {
	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocument"}).Debugf("Service requested: %#v", in.GetServiceName())
	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocument"}).Tracef("FileName: %v", in.GetFileName())
	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocument"}).Tracef("Option: %v", in.GetOptions())
	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocument"}).Tracef("Received: %v", string(in.GetDocument()))

	// Let's find out whether we find the server that can serve this service.
	a, err := dtas.resolver.GetApplication(in.GetServiceName())
	if err != nil || len(a.Instances) <= 0 {
		log.Errorf("Couldn't find server for app %s", in.GetServiceName())
		return &pb.TransformDocumentResponse{
			TransDocument: []byte{},
			TransOutput:   []string{},
			Error:         []string{"Could not find service", "Service requested: " + in.GetServiceName()},
		}, nil
	}
	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocument"}).Debugf("Connecting to: %s:%s", a.Instances[0].IpAddr, a.Instances[0].Port.Port)
	conn, err := grpc.Dial(a.Instances[0].IpAddr+":"+a.Instances[0].Port.Port, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocument"}).Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDTAServerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.TransformDocument(ctx, in)
	if err != nil {
		log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocument"}).Fatalf("Failed to transform: %s", err.Error())
	}
	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "TransformDocumentResult"}).Tracef("%s\n", string(r.GetTransDocument()))

	return r, err
}

// ListServices returns all the services visible for this gateway via the resolver
func (dtas *Gateway) ListServices(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServicesResponse, error) {
	// ListServices implements dtaservice.DTAServer
	a, _ := dtas.resolver.GetApplications()

	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": dtas.dts.AppName, "Status": "ListServices"}).Infof("Known Services registered with EUREKA: %v", a)
	services := (&pb.ListServicesResponse{}).Services
	for _, s := range a.Applications {
		services = append(services, s.Name)
	}
	return &pb.ListServicesResponse{Services: services}, nil
}

// ApplicationName returns just the application name
func (dtas *Gateway) ApplicationName() string {
	return appName
}

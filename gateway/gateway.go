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
	"github.com/theovassiliou/doctrans/dtaservice"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
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
	// FIXME this settings are not beeing preserved. Somehow overwritten by loading the config file
	dts := &pb.DocTransServer{
		AppName:     appName,
		CfgFile:     workingHomeDir + "/.dta/" + appName + "/config.json",
		LogLevel:    log.WarnLevel,
		DtaType:     "Gateway",
		ResolverURL: "http://eureka:8761/eureka",
	}

	// (1) SetUp Configuration
	dts = pb.SetupConfiguration(dts, workingHomeDir, VERSION)
	dts.DtaType = "Gateway"
	dts.ResolverURL = "http://eureka:8761/eureka"
	if dts.AppName == "" {
		dts.AppName = appName
	}

	// init the resolver so that we have access to the list of apps
	gateway := &Gateway{
		dts: dts,
	}

	// (2) Init and register GRPC Service
	lis := pb.GrpcLisInit(gateway.dts)

	if gateway.dts.Register {
		ip, _ := aux.ExternalIP()
		gateway.dts.RegisterAtRegistry(gateway.dts.HostName, gateway.dts.AppName, dtaservice.DefaultOrNot(ip, os.Getenv("IP")), gateway.dts.PortToListen, gateway.dts.DtaType, gateway.dts.TTL, gateway.dts.IsSSL)
	}

	// In case of a service the resolver is the registry. Might be nil in case no registration was configured

	// Build Eureka Configuration
	gateway.resolver = eureka.NewClient([]string{
		dts.ResolverURL,
		// add others servers here
	})

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

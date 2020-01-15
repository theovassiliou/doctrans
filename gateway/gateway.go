package main

import (
	"context"
	"net/http"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "DE.TU-BERLIN.QDS.GW"
)

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
		RegistrarURL: "http://127.0.0.1:8761/eureka",
		PortToListen: "50051",
		CfgFile:      workingHomeDir + "/.dta/" + appName + "/config.json",
	}

	// Parse to fill the configuration
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
	gateway := &Gateway{
		dts: dts,
		resolver: eureka.NewClient([]string{
			dts.RegistrarURL,
			// add others servers here
		}),
	}

	a := make(chan string)
	go startGrpcServer(gateway, a)
	grpcPort := <-a // receive the port it has registered at

	// Let's instanciate the the HTTP Server
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	log.Debugf("GRPC Endpoint localhost:%s\n", grpcPort)
	err = pb.RegisterDTAServerHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)

	if err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to register: %v", err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	if err = http.ListenAndServe(":8081", mux); err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}
	return

}

func startGrpcServer(gateway *Gateway, a chan string) {
	// We first create the listener to know the dynamically allocated port we listen on
	const maxPortSeek int = 20
	_configuredPort := gateway.dts.PortToListen

	lis := gateway.dts.CreateListener(maxPortSeek) // for the service

	if _configuredPort != gateway.dts.PortToListen {
		log.Warnf("Listing on port %v instead on configured, but used port %v\n", gateway.dts.PortToListen, _configuredPort)
	}

	a <- gateway.dts.PortToListen

	s := grpc.NewServer()

	// We register ourselfs by using the dyn.port
	gateway.dts.RegisterAtRegistry(gateway.dts.HostName, gateway.dts.AppName, aux.GetIPAdress(), gateway.dts.PortToListen, "Gateway", gateway.dts.TTL, gateway.dts.IsSSL)

	pb.RegisterDTAServerServer(s, gateway)
	// Start dta service by using the listener
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}

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

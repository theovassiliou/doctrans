package main

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/jpillora/opts"
	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "DE.TU-BERLIN.QDS.GW"
	dtaType = "Gateway"
)

// Gateway holds the infrastructure for performing the service
type Gateway struct {
	pb.UnimplementedDTAServerServer
	pb.DocTransServer
	resolver *eureka.Client
	listener net.Listener
}

func NewGateway(options GatewayOptions, appName, proto string) Gateway {
	var gw = Gateway{
		DocTransServer: pb.DocTransServer{
			AppName: appName,
			DtaType: dtaType,
			Proto:   proto,
		},
	}
	return gw
}

type GatewayOptions struct {
	pb.DocTransServerOptions
	PublicIP             bool   `opts:"group=Registrar" help:"Using the public IP Address if set, else external IP address"`
	ResolverURL          string `opts:"group=Resolver" help:"Resolver URL"`
	ResolverRegistration bool   `opts:"group=Resolver" help:"Register in addition also to the resolver"`
	pb.DocTransServerGenericOptions
}

func calcStatusURL(url, appName, instanceId string) string {
	return url + "/apps/" + appName + "/" + instanceId
}

func main() {
	workingHomeDir, _ := homedir.Dir()
	homepageURL := "https://github.com/theovassiliou/doctrans/blob/master/gateway/README.md"
	gwOptions := GatewayOptions{}
	gwOptions.CfgFile = workingHomeDir + "/.dta/" + appName + "/config.json"
	gwOptions.LogLevel = log.WarnLevel
	gwOptions.HostName = aux.GetHostname()
	gwOptions.ResolverURL = "http://eureka:8761/eureka"
	gwOptions.RegistrarURL = "http://eureka:8762/eureka"

	opts.New(&gwOptions).
		Repo("github.com/theovassiliou/doctrans").
		ConfigPath(gwOptions.CfgFile).
		Version(VERSION).
		Parse()

	if gwOptions.LogLevel != 0 {
		log.SetLevel(gwOptions.LogLevel)
	}

	// Calc Configuration
	registerGRPC, registerHTTP := determineServerConfig(gwOptions)

	var _httpListener net.Listener
	var _httpPort int

	// create GRPC Listener
	// -- take initial port
	_initialPort := gwOptions.Port
	// -- start listener and save used grpc port
	_grpcListener, _grpcPort := pb.InitListener(_initialPort)

	// create HTTP Listener (optional)
	if registerHTTP {
		// -- take GRPC port + 1
		// -- start listener and save used http port
		_httpListener, _httpPort = pb.InitListener(_grpcPort + 1)
	}

	extIP, _ := aux.ExternalIP()
	publicIP, _ := aux.PublicIP()
	ipAddressUsed := extIP

	if gwOptions.PublicIP && extIP != publicIP {
		ipAddressUsed = publicIP
		log.Debugf("External IP (%s) and Public IP (%s) differ. Using publicIP at registrar, if any", extIP, publicIP)
	}

	grpcGateway := NewGateway(gwOptions, appName, "grpc")
	grpcGateway.NewInstanceInfo("grpc@"+gwOptions.HostName, appName, ipAddressUsed, _grpcPort,
		0, false, dtaType, "grpc",
		homepageURL,
		calcStatusURL(gwOptions.RegistrarURL, appName, "grpc@"+gwOptions.HostName),
		"")

	httpGateway := NewGateway(gwOptions, appName, "http")
	httpGateway.NewInstanceInfo("http@"+gwOptions.HostName, appName, ipAddressUsed, _httpPort,
		0, false, dtaType, "http",
		homepageURL,
		calcStatusURL(gwOptions.RegistrarURL, appName, "http@"+gwOptions.HostName),
		"")

	// create client resolver
	grpcGateway.resolver = eureka.NewClient([]string{
		gwOptions.ResolverURL,
	})

	var wg sync.WaitGroup

	// Register at registrar
	// -- Register service with GRPC protocol
	log.Tracef("RegistrarURL: %s\n", gwOptions.RegistrarURL)
	if registerGRPC && gwOptions.RegistrarURL != "" {
		grpcGateway.RegisterAtRegistry(gwOptions.RegistrarURL)
	}
	if registerGRPC {
		go pb.StartGrpcServer(_grpcListener, &grpcGateway)
		pb.CaptureSignals(&grpcGateway.DocTransServer, gwOptions.RegistrarURL, &wg)
		wg.Add(1)
	}

	// -- Register service with HTTP protocol (optional)
	if registerHTTP && gwOptions.RegistrarURL != "" {
		httpGateway.RegisterAtRegistry(gwOptions.RegistrarURL)
	}

	if registerHTTP {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		go pb.MuxHTTPGrpc(ctx, _httpListener, _grpcPort)
		pb.CaptureSignals(&httpGateway.DocTransServer, gwOptions.RegistrarURL, &wg)
		wg.Add(1)
	}

	wg.Wait()
	return
}

func determineServerConfig(gwOptions GatewayOptions) (registerGRPC, registerHTTP bool) {
	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.GRPC {
		registerGRPC = true
	}

	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.HTTP {
		registerHTTP = true
	}
	return
}

// TransformDocument looks up the requested services via the resolver and forwards the request to the resolved service.
func (dtas *Gateway) TransformDocument(ctx context.Context, in *pb.DocumentRequest) (*pb.TransformDocumentResponse, error) {
	log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Debugf("Service requested: %#v", in.GetServiceName())
	log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Tracef("FileName: %v", in.GetFileName())
	log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Tracef("Option: %v", in.GetOptions())
	log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Tracef("Received: %v", string(in.GetDocument()))

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
	log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Debugf("Connecting to: %s:%s", a.Instances[0].IpAddr, a.Instances[0].Port.Port)
	conn, err := grpc.Dial(a.Instances[0].IpAddr+":"+a.Instances[0].Port.Port, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDTAServerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.TransformDocument(ctx, in)
	if err != nil {
		log.WithFields(log.Fields{"Service": dtas.DocTransServer.AppName, "Status": "TransformDocument"}).Fatalf("Failed to transform: %s", err.Error())
	}
	log.WithFields(log.Fields{"Service": dtas.DocTransServer.AppName, "Status": "TransformDocumentResult"}).Tracef("%s\n", string(r.GetTransDocument()))

	return r, err
}

// ListServices returns all the services visible for this gateway via the resolver
func (dtas *Gateway) ListServices(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServicesResponse, error) {
	// ListServices implements dtaservice.DTAServer
	a, _ := dtas.resolver.GetApplications()

	log.WithFields(log.Fields{"Service": dtas.DocTransServer.AppName, "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": dtas.DocTransServer.AppName, "Status": "ListServices"}).Infof("Known Services registered with EUREKA: %v", a)
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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"regexp"
	"sync"

	"github.com/jpillora/opts"
	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "DE.TU-BERLIN.QDS.COUNT"
	dtaType = "Service"
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

type ServiceOptions struct {
	pb.DocTransServerOptions
	pb.DocTransServerGenericOptions
}

func calcStatusURL(url, appName, instanceId string) string {
	return url + "/apps/" + appName + "/" + instanceId
}

func NewDtaService(options ServiceOptions, appName, proto string) DtaService {
	var gw = DtaService{
		DocTransServer: pb.DocTransServer{
			AppName: appName,
			DtaType: dtaType,
			Proto:   proto,
		},
	}
	return gw
}

// DtaService holds the infrastructure for performing the service.
type DtaService struct {
	pb.UnimplementedDTAServerServer
	pb.DocTransServer
	resolver *eureka.Client
	listener net.Listener
}

func main() {
	workingHomeDir, _ := homedir.Dir()
	homepageURL := "https://github.com/theovassiliou/doctrans/blob/master/gateway/README.md"

	serviceOptions := ServiceOptions{}
	serviceOptions.CfgFile = workingHomeDir + "/.dta/" + appName + "/config.json"
	serviceOptions.LogLevel = log.WarnLevel
	serviceOptions.HostName = aux.GetHostname()
	serviceOptions.RegistrarURL = "http://eureka:8761/eureka"

	opts.New(&serviceOptions).
		Repo("github.com/theovassiliou/doctrans").
		ConfigPath(serviceOptions.CfgFile).
		Version(VERSION).
		Parse()

	if serviceOptions.LogLevel != 0 {
		log.SetLevel(serviceOptions.LogLevel)
	}

	// Calc Configuration
	registerGRPC, registerHTTP := determineServerConfig(serviceOptions)

	var _httpListener net.Listener
	var _httpPort int

	// create GRPC Listener
	// -- take initial port
	_initialPort := serviceOptions.Port
	// -- start listener and save used grpc port
	_grpcListener, _grpcPort := pb.InitListener(_initialPort)

	// create HTTP Listener (optional)
	if registerHTTP {
		// -- take GRPC port + 1
		// -- start listener and save used http port
		_httpListener, _httpPort = pb.InitListener(_grpcPort + 1)
	}

	ipAddressUsed, _ := aux.ExternalIP()

	grpcGateway := NewDtaService(serviceOptions, appName, "grpc")
	grpcGateway.NewInstanceInfo("grpc@"+serviceOptions.HostName, appName, ipAddressUsed, _grpcPort,
		0, false, dtaType, "grpc",
		homepageURL,
		calcStatusURL(serviceOptions.RegistrarURL, appName, "grpc@"+serviceOptions.HostName),
		"")

	httpGateway := NewDtaService(serviceOptions, appName, "http")
	httpGateway.NewInstanceInfo("http@"+serviceOptions.HostName, appName, ipAddressUsed, _httpPort,
		0, false, dtaType, "http",
		homepageURL,
		calcStatusURL(serviceOptions.RegistrarURL, appName, "http@"+serviceOptions.HostName),
		"")

	var wg sync.WaitGroup

	// Register at registrar
	// -- Register service with GRPC protocol
	log.Tracef("RegistrarURL: %s\n", serviceOptions.RegistrarURL)
	if registerGRPC && serviceOptions.RegistrarURL != "" {
		grpcGateway.RegisterAtRegistry(serviceOptions.RegistrarURL)
	}
	if registerGRPC {
		go pb.StartGrpcServer(_grpcListener, &grpcGateway)
		pb.CaptureSignals(&grpcGateway.DocTransServer, serviceOptions.RegistrarURL, &wg)
		wg.Add(1)
	}

	// -- Register service with HTTP protocol (optional)
	if registerHTTP && serviceOptions.RegistrarURL != "" {
		httpGateway.RegisterAtRegistry(serviceOptions.RegistrarURL)
	}

	if registerHTTP {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		go pb.MuxHTTPGrpc(ctx, _httpListener, _grpcPort)
		pb.CaptureSignals(&httpGateway.DocTransServer, serviceOptions.RegistrarURL, &wg)
		wg.Add(1)
	}

	wg.Wait()
	return
}

func determineServerConfig(gwOptions ServiceOptions) (registerGRPC, registerHTTP bool) {
	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.GRPC {
		registerGRPC = true
	}

	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.HTTP {
		registerHTTP = true
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
	return s.AppName
}

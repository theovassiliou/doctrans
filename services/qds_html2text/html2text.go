package main

// A simple implemenation of using the Golang DocTrans Framework
import (
	"github.com/jpillora/opts"
	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"jaytaylor.com/html2text"


	log "github.com/sirupsen/logrus"
	dta "github.com/theovassiliou/doctrans/dtaservice"
	aux "github.com/theovassiliou/doctrans/ipaux"
	service "github.com/theovassiliou/doctrans/services/qds_html2text/serviceimplementation"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const (
	appName = "DE.TU-BERLIN.QDS.HTML2TEXT"
	dtaType = "Service"
)

// Work returns a nicely formatted text from a HTML input
func Work(input []byte, options *structpb.Struct) (string, []string, error) {
	text, err := html2text.FromString(string(input), html2text.Options{PrettyTables: true})
	return string(text), []string{}, err
}

type ServiceOptions struct {
	pb.DocTransServerOptions
	pb.DocTransServerGenericOptions
}

func calcStatusURL(url, appName, instanceId string) string {
	return url + "/apps/" + appName + "/" + instanceId
}

func NewDtaService(options ServiceOptions, appName, proto string) pb.IDocTransServer {
	var gw = DtaService{
		GenDocTransServer: pb.GenDocTransServer{
			AppName: appName,
			DtaType: dtaType,
			Proto:   proto,
		},
	}
	return &gw
}

// DtaService holds the infrastructure for performing the service.
type DtaService struct {
	pb.UnimplementedDTAServerServer
	pb.GenDocTransServer
	resolver *eureka.Client
	listener net.Listener
	pb.IDocTransServer
}

func main() {
	workingHomeDir, _ := homedir.Dir()
	homepageURL := "https://github.com/theovassiliou/doctrans/blob/master/gateway/README.md"

	serviceOptions := serviceCmdLineOptions{}
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

	if serviceOptions.LocalExecution != "" {
		s := service.DtaService{}
		s.AppName = appName

		service.ExecuteWorkerLocally(s, serviceOptions.LocalExecution)
		return
	}

	var _grpcGateway, _httpGateway dta.IDocTransServer
	// Calc Configuration
	registerGRPC, registerHTTP := determineServerConfig(serviceOptions)
	if registerGRPC {
		_grpcGateway = newDtaService(serviceOptions, appName, "grpc")
	}
	if registerHTTP {
		_httpGateway = newDtaService(serviceOptions, appName, "http")
	}

	dta.LaunchServices(_grpcGateway, _httpGateway, appName, dtaType, homepageURL, serviceOptions.DocTransServerOptions)
}

func newDtaService(options serviceCmdLineOptions, appName, proto string) dta.IDocTransServer {
	gw := service.DtaService{
		GenDocTransServer: dta.GenDocTransServer{
			AppName: appName,
			DtaType: dtaType,
			Proto:   proto,
		},
	}
	gw.AppName = appName

	return &gw
}

func determineServerConfig(gwOptions serviceCmdLineOptions) (registerGRPC, registerHTTP bool) {
	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.GRPC {
		registerGRPC = true
	}

	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.HTTP {
		registerHTTP = true
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
		Document: []byte(transResult),
		Output:   stdOut,
		Error:    errorS,
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

func (*DtaService) TransformPipe(ctx context.Context, req *pb.TransformPipeRequest) (*pb.TransformPipeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}

func (*DtaService) Options(context.Context, *pb.OptionsRequest) (*pb.OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}

// ApplicationName returns the name of the service application
func (s *DtaService) ApplicationName() string {
	return s.AppName
}

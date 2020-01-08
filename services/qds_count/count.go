package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"regexp"

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
	appName = "DE.TU-BERLIN.QDS.COUNT"
)

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
func Work(in *pb.DocumentRequest) (string, []string, error) {

	input := in.GetDocument()
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
		LogLevel:     log.WarnLevel,
	}

	// Parse to fill the defaults
	opts.New(dts).
		Repo("github.com/theovassiliou/doctrans").
		Version(VERSION).
		Parse()

	if dts.LogLevel != 0 {
		log.SetLevel(dts.LogLevel)
	}

	if dts.AppName != "" && dts.CfgFile != "" {
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
		log.Infoln("No config file found. Consider creating one using --init option.")
	}

	// Parse command line parameters again to insist on config parameters
	opts.New(dts).Parse()
	if dts.LogLevel != 0 {
		log.SetLevel(dts.LogLevel)
	}

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
	_configuredPort := gateway.dts.PortToListen

	lis := gateway.dts.CreateListener(maxPortSeek) // for the service

	if _configuredPort != gateway.dts.PortToListen {
		log.Warnf("Listing on port %v instead on configured, but used port %v\n", gateway.dts.PortToListen, _configuredPort)
	}

	s := grpc.NewServer()

	// We register ourselfs by using the dyn.port
	dts.RegisterAtRegistry(dts.HostName, dts.AppName, pb.GetIPAdress(), dts.PortToListen, "Service", dts.TTL, dts.IsSSL)

	pb.RegisterDTAServerServer(s, gateway)
	// Start dta service by using the listener
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}

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
func (s *DtaService) TransformDocument(ctx context.Context, in *pb.DocumentRequest) (*pb.TransformDocumentReply, error) {
	l, sOut, sErr := Work(in)
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	} else {
		errorS = []string{}
	}
	log.WithFields(log.Fields{"Service": "count", "Status": "TransformDocument"}).Tracef("Received document: %s and has lines %s", string(in.GetDocument()), l)

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
